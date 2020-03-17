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
	addBlock --data DATA			"添加区块"
	printChain			"打印区块链"
	printChainR			"反向打印区块链"
	getBalance --address ADDRESS			"指定地址查找余额"
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
		fmt.Println("正向打印区块链")
		//打印区块
		cli.PrintBlockChain()
	case "printChainR":
		fmt.Println("反向打印区块链")
		//反向打印区块
		cli.PrintBlockChainReverse()
	case "getBalance":
		//指定地址获取余额
		//确认参数正确
		if len(args) == 4 && args[2] == "--address" {
			fmt.Println("获取余额")
			//获取地址
			address := args[3]
			//获取余额
			cli.GetBalance(address)
		}else {
			fmt.Println("请输入正确的地址参数")
			fmt.Println(Usage)
		}
	default:
		fmt.Println(Usage)
	}
}
