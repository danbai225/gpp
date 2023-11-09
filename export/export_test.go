package main

import (
	"syscall"
	"testing"
	"time"
)

func TestDll(t *testing.T) {
	ddDLL := syscall.NewLazyDLL("C:\\code\\gpp\\export\\bin\\windows\\gpp.dll")
	start := ddDLL.NewProc("Start")
	r1, _, _ := start.Call()
	if r1 != 0 {
		t.Error("启动失败")
	} else {
		t.Log("启动成功")
	}
	time.Sleep(10 * time.Second)
	stop := ddDLL.NewProc("Stop")
	r1, _, _ = stop.Call()
	if r1 != 0 {
		t.Error("停止失败")
	} else {
		t.Log("启动成功")
	}
	select {}
}
