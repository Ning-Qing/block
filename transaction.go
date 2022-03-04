package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

const subsidy = 10 // subsidu 发币量

type Transaction struct {
	ID   []byte     // 交易ID
	Vin  []TXInput  // 交易的输入集
	Vout []TXOutput // 交易的输出集
}

// SetID 设置交易的 ID
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic(err)
	}
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// IsCoinbase 检查交易是否是 coinbase
func (tx Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

type TXOutput struct {
	Value        int    // 输出的值
	ScriptPubKey string //
}

// CanBeUnlockedWith 检查输出是否可以使用提供的数据解锁
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

// CanUnlockOutputWith 检查地址是否发起了交易
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

// NewCoinbaseTX 创建一个coinbase交易
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{
		Txid:      []byte{},
		Vout:      -1,
		ScriptSig: data,
	}
	txout := TXOutput{
		Value:        subsidy,
		ScriptPubKey: to,
	}
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{txout}}
	tx.SetID()
	return &tx
}

// NewUTXOTransaction 创建一个新交易
func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var inputs []TXInput
	var outpusts []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}
	// 这里并未按需取未使用输出，而是将所有未使用输出作为当前交易的输入
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}

		for _, out := range outs {
			input := TXInput{
				Txid:      txID,
				Vout:      out,
				ScriptSig: from,
			}
			inputs = append(inputs, input)
		}
	}
	// 一个地址的未使用的多个输出最终演变为一个未使用输出，提高了效率
	outpusts = append(
		outpusts,
		TXOutput{
			Value:        amount,
			ScriptPubKey: to,
		},
	)
	if acc > amount {
		outpusts = append(
			outpusts,
			TXOutput{
				Value:        acc - amount,
				ScriptPubKey: from,
			},
		)
	}
	tx := Transaction{
		ID:   nil,
		Vin:  inputs,
		Vout: outpusts,
	}
	tx.SetID()
	return &tx
}
