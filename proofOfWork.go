package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

//定义一个工作量证明结构ProofOfWork
type ProofOfWork struct {
	block *Block
	target *big.Int
}

//创建pow函数
func NewProofOfWork(block *Block)*ProofOfWork  {
	pow := ProofOfWork{
		block:  block,
	}
	//自定义难度，写成固定值,现在是string类型
	targetString := "0000100000000000000000000000000000000000000000000000000000000000"
	//引入辅助变量，目的是将上面的string转换成big.int类型
	bigIntTmp := big.Int{}
	bigIntTmp.SetString(targetString,16)
	pow.target = &bigIntTmp
	return &pow
}

//不断运行计算哈希函数
func (pow *ProofOfWork)Run()([]byte,uint64)  {
	//定义一个随机数nonce
	var nonce uint64
	var hash [32]byte
	for {
		//1.拼装数据（区块的数据，还有不断变化的随机数）
		block := pow.block
		temp := [][]byte{
			block.PrevHash,
			block.MerkleRoot,
			Uint64ToByte(block.Version),
			Uint64ToByte(block.TimeStamp),
			Uint64ToByte(block.Difficult),
			Uint64ToByte(nonce),
			//只对区块头做哈希，区块体通过梅克尔根产生影响
			//block.Data,
		}
		blockInfo := bytes.Join(temp,[]byte{})
		//2.进行哈希运算
		hash = sha256.Sum256(blockInfo)
		//3.与pow中的target进行比较
		//将得到的hash转换成big.int进行比较
		tmpInt := big.Int{}
		tmpInt.SetBytes(hash[:])
		if tmpInt.Cmp(pow.target) < 1 {
			fmt.Printf("挖矿成功！hash:%x nonce:%d\n",hash,nonce)
			return hash[:],nonce
			//找到nonce值，退出循环
		}else {
			//没找到，nonce加一，继续循环
			nonce += 1
		}
	}
}
