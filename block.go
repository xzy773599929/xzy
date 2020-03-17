package main

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"os"
	"time"
)

//1 .	 定义结构
type Block struct {
	Version uint64 //版本号
	PrevHash []byte  //前区块哈希
	MerkleRoot []byte //梅克尔根（就是一个哈希值）
	TimeStamp uint64 //时间戳
	Difficult uint64 //难度值
	Nonce uint64 //随机数，挖矿时要找的数
	Hash []byte  //当前区块哈希
	//Data []byte  //数据
	//真实的交易数据
	Transactions []*Transaction
}

//2 .	 创建区块
func NewBlock(txs []*Transaction,PrevBlockHash []byte)*Block  {
	block := Block{
		Version:    00,
		PrevHash:   PrevBlockHash,
		MerkleRoot: []byte{},
		TimeStamp:  uint64(time.Now().Unix()),
		Difficult:  100,
		Nonce:      100,
		Hash:       []byte{},
		Transactions: txs,
	}
	block.MerkleRoot = block.MakeMerkelRoot()
	//block.SetHash()
	//创建一个pow对象
	pow := NewProofOfWork(&block)
	//不断查找随机数，不断进行哈希运算
	hash,nonce := pow.Run()
	//根据挖矿结果对区块数据进行更新
	block.Hash = hash
	block.Nonce = nonce
	return &block
}

/*//生成哈希值
func (block *Block)SetHash()  {
	//拼装数据
	temp := [][]byte{
		block.PrevHash,
		block.MerkleRoot,
		Uint64ToByte(block.Version),
		Uint64ToByte(block.TimeStamp),
		Uint64ToByte(block.Difficult),
		Uint64ToByte(block.Nonce),
		block.Data,
	}
	blockInfo := bytes.Join(temp,[]byte{})
	//sha256
	hash := sha256.Sum256(blockInfo)
	block.Hash = hash[:]
}*/

//block类型转换成[]byte类型，序列化
func (block *Block)Serialize() []byte {
	//使用gob进行序列化（编码）得到字节流
	var buffer bytes.Buffer
	//定义一个编码器
	encode := gob.NewEncoder(&buffer)
	//使用编码器进行编码
	err := encode.Encode(&block)
	if err != nil {
		fmt.Println("编码失败:",err)
		os.Exit(1)
	}
	return buffer.Bytes()
}

//反序列化
func Deserialize(data []byte) Block {
	var buffer bytes.Buffer
	var block Block

	//将data写入buffer
	_,err := buffer.Write(data)
	if err != nil {
		fmt.Println("buffer.Read failed:",err)
		os.Exit(1)
	}
	//创建decoder
	decoder := gob.NewDecoder(&buffer)
	//将buffer数据转换成block
	err = decoder.Decode(&block)
	if err != nil {
		fmt.Println("decode failed:",err)
		os.Exit(1)
	}
	return block
}

//创世区块
func GenesisBlock(address string) *Block {
	coinbase := NewCoinBase(address,"创世区块")
	return NewBlock([]*Transaction{coinbase},[]byte{})
}

//辅助函数，将uint6464转换成[]byte
func Uint64ToByte(num uint64)[]byte  {
	var buffer bytes.Buffer

	err := binary.Write(&buffer,binary.BigEndian,num)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

//模拟计算MerkelRoot
func (block *Block)MakeMerkelRoot() []byte{
	//TODO
	return []byte{}
}