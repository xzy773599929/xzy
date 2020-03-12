package main

import (
	"crypto/sha256"
	"fmt"
)

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


func main()  {
	block := NewBlock("许泽洋获得50个BTC",[]byte{})

	fmt.Printf("前一个区块哈希:%x\n",block.PrevHash)
	fmt.Printf("当前区块哈希:%x\n",block.Hash)
	fmt.Printf("区块内交易数据:%s\n",block.Data)
	//3 .	 ⽣成哈希
	//4 .	 引⼊区块链
	//5 .	 添加区块
	//6 .	 重构代码
}