package main

import (
	"fmt"
	"time"
)

//打印区块链
func (cli *CLI)PrintBlockChain()  {
	bc := cli.bc
	//创建迭代器
	it := bc.NewIterator()
	//先利用i获取区块链长度
	itx := bc.NewIterator()
	i := 1
	for {
		blockx := itx.Next()
		if len(blockx.PrevHash) == 0 {
			break
		}
		i += 1
	}
	//开始循环读取数据库数据bl
	for {
		block := it.Next()
		fmt.Printf("\n============================区块高度:%d\n",i-1)
		i = i-1
		fmt.Printf("版本号:%d\n",block.Version)
		fmt.Printf("前一区块哈希:%x\n",block.PrevHash)
		fmt.Printf("梅克尔根:%x\n",block.MerkleRoot)
		timeFormat := time.Unix(int64(block.TimeStamp),0).Format("2006-02-02 15:04:05")
		fmt.Printf("时间戳:%s\n",timeFormat)
		fmt.Printf("难度值:%d\n",block.Difficult)
		fmt.Printf("随机数:%d\n",block.Nonce)
		fmt.Printf("当前区块哈希:%x\n",block.Hash)
		fmt.Printf("区块内交易数据:%s\n",block.Transactions[0].TXInputs[0].Signature)
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
	//先利用i获取区块链长度
	i := 1
	for {
		block := it.Next()
		if len(block.PrevHash) == 0 {
			break
		}
		i += 1
	}
	//定义数组blocks，将区块逐个放入
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
	//遍历数组，将区块反向打印
	for m := len(blocks)-1 ; m >= 0 ; m-- {
		fmt.Printf("\n============================区块高度:%d\n",len(blocks)-1-m)
		fmt.Printf("版本号:%d\n",blocks[m].Version)
		fmt.Printf("前一区块哈希:%x\n",blocks[m].PrevHash)
		fmt.Printf("梅克尔根:%x\n",blocks[m].MerkleRoot)
		timeFormat := time.Unix(int64(blocks[m].TimeStamp),0).Format("2006-02-02 15:04:05")
		fmt.Printf("时间戳:%s\n",timeFormat)
		fmt.Printf("难度值:%d\n",blocks[m].Difficult)
		fmt.Printf("随机数:%d\n",blocks[m].Nonce)
		fmt.Printf("当前区块哈希:%x\n",blocks[m].Hash)
		fmt.Printf("区块内交易数据:%s\n",blocks[m].Transactions[0].TXInputs[0].Signature)
	}
	fmt.Println("区块链遍历完毕")
}

//指定地址获取余额
func (cli *CLI)GetBalance(address string) {
	bc := cli.bc
	//校验地址
	//TODO
	//先根据地址获得公钥哈希
	pubKeyHash := GetPubKeyHashFromAddress(address)
	utxos := bc.FindUTXOs(pubKeyHash)
	//定义总余额,遍历utxos中的value进行累加
	total := 0.0
	for _,utxo := range utxos {
		total += utxo.Value
	}
	fmt.Printf("\"%s\"余额为:%f\n",address,total)
}

//转账交易
func (cli *CLI)send(from,to string,amount float64,miner,data string)  {
	fmt.Printf("from:%s\n",from)
	fmt.Printf("to:%s\n",to)
	fmt.Printf("amount:%f\n",amount)
	fmt.Printf("miner:%s\n",miner)
	fmt.Printf("data:%s\n",data)

	//1.创建挖矿交易
	coinbase := NewCoinBase(miner,data)
	//2.创建普通交易
	transaction := NewTransaction(from,to,amount,cli.bc)
	if transaction == nil {
		return
	}
	//3.添加区块
	cli.bc.AddBlock([]*Transaction{coinbase,transaction})
	fmt.Println("转账成功!")
}

//创建新钱包
func (cli *CLI)newWallet()  {
	//wallet := NewWallet()
	//address := wallet.NewAddress()
	ws := NewWallets()
	address := ws.CreatWallet()
	fmt.Println(address)
	//for address := range wallets.WalletsMap {
	//	fmt.Printf("地址:%s\n",address)
	//}
	//fmt.Printf("私钥:%v\n",wallet.privateKey)
	//fmt.Printf("公钥:%v\n",wallet.pubKey)
}

//获取钱包所有地址
func (cli *CLI)ListAllAddress()  {
	ws := NewWallets()
	addresses := ws.ListAllAddress()
	for i,address := range addresses {
		fmt.Printf("地址%d:%s\n",i+1,address)
	}
}