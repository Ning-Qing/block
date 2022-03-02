package main

import (
	"time"
)

type Block struct {
	Timestamp     int64 // 创建区块的当前事件戳
	Data          []byte
	PrevBlockHash []byte // 上一个区块的hash
	Hash          []byte
	Nonce         int // 工作量证明产生的随机值
}

// NewBlock 创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Data:          []byte(data),
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock 创建创世区块
func NewGenesisBlock() *Block {
	return NewBlock("GenesisBlock", []byte{})
}
