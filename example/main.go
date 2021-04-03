package main

import (
	"fmt"

	"github.com/rfyiamcool/go-snowflake"
)

func method1() {
	id, err := snowflake.Next()
	if err != nil {
		panic(err)
	}
	ts := snowflake.GetTimeFromID(id)
	fmt.Println("id: ", id, " timestamp: ", ts, " (ms) ")
}

func method2() {
	workerID := 111
	snowflake.Init(int64(workerID))
	id, err := snowflake.Next()
	if err != nil {
		panic(err)
	}
	ts := snowflake.GetTimeFromID(id)
	fmt.Println("id: ", id, " timestamp: ", ts, " (ms) ")
}

func method3() {
	sf := snowflake.New(222)
	id, err := sf.Next()
	if err != nil {
		panic(err)
	}
	ts := sf.GetTimeFromID(id)
	fmt.Println("id: ", id, " timestamp: ", ts, " (ms) ")
}

func main() {
	method1()
	method2()
	method3()
}
