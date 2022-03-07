package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const (
	version            = byte(0x00)
	addressChecksumLen = 4
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey // 私钥
	PublicKey  []byte           // 公钥
}

// GetAddress 生成钱包地址
// 地址有version、pubkey、checksum构成
// version 1字节
// pubkeyhash  20字节
// checksum 长度取决于addressChecksumLen，默认4字节
func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)

	// 将1字节的版本号作为前缀加入
	versionedPayload := append([]byte{version}, pubKeyHash...)
	// 生成公钥校验和
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	return Base58Encode(fullPayload)
}

// Checksum 生成公钥的校验和
// 两次sha256加密，取前addressChecksumLen个字节作为校验和
func checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

// HashPubKey 生成公钥hash
// 两次(sh256->pripemd160)加密产生20字节的hash
func HashPubKey(pubKey []byte) []byte {
	// 生成32字节的hash
	publicSHA256 := sha256.Sum256(pubKey)
	RIPEMD160Hasher := ripemd160.New()
	// 生成20字节的hash
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		log.Panic(err)
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}

// ValidateAddress 检查地址是否有效
// 通过checksum校验地址
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

// newKeyPair 创建密钥对
// 公钥是椭圆上x1,x2,x3.....,y1,y2,y3.......的组合
func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	pubKey := append(private.PublicKey.X.Bytes(), private.Y.Bytes()...)
	return *private, pubKey
}

// NewWallet 创建一个钱包
func NewWallet() *Wallet {
	private, pubilc := newKeyPair()
	return &Wallet{
		PrivateKey: private,
		PublicKey:  pubilc,
	}
}
