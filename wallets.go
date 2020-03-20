package main

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
)

//定义一个Wallets结构，它保存所有wallet以及它的地址
type Wallets struct {
	//map[地址]*Wallet
	WalletsMap map[string]*Wallet
}

//创建方法
func NewWallets() *Wallets  {
	var ws Wallets

	ws.WalletsMap = make(map[string]*Wallet)

	return &ws
}

func (ws *Wallets)CreatWallet() string {
	wallet := NewWallet()
	address := wallet.NewAddress()

	ws.WalletsMap[address] = wallet

	ws.saveToFile()
	return address
}

//保存方法，把新建的wallet加进去
func (ws *Wallets)saveToFile()  {
	var buffer bytes.Buffer

	//如果gob的Encode类型是interface或者struct中的某些字段是interface类型，需要在gob中注册可能的所有实现或者可能类型，否则会报错。
	gob.Register(elliptic.P256())
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(ws)
	if err != nil {
		log.Panic(err)
	}

	err = ioutil.WriteFile("wallet.dat",buffer.Bytes(),0600)
	if err != nil  {
		log.Panic(err)
	}
}

//读取文件方法，把所有的wallet读出来
