package main

import "crypto/sha256"

//1 .	 定义结构
type Block struct {
	PrevHash []byte  //前区块哈希
	Hash []byte  //当前区块哈希
	Data []byte  //数据
}

//2 .	 创建区块
func NewBlock(data string,PrevBlockHash []byte)*Block  {
	block := Block{
		PrevHash: PrevBlockHash,
		Hash:     []byte{},
		Data:     []byte(data),
	}
	block.SetHash()
	return &block
}

//生成哈希值
func (block *Block)SetHash()  {
	//拼装数据
	blcokInfo := append(block.PrevHash,block.Data...)
	//sha256
	hash := sha256.Sum256(blcokInfo)
	block.Hash = hash[:]
}

//创世区块
func GenesisBlock() *Block {
	return NewBlock("创世区块",[]byte{})
}