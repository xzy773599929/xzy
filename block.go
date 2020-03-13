package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"time"
)

//1 .	 定义结构
type Block struct {
	Version uint64 //版本号
	PrevHash []byte  //前区块哈希
	MerkleRoot []byte //梅克尔根（就是一个哈希值）
	TimeStamp uint64 //时间戳
	Difficult uint64 //难度值
	Nonce uint64 //随机数，挖矿时要找的数
	Hash []byte  //当前区块哈希
	Data []byte  //数据

}

//2 .	 创建区块
func NewBlock(data string,PrevBlockHash []byte)*Block  {
	block := Block{
		Version:    00,
		PrevHash:   PrevBlockHash,
		MerkleRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficult:  100,
		Nonce:      100,
		Hash:     []byte{},
		Data:     []byte(data),
	}
	//block.SetHash()
	//创建一个pow对象
	pow := NewProofOfWork(&block)
	//不断查找随机数，不断进行哈希运算
	hash,nonce := pow.Run()
	//根据挖矿结果对区块数据进行更新
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

//生成哈希值
func (block *Block)SetHash()  {
	//拼装数据
	temp := [][]byte{
		block.PrevHash,
		block.MerkleRoot,
		Uint64ToByte(block.Version),
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficult),
		Uint64ToByte(block.Nonce),
		block.Data,
	}
	blockInfo := bytes.Join(temp,[]byte{})
	//sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}

//block类型转换成[]byte类型
func (block *Block)ToByte() []byte {
	//TODO
	return []byte{}
}

//创世区块
func GenesisBlock() *Block {
	return NewBlock("创世区块",[]byte{})
}

//辅助函数，将uint6464转换成[]byte
func Uint64ToByte(num uint64)[]byte  {
	var buffer bytes.Buffer

	err := binary.Write(&buffer,binary.BigEndian,num)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}