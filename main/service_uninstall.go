package main

import (
	"github.com/TeaWeb/code/teautils"
	"log"
)

// 卸载服务
func main() {
	log.Println("uninstalling ...")
	manager := teautils.NewServiceManager("TeaWeb", "TeaWeb Server")
	err := manager.Uninstall()
	if err != nil {
		log.Println("ERROR: " + err.Error())
		manager.Close()
		manager.PauseWindow()
		return
	}

	log.Println("uninstalled service successfully")
	log.Println("done.")
	
	manager.PauseWindow()
}
