package main

import (
	"Block/bolt"
	"bytes"
	"crypto/ecdsa"
	"errors"
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
	//校验地址
	if !IsValidAddress(address) {
		fmt.Printf("地址无效:%s\n",address)
		panic("创世地址无效")
	}
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
	//先校验交易是否合法
	for i,tx := range txs {
		if !bc.VerifyTransaction(tx) {
			fmt.Println("矿工发现无效交易!",i)
			return
		}
	}
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
	fmt.Println("转账成功!")
}

//找到指定地址的所有UTXO
func (bc *BlockChain)FindUTXOs(pubKeyHash []byte)[]TXOutput  {
	var UTXO []TXOutput
	txs := bc.FindUTXOTransaction(pubKeyHash)
	for _,tx := range txs {
		for _,output := range tx.TXOutputs {
			if bytes.Equal(output.PubKeyHash,pubKeyHash) {
				UTXO = append(UTXO, output)
			}
		}
	}
	return UTXO
}

//寻找交易所需UTXOs
func (bc *BlockChain)FindNeedUXTOs(senderPubKeyHash []byte,amount float64)(map[string][]uint64,float64) {
	//找到的utxo集合
	utxos := make(map[string][]uint64)
	//找到的utxo里面包含的余额总数
	var cacl float64

	txs := bc.FindUTXOTransaction(senderPubKeyHash)
	for i,tx := range txs {
		for _,output := range tx.TXOutputs{
			if bytes.Equal(output.PubKeyHash,senderPubKeyHash) {
				if cacl < amount {
					//1.把utxo加进来
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)],uint64(i))
					//2.统计当前utxo总价
					cacl += output.Value
					//加完后满足条件
					//3.比较一下是否满足转账需求
					//	a.满足的话直接返回utxos，cacl
					//	b.不满足继续统计
					if cacl >= amount {
						//break
						fmt.Printf("找到了满足的金额:%f\n",cacl)
						return utxos,cacl
					}
				}else {
					fmt.Printf("不满足转账金额，当前总额:%f,  目标总额:%f\n",cacl,amount)
				}
			}
		}
	}
	return utxos,cacl
}

//寻找所需UTXO
func (bc *BlockChain)FindUTXOTransaction(senderPubKeyHash []byte) []*Transaction  {
	//var UTXO []TXOutput
	var txs []*Transaction
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
				if bytes.Equal(output.PubKeyHash,senderPubKeyHash){
					//UTXO = append(UTXO, output)
					txs = append(txs,tx)
				}

			}

			//如果当前交易是挖矿交易的话，那么直接跳过，不做遍历
			if !tx.IsCoinbase() {
				//4.遍历input，找到自己花费过的utxo集合（把自己消耗过的标记出来）
				for _,input := range tx.TXInputs {
					//判断当前input的签名是否属于自己，如果和自己的地址一致，说明这个消费是自己的
					pubKeyHash := HashPubkey(input.PubKey)
					//if input.Sig == address {
					if bytes.Equal(pubKeyHash,senderPubKeyHash) {
						indexArray := spentOutputs[string(input.TXid)]
						indexArray = append(indexArray,input.Index)
						spentOutputs[string(input.TXid)] = indexArray
					}
				}
			}
		}
		//跳出遍历区块的循环
		if len(block.PrevHash) == 0 {
			break
		}
	}
	return txs
}

//根据id查找交易本身，需要遍历整个区块链
func (bc *BlockChain)FindTransactionByTXid(id []byte)(Transaction,error)  {
	it := bc.NewIterator()
	for {
		//1.遍历区块链
		block := it.Next()
		//2.遍历交易
		for _ , tx := range block.Transactions {
			//3.比较交易，找到了直接退出
			if bytes.Equal(tx.TXID,id) {
				return *tx,nil
			}
		}

		if len(block.PrevHash) == 0 {
			fmt.Printf("区块链遍历完毕！\n")
			break
		}
	}
	//4.如果没找到，返回空Transaction，同时返回错误状态
	return Transaction{},errors.New("无效的交易id，请检查！")
}

func (bc *BlockChain)SignTransaction(tx *Transaction,privateKey *ecdsa.PrivateKey)  {
	//签名,交易创建的最后进行签名
	prevTXs := make(map[string]Transaction)
	//找到所有引用的交易
	//1.根据inputs来找，有多少input，就遍历多少次
	for _,input := range tx.TXInputs {
		//2.找到目标交易（根据TXid来找）
		txx,err := bc.FindTransactionByTXid(input.TXid)
		if err != nil {
			log.Panic(err)
		}
		//3.追加到prevTXs里面
		prevTXs[string(input.TXid)] = txx
	}

	tx.Sign(privateKey,prevTXs)
}

//签名验证
func (bc *BlockChain)VerifyTransaction(tx *Transaction)bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Transaction)
	//找到所有引用的交易
	//1.根据inputs来找，有多少input，就遍历多少次
	for _,input := range tx.TXInputs {
		//2.找到目标交易（根据TXid来找）
		txx,err := bc.FindTransactionByTXid(input.TXid)
		if err != nil {
			log.Panic(err)
		}
		//3.追加到prevTXs里面
		prevTXs[string(input.TXid)] = txx
	}
	return tx.Verify(prevTXs)
}