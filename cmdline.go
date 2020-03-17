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
		fmt.Printf("============================\n\n")
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
	//TODO
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