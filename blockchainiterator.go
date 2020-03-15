package main

import (
	"Block/bolt"
	"log"
)

//定义一个读取区块链的迭代器
type BlockChainIterator struct {
	db *bolt.DB
	//哈希指针，用于不断索引
	currentHashPointer []byte
}
//创建迭代器函数
func (bc *BlockChain)NewIterator() *BlockChainIterator  {
	
	return &BlockChainIterator{
		db:                 bc.db,
		//最初指向区块链中最后一个区块，随着调用Next方法，哈希指针不断前移
		currentHashPointer: bc.tail,
	}
}
//实现Next函数
func (it *BlockChainIterator)Next() *Block  {
	//获取数据库
	db := it.db
	block := Block{}
	//读数据库view
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("读bucket为空")
		}
		//根据哈希指针获取数据
		data := bucket.Get(it.currentHashPointer)
		//数据反序列化获得block结构体
		block = Deserialize(data)
		//哈希指针前移
		it.currentHashPointer = block.PrevHash
		return nil
	})
	return &block
}