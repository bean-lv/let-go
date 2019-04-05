package main

import (
	"fmt"
	"letgo"
)

func main() {
	filename := "conf/conf.conf"
	err := letgo.DefaultConfig(filename)
	if err != nil {
		panic(err)
	}
	fmt.Println(letgo.Config)
	fmt.Println(letgo.Config.Get("url"))
	letgo.Config.Set("test", "mytest init conf ")
	fmt.Println(letgo.Config.Get("test"))
	fmt.Println(letgo.Config.Get("中文"))
}
