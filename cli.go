package main

import (
	"fmt"
	"os"
	"strconv"
)
//这是一个用来接收命令行参数，并控制区块链操作的文件
//定义一个用来操作命令行的结构体
type CLI struct {
	bc *BlockChain
}
//提示语，提示可用命令
const Usage = `
	printChain			"打印区块链"
	printChainR			"反向打印区块链"
	getBalance --address ADDRESS			"指定地址查找余额"
	send FROM TO AMOUNT MINER DATA		"由FROM转AMOUNT金额给TO，由MINER挖矿，同时写入DATA"
	newWallet			"创建一个新钱包"
	listAddress			"显示所有地址"
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
	/*case "addBlock":
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
		}*/
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
	case "send":
		fmt.Println("转账开始...")
		//block.exe send FROM TO AMOUNT MINER DATA
		if len(args) == 7 {
			from := args[2]
			to := args[3]
			amount,_ := strconv.ParseFloat(args[4],64) // 字符串转float
			miner := args[5]
			data := args[6]
			cli.send(from,to,amount,miner,data)
		}else {
			fmt.Println("请输入正确的转账参数!")
			fmt.Println(Usage)
		}
	case "newWallet":
		fmt.Println("正在创建钱包...")
		cli.newWallet()
	case "listAddresses":
		fmt.Println("正在获取所有地址...")
		cli.ListAllAddress()
	default:
		fmt.Println("无效的命令!")
		fmt.Println(Usage)
	}
}
