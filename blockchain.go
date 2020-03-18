package main

import (
	"Block/bolt"
	"fmt"
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
	//定义一个map来保存消费过的output，key是这个交易的id，value是这个交易中的索引值的数组,因为一笔交易可能有多个output都是同个地址的
	//map[交易id][]索引值
	spentOutput := make(map[string][]int64)
	//创建迭代器
	it := bc.NewIterator()
	for {
		//1.遍历区块
		block := it.Next()
		//2.遍历交易
		for _,tx := range block.Transactions {

			//3.遍历output，找到和自己地址相关的utxo（在添加output之前检查是否已经消耗过）
			for i,output := range tx.TXOutputs {
				fmt.Printf("corrent index %d\n",i)

				//这个output和我们的目标地址相同，满足条件，添加到UTXO数组中
				if output.PubKeyHash == address {
					UTXO = append(UTXO, output)
				}
			}

			//4.遍历input，找到自己花费过的utxo集合（把自己消耗过的标记出来）
			for _,input := range tx.TXInputs {
				//判断当前input的签名是否属于自己，如果和自己的地址一致，说明这个消费是自己的
				if input.Sig == address {
					indexArray := spentOutput[string(input.TXid)]
					indexArray = append(indexArray,input.Index)
				}
			}
		}
		//跳出遍历区块的循环
		if len(block.PrevHash) == 0 {
			break
		}
	}



	return UTXO
}