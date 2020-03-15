package main

import "Block/bolt"

//定义一个读取区块链的迭代器
type BlockChainIterator struct {
	db *bolt.DB
	//哈希指针，用于不断索引
	currentHashPointer []byte
}

func (bc *BlockChain)NewIterator() *BlockChainIterator  {
	
	return &BlockChainIterator{
		db:                 bc.db,
		//最初指向区块链中最后一个区块，随着调用Next方法，哈希指针不断前移
		currentHashPointer: bc.tail,
	}
}