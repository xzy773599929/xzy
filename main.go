package main

import "fmt"

func main()  {
	bc := NewBlockChain()
	bc.AddBlock("许泽洋获得BTC100个")
	bc.AddBlock("许泽洋获得BTC167个")
	bc.AddBlock("许泽洋获得BTC24个")

	for i , block := range bc.Blocks {
		fmt.Printf("=====区块高度:%d=====\n",i)
		fmt.Printf("前一区块哈希:%x\n",block.PrevHash)
		fmt.Printf("当前区块哈希:%x\n",block.Hash)
		fmt.Printf("区块内交易数据:%s\n",block.Data)
	}
}