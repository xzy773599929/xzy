package main

import (
	"fmt"
	"os"
)

func main()  {
	len1 := len(os.Args)
	fmt.Println(len1)
	for i,cmd := range os.Args {
		fmt.Printf("org[%d]:%s\n",i,cmd)
	}
}