package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
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
		Signature:[]byte(data),
		PubKey:nil,
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
	//privateKey := wallet.PrivateKey
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
	return &tx
}
