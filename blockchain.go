package main

type BlockChain struct {
	Blocks []*Block
}

//定义一个区块链
func NewBlockChain() *BlockChain  {
	//创建一个创世区块，并添加到区块链
	gensisBlock := GenesisBlock()
	return &BlockChain{
		Blocks : []*Block{gensisBlock},
	}
}

// 添加区块
func (bc *BlockChain)AddBlock(data string)  {
	//获取区块链中最后一个区块
	lastBlock := bc.Blocks[len(bc.Blocks)-1]
	//获取前一个区块哈希值
	prevHash := lastBlock.Hash
	//创建一个新的区块
	block := NewBlock(data,prevHash[:])
	//添加区块到区块链中
	bc.Blocks = append(bc.Blocks,block)
}
