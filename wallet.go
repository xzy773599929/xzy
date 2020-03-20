package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
	"github.com/iocn-io/ripemd160"
	"log"
)

//定义钱包结构体，每一个钱包保存了私钥公钥对
type Wallet struct {
	//私钥
	PrivateKey *ecdsa.PrivateKey
	//公钥，这里不是保存原始公钥，而是公钥结构体中的X和Y拼接的字符串，在校验时重新拆分（参考r，s传递)
	PubKey []byte
}

//创建钱包函数
func NewWallet() *Wallet {
	//创建曲线
	curve := elliptic.P256()
	//生成私钥
	privateKey ,err := ecdsa.GenerateKey(curve,rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	//生成公钥
	pubKeyOrig := privateKey.PublicKey
	pubKey := append(pubKeyOrig.X.Bytes(),pubKeyOrig.Y.Bytes()...)
	//创建钱包
	wallet := Wallet{
		PrivateKey: privateKey,
		PubKey:     pubKey,
	}
	return &wallet
}

//生成地址
func (w *Wallet)NewAddress()string {
	pubKey := w.PubKey
	version := byte(00)
	//a.对公钥进行哈希处理：RIPEMD160（sha256())
	ripemdHash := HashPubkey(pubKey)
	payload := append([]byte{version},ripemdHash[:]...)

	//b.获取校验码:checksum()
	checkCode := checksum(payload)

	//c.拼接:version + hash + checksum
	pubKeyHash := append(payload,checkCode...)

	//d.base58,该方法来自go语言的比特币库
	address := base58.Encode(pubKeyHash)
	return address
}

//公钥哈希函数
func HashPubkey(pubKey []byte)[]byte  {
	hash256 := sha256.Sum256(pubKey)
	ripemd160Hasher := ripemd160.New()
	_,err := ripemd160Hasher.Write(hash256[:])
	if err != nil {
		log.Panic(err)
	}
	ripemdHash := ripemd160Hasher.Sum(nil)
	return ripemdHash
}
//checksum
func checksum(payload []byte)[]byte  {
	hashFirst := sha256.Sum256(payload)
	hashSecond := sha256.Sum256(hashFirst[:])
	checkCode := hashSecond[:4]
	return checkCode
}