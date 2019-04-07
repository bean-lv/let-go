package main

import (
	"fmt"
	"letgo"
)

func main() {
	testLog()
	// testConfig()
}

func testLog() {
	msg := "This is a test message."
	err := letgo.DefaultLogger()
	if err != nil {
		panic(err)
	}
	letgo.Logger.Error(msg)
	letgo.Logger.Fatal(msg)
	letgo.Logger.Crash(msg)
	letgo.Logger.Warn(msg)
	letgo.Logger.Info(msg)
	letgo.Logger.Notice(msg)
	letgo.Logger.Status(msg)
	letgo.Logger.Debug(msg)
	letgo.Logger.Trace(msg)
}

func testConfig() {
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
