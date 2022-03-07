# Block
Block 是基于go实现的一个简易区块链

## Part 1 基本原型
实现了区块链的基本原型，区块及链

## Part 2 工作量证明
Hashcash工作量证明算法,分为以下步骤
- 采取一些公开的数据`data`
- 向其添加一个从0开始的计数器`counter`
- 获取组合(`data+counter`)的`hash`
- 检查`hash`是否满足某些要求
    - 满足，即完成
    - 不满足,增加计数器，并重复步骤3和4

## Part 3 持久化
![持久化](./static/持久化.svg)

## Part 4 交易

Transactions:
 - Coinbase 创世块交易，只有输出没有输入
 - UTXOTransaction 点对点的交易, 由一个或多个未使用输出作为输入产生一个或两个输出

Balance:
- UTXO模型
- 余额通过遍历整个交易记录得来

## Part 5 地址与钱包

### 地址

### 密钥对

### 签名

### 算法
- ecdsa 椭圆曲线数字签名算法
- sha256、ripemd160 加密算法
- base58 编码算法

# 
```go
// wallet
// NewWallet 创建一个钱包
func NewWallet() *Wallet {
	private, pubilc := newKeyPair()
	return &Wallet{
		PrivateKey: private,
		PublicKey:  pubilc,
	}
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

// sign
// privKey wallet.privatekey
// txCopy.id transaction.id
r,s,err := ecdsa.Sign(rand.Reader,&privKey,txCopy.ID)
if err!= nil {
	log.Panic(err)
}
signature := append(r.Bytes(),s.Bytes()...)
tx.Vin[inID].Signature = signature

// verify
// pubkey wallet.pubilc
r := big.Int{}
s := big.Int{}
sigLen := len(vin.Signature)
r.SetBytes(vin.Signature[:(sigLen / 2)])
s.SetBytes(vin.Signature[(sigLen / 2):])

x := big.Int{}
y := big.Int{}
keyLen := len(vin.PubKey)
x.SetBytes(vin.PubKey[:(keyLen / 2)])
y.SetBytes(vin.PubKey[(keyLen / 2):])

rawPubKey := ecdsa.PublicKey{curve, &x, &y}
if ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) == false {
 return false
}

````