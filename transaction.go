package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

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
	//引用的Output所引值
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



//2.提供创建交易的方法
//3.创建挖矿交易

