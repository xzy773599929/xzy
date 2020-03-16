package main

import (
	"fmt"
	"os"
)
//这是一个用来接收命令行参数，并控制区块链操作的文件
//定义一个用来操作命令行的结构体
type CLI struct {
	bc *BlockChain
}
//提示语，提示可用命令
const Usage = `
	addBlock --data DATA			"add data to blockchain"
	printChain			"print all blockchain data"
`

func (cli *CLI)Run()  {
	//1.得到所有命令
	args := os.Args
	if len(args) < 2 {
		fmt.Println(Usage)
		return
	}
	//2.分析命令
	cmd := args[1]
	switch cmd {
	case "addBlock":
		//3.执行相应操作
		//添加区块,确认添加区块参数正确
		if len(args) == 4 && args[2] == "--data" {
			//获取数据
			data := args[3]
			//写入区块链
			cli.AddBlock(data)
		}else {
			fmt.Println("请输入正确的添加区块参数")
			fmt.Println(Usage)
		}
	case "printChain":
		//打印区块
		cli.PrintBlockChain()
	default:
		fmt.Println(Usage)
	}
}
