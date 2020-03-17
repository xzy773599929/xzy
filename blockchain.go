package main

import (
	"Block/bolt"
	"log"
)

type BlockChain struct {
	//Blocks []*Block
	db *bolt.DB
	tail []byte //存储最后一个区块的哈希
}

const blockBucket  = "blockBucket"
const blockChainDB  = "blockChain.db"
//定义一个区块链
func NewBlockChain(address string) *BlockChain  {
	//创建一个创世区块，并添加到区块链
	var lastHash []byte
	//1.打开数据库
	db, err := bolt.Open(blockChainDB, 0600, nil)
	if err != nil {
		log.Panic("打开数据库失败")
	}
	//defer db.Close()
	//将要操作数据库
	//创建表，写操作
	db.Update(func(tx *bolt.Tx) error {
		//2.找到抽屉bucket（如果没有就创建）
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			//没有抽屉需要创建
			bucket, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("创建bucket(b1)失败！")
			}
			genesisBlock := GenesisBlock(address)
			bucket.Put(genesisBlock.Hash, genesisBlock.Serialize())
			bucket.Put([]byte("LastHashKey"), genesisBlock.Hash)
			//赋值最后一个区块哈希值
			lastHash = genesisBlock.Hash
		}else {
			lastHash = bucket.Get([]byte("LastHashKey"))
		}
		return nil
	})
	return &BlockChain{db,lastHash}
}

// 添加区块
func (bc *BlockChain)AddBlock(txs []*Transaction)  {
	//获取前一个区块哈希值
	db := bc.db //区块链数据库
	lasthash := bc.tail //当前链中最后一个哈希值

	db.Update(func(tx *bolt.Tx) error {
		//完成数据添加
		//找到bucket
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("bucket为空")
		}

		//a.创建新的区块
		block := NewBlock(txs,lasthash)

		//b.添加区块到区块链中
		//key为block哈希值，value为block的序列化
		bucket.Put(block.Hash,block.Serialize())
		bucket.Put([]byte("LastHashKey"), block.Hash)

		//c.更新数据库中最后一个区块哈希值
		bc.tail = block.Hash

		return nil
	})
}

//找到指定地址的所有UTXO
func (bc *BlockChain)FindUTXOs(address string)[]TXOutput  {
	var UTXO []TXOutput
	//TODO

	return UTXO
}