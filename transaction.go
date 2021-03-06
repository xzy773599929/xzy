package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"math/big"
	"strings"
)

//挖矿奖励金额
const reward = 12.5 

//1.定义交易结构
type Transaction struct {
	//交易ID
	TXID []byte
	//交易输入数组
	TXInputs []TXInput
	//交易输出数组
	TXOutputs []TXOutput
}
//交易输入
type TXInput struct {
	//引用的交易ID
	TXid []byte
	//引用的Output索引值
	Index int64
	//解锁脚本，用地址模拟
	//Sig string

	//真正的数字签名，由r，s拼接成的[]byte
	Signature []byte

	//公钥，这里不是保存原始公钥，而是公钥结构体中的X和Y拼接的字符串，在校验时重新拆分（参考r，s传递)
	PubKey []byte
}
//交易输出
type TXOutput struct {
	//转账金额
	Value float64
	//锁定脚本，用地址模拟
	//PubKeyHash string

	//收款方的公钥的哈希，注意：是哈希不是公钥，也不是地址
	PubKeyHash []byte
}

//由于现在储存的字段是地址的公钥哈希，所以无法直接创建TXOutput，为了能得到公钥哈希，写一个Lock方法将地址转换成公钥哈希
func (output *TXOutput)Lock(address string)  {

	//真正的锁定动作！！！
	output.PubKeyHash = GetPubKeyHashFromAddress(address)
}

//定义一个创建TXOutput的方法，否则无法调用lock
func NewOutput(address string,value float64) *TXOutput {
	var output TXOutput
	output.Value = value
	output.Lock(address)
	return &output
}

//设置交易ID
func (tx *Transaction)SetHash()  {
	var buffer bytes.Buffer
	//创建一个gob编码器，进行序列化
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic(err)
	}

	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}

//实现一个函数，判断当前交易是否挖矿交易,如果是就返回true
func (tx *Transaction)IsCoinbase()bool  {
	//1.交易input只有一个
	input := tx.TXInputs[0]
	//2.交易id为空
	//3.交易的index为-1
	if len(tx.TXInputs) == 1 && bytes.Equal(input.TXid,[]byte{}) && input.Index == -1 {
		return true
	}
	return false
}

//2.提供创建交易的方法
//创建挖矿交易
func NewCoinBase(address string,data string)*Transaction  {
	//挖矿交易的特点：1.只有一个input，2.无须引用交易id，3.无须引用index，4.无须指定签名，所以PubKey字段可以有矿工自己指定，一般写矿池名
	//签名先填写为空，后面创建完整交易后，最后做一次签名即可
	input := TXInput{
		TXid:  []byte{},
		Index: -1,
		Signature:nil,
		PubKey:[]byte(data),
	}
	//output := TXOutput{
	//	Value:      reward,
	//	PubKeyHash: address,
	//}
	//新的创建方法
	output := NewOutput(address,reward)
	//设置交易内容
	tx := Transaction{
		TXID:      []byte{},
		TXInputs:  []TXInput{input},
		TXOutputs: []TXOutput{*output},
	}
	//设置交易ID
	tx.SetHash()

	return &tx
}

