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
	block.SetHash()
	return &block
}

//生成哈希值
func (block *Block)SetHash()  {
	//拼装数据
	/*blcokInfo := append(block.PrevHash,block.Data...)
	blcokInfo = append(blcokInfo,block.MerkleRoot...)
	blcokInfo = append(blcokInfo,Uint64ToByte(block.Version)...)
	blcokInfo = append(blcokInfo,Uint64ToByte(block.TimeStamp)...)
	blcokInfo = append(blcokInfo,Uint64ToByte(block.Difficult)...)
	blcokInfo = append(blcokInfo,Uint64ToByte(block.Nonce)...)*/

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