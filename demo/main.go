package main

import (
	"Block/bolt"
	"fmt"
	"log"
)

func main() {
	//1.打开数据库
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		log.Panic("打开数据库失败")
	}
	defer db.Close()
	//将要操作数据库
	//创建表，写操作
	db.Update(func(tx *bolt.Tx) error {
		//2.找到抽屉bucket（如果没有就创建）
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil {
			//没有抽屉需要创建
			bucket, err = tx.CreateBucket([]byte("b1"))
			if err != nil {
				log.Panic("创建bucket(b1)失败！")
			}
		}
		bucket.Put([]byte("111"), []byte("hello"))
		bucket.Put([]byte("222"), []byte("world"))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	//查询表操作
	db.View(func(tx *bolt.Tx) error {
		//2.找到抽屉bucket（如果没有panic）
		bucket := tx.Bucket([]byte("b1"))
		if bucket == nil {
			log.Panic("The bucket b1 is not existed")
		}
		v1 := bucket.Get([]byte("111"))
		v2 := bucket.Get([]byte("222"))
		fmt.Println(string(v1), string(v2))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}
