package main

import "math/big"

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
	//TODO
	return []byte("sssssss"),10
}
