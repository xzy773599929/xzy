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
	Sig string
}
//交易输出
type TXOutput struct {
	//转账金额
	Value float64
	//锁定脚本，用地址模拟
	PubKeyHash string
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
	if len(tx.TXInputs) == 1 {
		input := tx.TXInputs[0]
		//2.交易id为空
		//3.交易的index为-1
		if bytes.Equal(input.TXid,[]byte{}) && input.Index == -1 {
			return true
		}
	}
	return false
}

//2.提供创建交易的方法
//创建挖矿交易
func NewCoinBase(address string,data string)*Transaction  {
	//挖矿交易的特点：1.只有一个input，2.无须引用交易id，3.无须引用index，4.无须指定签名，所以sig字段可以有矿工自己指定，一般写矿池名
	input := TXInput{
		TXid:  []byte{},
		Index: -1,
		Sig:   data,
	}
	output := TXOutput{
		Value:      reward,
		PubKeyHash: address,
	}
	//设置交易内容
	tx := Transaction{
		TXID:      []byte{},
		TXInputs:  []TXInput{input},
		TXOutputs: []TXOutput{output},
	}
	//设置交易ID
	tx.SetHash()

	return &tx
}

//创建普通交易

func NewTransaction(from,to string,amount float64,bc *BlockChain) *Transaction  {
	//1.找到最合理的UTXO集合 map[string][]uint64
	utxos,resVal := bc.FindNeedUXTOs(from,amount)
	if resVal < amount {
		fmt.Println("余额不足，交易失败！！")
		return nil
	}

	var inputs []TXInput
	var outputs []TXOutput
	//2.创建交易输入，将这些UTXO逐一转换成input
	for id,indexArray := range utxos {
		for _,index := range indexArray {
			input := TXInput{[]byte(id),int64(index),from}
			inputs = append(inputs, input)
		}
	}

	//3.创建outputs
	output := TXOutput{amount,to}
	outputs = append(outputs, output)

	//4.如果有零钱，则要找零
	if resVal > amount {
		//找零
		outputs = append(outputs, TXOutput{resVal-amount,from})
	}

	//返回交易
	tx := Transaction{[]byte{},inputs,outputs}
	tx.SetHash()
	return &tx
}
