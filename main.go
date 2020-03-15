package main

import "fmt"

func main()  {
	bc := NewBlockChain()
	bc.AddBlock("许泽洋获得BTC100个")
	bc.AddBlock("许泽洋获得BTC167个")
	bc.AddBlock("许泽洋获得BTC24个")

	//创建迭代器
	it := bc.NewIterator()
	//开始循环读取数据库数据
	for {
		block := it.Next()
		fmt.Printf("============================\n\n")
		fmt.Printf("前一区块哈希:%x\n",block.PrevHash)
		fmt.Printf("当前区块哈希:%x\n",block.Hash)
		fmt.Printf("区块内交易数据:%s\n",block.Data)
		if len(block.PrevHash) == 0 {
			fmt.Println("迭代器迭代完毕")
			break
		}
	}
}