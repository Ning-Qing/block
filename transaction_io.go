package main

import "bytes"

type TXOutput struct {
	Value      int    // 输出的值
	PubKeyHash []byte // 公钥产生的hash
}

// Lock 签署输出
// 从地址中获取公钥的hash
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey 检查输出是否可以被 pubkey 的所有者使用
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// NewTXOutput 常见一个新的输出
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{
		Value:      value,
		PubKeyHash: nil,
	}
	txo.Lock([]byte(address))

	return txo
}

type TXInput struct {
	Txid      []byte // 一个输入引用了之前交易的一个输出,所引用的输出的交易的 ID
	Vout      int // 引用的输出在其所在交易的索引
	Signature []byte
	PubKey    []byte
}

// UsesKey 检查pubKeyHash所有者是否发起了交易
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)
	// 比较pubKey与lockingHash
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
