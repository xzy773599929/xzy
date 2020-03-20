package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"math/big"
)

func main()  {
	//创建曲线
	curve := elliptic.P256()
	//生成私钥
	privateKey ,err := ecdsa.GenerateKey(curve,rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	//生成公钥
	pubKey := privateKey.PublicKey

	data := "hello world !"
	hash := sha256.Sum256([]byte(data))

	//签名
	r,s,err := ecdsa.Sign(rand.Reader,privateKey,hash[:])
	if err != nil {
		log.Panic(err)
	}

	//把r，s进行序列化传输
	signature := append(r.Bytes(),s.Bytes()...)

	//定义两个辅助的big.Int
	r1 := big.Int{}
	s1 := big.Int{}

	//拆分signature，平均分成两半，前半部分为r，后半部分为s
	r1.SetBytes(signature[0:len(signature)/2])
	s1.SetBytes(signature[len(signature)/2:])

	//校验需要3个东西，数据，签名，公钥
	res := ecdsa.Verify(&pubKey,hash[:],&r1,&s1)
	fmt.Printf("校验结果:%v",res)
}