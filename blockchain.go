package main

type BlockChain struct {
	blocks []*Block
}

// AddBlock 在链上添加区块
func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

// NewBlockChain 创建链
func NewBlockChain() *BlockChain {
	return &BlockChain{
		blocks: []*Block{NewGenesisBlock()},
	}
}