//创建普通交易
func NewTransaction(from,to string,amount float64,bc *BlockChain) *Transaction  {
	//1.创建交易之后要进行数字签名，所以需要私钥，需要打开钱包NewWallets()
	ws := NewWallets()
	//2.找到自己的钱包，根据地址返回自己的wallet
	wallet := ws.WalletsMap[from]
	if wallet == nil {
		fmt.Println("地址对应的钱包不存在")
		return nil
	}
	//3.得到对应的公钥私钥
	privateKey := wallet.PrivateKey
	pubKey := wallet.PubKey
	//传递公钥的哈希而不是地址
	pubKeyHash := HashPubkey(pubKey)

	//1.找到最合理的UTXO集合 map[string][]uint64
	utxos,resVal := bc.FindNeedUXTOs(pubKeyHash,amount)
	if resVal < amount {
		fmt.Println("余额不足，交易失败！！")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput
	//2.创建交易输入，将这些UTXO逐一转换成input
	for id,indexArray := range utxos {
		for _,index := range indexArray {
			input := TXInput{[]byte(id),int64(index),nil,pubKey}
			inputs = append(inputs, input)
		}
	}

	//3.创建outputs
	//output := TXOutput{amount,to}
	output := NewOutput(to,amount)
	outputs = append(outputs, *output)

	//4.如果有零钱，则要找零
	if resVal > amount {
		//找零
		output := NewOutput(from, resVal - amount)
		outputs = append(outputs, *output)
	}
	//返回交易
	tx := Transaction{[]byte{},inputs,outputs}
	tx.SetHash()
	//签名
	bc.SignTransaction(&tx,privateKey)

	return &tx
}

//签名的具体实现，参数为：私钥，inputs里面所有引用的交易结构map[交易Id]Transaction
func (tx *Transaction)Sign(private *ecdsa.PrivateKey,prevTXs map[string]Transaction) {
	//挖矿交易不需要签名
	if tx.IsCoinbase() {
		return
	}
	//1.创建一个当前交易的txCopy：TrimmedCopy()，要把signature和PubKey字段设为nil
	txCopy := tx.TrimmedCopy()
	//2.循环遍历txCopy的inputs，得到这个input所引用的output的公钥哈希,将这个公钥哈希暂时填充到txCopy中的input的PubKey字段中，用于签名用
	for i,input := range txCopy.TXInputs {
		prevTX := prevTXs[string(input.TXid)]
		if len(prevTX.TXID) == 0 {
			log.Panic("所引用的交易无效!")
		}
		txCopy.TXInputs[i].PubKey = prevTX.TXOutputs[input.Index].PubKeyHash
		//3.生成要签名的数据，要签名的数据一定是一个哈希值
			//a.对每一个input都要签名一次，签名的数据是当前input引用的output的哈希+当前的outputs(都在txCopy里面)
			//b.要对这个拼好的txCopy进行哈希处理，SetHash得到TXID,这个TXID就是我们要签名的最终数据
		txCopy.SetHash()
		//还原,以免影响后面的input的签名
		txCopy.TXInputs[i].PubKey = nil
		signDataHash := txCopy.TXID
		//4.执行签名动作得到r，s字节流
		r,s,err := ecdsa.Sign(rand.Reader,private,signDataHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(),s.Bytes()...)
		//5.放到我们所签名的input的Signature中
		tx.TXInputs[i].Signature = signature
	}
}

//创建交易副本,把signature和PubKey字段设为nil
func (tx *Transaction)TrimmedCopy() Transaction  {
	var inputs []TXInput
	for _,input := range tx.TXInputs{
		inputs = append(inputs, TXInput{input.TXid,input.Index,nil,nil})
	}
	return Transaction{tx.TXID,inputs,tx.TXOutputs}
}

//签名校验
func (tx *Transaction)Verify(prevTXs map[string]Transaction)bool {
	//挖矿交易不需要验证
	if tx.IsCoinbase() {
		return true
	}
	//1.得到签名的数据
	txCopy := tx.TrimmedCopy()
	//与签名时类似
	for i ,input := range tx.TXInputs {
		prevTx := prevTXs[string(input.TXid)]
		if len(prevTx.TXID) == 0 {
			log.Panic("所引用的交易无效!")
		}
		txCopy.TXInputs[i].PubKey = prevTx.TXOutputs[input.Index].PubKeyHash
		txCopy.SetHash()
		//还原,以免影响后面的的签名验证
		txCopy.TXInputs[i].PubKey = nil
		dataHash := txCopy.TXID
		//2.得到signature，反推回r,s
		signature := input.Signature
		//3.拆解PubKey，X,Y得到原生公钥
		pubKey := input.PubKey

		//定义两个辅助的big.Int
		r := big.Int{}
		s := big.Int{}
		//拆分signature，平均分成两半，前半部分为r，后半部分为s
		r.SetBytes(signature[0:len(signature)/2])
		s.SetBytes(signature[len(signature)/2:])

		//定义两个辅助的big.Int
		X := big.Int{}
		Y := big.Int{}
		//拆分pubKey，平均分成两半，前半部分为X，后半部分为Y
		X.SetBytes(pubKey[0:len(pubKey)/2])
		Y.SetBytes(pubKey[len(pubKey)/2:])
		pubKeyOrigin := ecdsa.PublicKey{elliptic.P256(),&X,&Y}

		//验证签名,如果验证失败，返回false
		if !ecdsa.Verify(&pubKeyOrigin,dataHash,&r,&s) {
			fmt.Println(i)
			return false
		}
	}
	return true
}

func(tx	Transaction)String()string	{
	var	lines	[]string
	lines	=	append(lines,	fmt.Sprintf("---Transaction	%x:",tx.TXID))
	for	i,	input	:=	range	tx.TXInputs	{
		lines	=	append(lines,	fmt.Sprintf("		Input		%d:",i))
		lines	=	append(lines,	fmt.Sprintf("		TXID:		%x",input.TXid))
		lines	=	append(lines,	fmt.Sprintf("		Out:		%d",input.Index))
		lines	=	append(lines,	fmt.Sprintf("		Signature:	%x",input.Signature))
		lines	=	append(lines,	fmt.Sprintf("		PubKey:		%x",input.PubKey))
	}
	for	i,	output	:=	range	tx.TXOutputs{
		lines	=	append(lines,	fmt.Sprintf("		Output		%d:",i))
		lines	=	append(lines,	fmt.Sprintf("		Value:		%f",output.Value))
		lines	=	append(lines,	fmt.Sprintf("		Script:		%x",output.PubKeyHash))
	}
	return	strings.Join(lines,	"\n")
}