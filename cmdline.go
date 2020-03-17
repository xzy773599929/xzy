package main

import "fmt"

func (cli *CLI)AddBlock(data string)  {
	//cli.bc.AddBlock(data)
	//TODO
	fmt.Println("区块写入成功！")
}

func (cli *CLI)PrintBlockChain()  {
	bc := cli.bc
	//创建迭代器
	it := bc.NewIterator()
	//开始循环读取数据库数据bl
	for {
		block := it.Next()
		fmt.Printf("\n============================\n")
		fmt.Printf("版本号:%d\n",block.Version)
		fmt.Printf("前一区块哈希:%x\n",block.PrevHash)
		fmt.Printf("梅克尔根:%x\n",block.MerkleRoot)
		fmt.Printf("时间戳:%d\n",block.TimeStamp)
		fmt.Printf("难度值:%d\n",block.Difficult)
		fmt.Printf("随机数:%d\n",block.Nonce)
		fmt.Printf("当前区块哈希:%x\n",block.Hash)
		fmt.Printf("区块内交易数据:%s\n",block.Transactions[0].TXInputs[0].Sig)
		if len(block.PrevHash) == 0 {
			fmt.Println("区块链遍历完毕")
			break
		}
	}
}

//反向打印区块链
func (cli *CLI)PrintBlockChainReverse()  {
	bc := cli.bc
	//创建迭代器
	it := bc.NewIterator()
	//开始循环读取数据库数据
	i := 1
	for {
		block := it.Next()
		if len(block.PrevHash) == 0 {
			break
		}
		i += 1
	}

	blocks := make([]*Block,i)
	it1 := bc.NewIterator()
	l := 0
	for {
		blocks[l] = it1.Next()
		if len(blocks[l].PrevHash) == 0 {
			break
		}
		l += 1
	}
	for m := len(blocks)-1 ; i >= 0 ; i-- {
		fmt.Printf("\n============================区块高度:%d\n",m)
		fmt.Printf("版本号:%d\n",blocks[m].Version)
		fmt.Printf("前一区块哈希:%x\n",blocks[m].PrevHash)
		fmt.Printf("梅克尔根:%x\n",blocks[m].MerkleRoot)
		fmt.Printf("时间戳:%d\n",blocks[m].TimeStamp)
		fmt.Printf("难度值:%d\n",blocks[m].Difficult)
		fmt.Printf("随机数:%d\n",blocks[m].Nonce)
		fmt.Printf("当前区块哈希:%x\n",blocks[m].Hash)
		fmt.Printf("区块内交易数据:%s\n",blocks[m].Transactions[0].TXInputs[0].Sig)
	}
	fmt.Println("区块链遍历完毕")
}

//指定地址获取余额
func (cli *CLI)GetBalance(address string) {
	bc := cli.bc
	utxos := bc.FindUTXOs(address)
	//定义总余额,遍历utxos中的value进行累加
	total := 0.0
	for _,utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("\"%s\"余额为:%f\n",address,total)
}