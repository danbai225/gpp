package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
*/
import "C"
import "fmt"

import (
	"github.com/danbai225/gpp/core"
	box "github.com/sagernet/sing-box"
)

var b *box.Box

//export Start
func Start() {
	if b != nil {
		fmt.Println("已启动")
	}
	var err error
	b, err = core.Client()
	if err != nil {
		fmt.Println(err)
	}
	err = b.Start()
	if err != nil {
		fmt.Println(err)
	}
}

//export Stop
func Stop() {
	if b == nil {
		fmt.Println("未启动")
	}
	err := b.Close()
	if err != nil {
		fmt.Println(err)
	}
}

//export Test
func Test() {
	fmt.Println("测试成功")
}

func main() {}
