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
	//定义一个map来保存消费过的output，key是这个交易的id，value是这个交易中的索引值的数组,因为一笔交易可能有多个output都是同个地址的
	//map[交易id][]索引值
	spentOutputs := make(map[string][]int64)
	//创建迭代器
	it := bc.NewIterator()
	for {
		//1.遍历区块
		block := it.Next()
		//2.遍历交易
		for _,tx := range block.Transactions {

		OUTPUT:
			//3.遍历output，找到和自己地址相关的utxo（在添加output之前检查是否已经消耗过）
			for i,output := range tx.TXOutputs {
				//在这里做一个过滤，过滤消耗过的output，进行对比
				//如果相同，则不添加
				//如果当前交易的id已存在于map的key中，则说明这个交易里有消耗过的output
				if spentOutputs[string(tx.TXID)] != nil  {
					for _,j := range spentOutputs[string(tx.TXID)] {
						if int64(i) == j {
							//相等说明当前output已经消耗了，不再添加
							continue OUTPUT
						}
					}
				}

				//这个output和我们的目标地址相同，满足条件，添加到UTXO数组中
				if output.PubKeyHash == address {
					UTXO = append(UTXO, output)
				}
			}

			//如果当前交易是挖矿交易的话，那么直接跳过，不做遍历
			if !tx.IsCoinbase() {
				//4.遍历input，找到自己花费过的utxo集合（把自己消耗过的标记出来）
				for _,input := range tx.TXInputs {
					//判断当前input的签名是否属于自己，如果和自己的地址一致，说明这个消费是自己的
					if input.Sig == address {
						indexArray := spentOutputs[string(input.TXid)]
						indexArray = append(indexArray,input.Index)
					}
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

//寻找交易所需UTXOs
func (bc *BlockChain)FindNeedUXTOs(from string,amount float64)(map[string][]uint64,float64) {
	//找到的utxo集合
	var utxos map[string][]uint64
	//找到的utxo里面包含的余额总数
	var cacl float64

	//1111111111111111111
	//map[交易id][]索引值
	spentOutputs := make(map[string][]int64)
	//创建迭代器
	it := bc.NewIterator()
	for {
		//1.遍历区块
		block := it.Next()
		//2.遍历交易
		for _,tx := range block.Transactions {

		OUTPUT:
			//3.遍历output，找到和自己地址相关的utxo（在添加output之前检查是否已经消耗过）
			for i,output := range tx.TXOutputs {
				//在这里做一个过滤，过滤消耗过的output，进行对比
				//如果相同，则不添加
				//如果当前交易的id已存在于map的key中，则说明这个交易里有消耗过的output
				if spentOutputs[string(tx.TXID)] != nil  {
					for _,j := range spentOutputs[string(tx.TXID)] {
						if int64(i) == j {
							//相等说明当前output已经消耗了，不再添加
							continue OUTPUT
						}
					}
				}

				//这个output和我们的目标地址相同，满足条件，添加到UTXO数组中
				if output.PubKeyHash == from {
					//UTXO = append(UTXO, output)
					if cacl < amount {
						//1.把utxo加进来
						utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)],uint64(i))
						//2.统计当前utxo总价
						cacl += output.Value
						//加完后满足条件
						//3.比较一下是否满足转账需求
						//	a.满足的话直接返回utxos，cacl
						//	b.不满足继续统计
						if cacl > amount {
							return utxos,cacl
						}
					}
				}
			}

			//如果当前交易是挖矿交易的话，那么直接跳过，不做遍历
			if !tx.IsCoinbase() {
				//4.遍历input，找到自己花费过的utxo集合（把自己消耗过的标记出来）
				for _,input := range tx.TXInputs {
					//判断当前input的签名是否属于自己，如果和自己的地址一致，说明这个消费是自己的
					if input.Sig == from {
						indexArray := spentOutputs[string(input.TXid)]
						indexArray = append(indexArray,input.Index)
					}
				}
			}
		}
		//跳出遍历区块的循环
		if len(block.PrevHash) == 0 {
			break
		}
	}
	//222222222222222
	return utxos,cacl
}