package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type Person struct {
	Name string
	Age int
}

func main()  {
	XM := Person{
		Name: "Hello",
		Age:  20,
	}
	//使用gob进行序列化（编码）得到字节流
	var buffer bytes.Buffer
	//定义一个编码器
	encode := gob.NewEncoder(&buffer)
	//使用编码器进行编码
	err := encode.Encode(&XM)
	if err != nil {
		fmt.Println("编码失败:",err)
	}
	fmt.Println("编码后的XM:",buffer)

	var LL Person
	//定义一个解码器
	decode := gob.NewDecoder(&buffer)
	//使用解码器进行解码
	err = decode.Decode(&LL)
	if err != nil {
		fmt.Println("解码失败:",err)
	}
	fmt.Println("解码后的LL:",LL)
}