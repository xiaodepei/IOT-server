package main

import (
	. "Crystalline/control"
	. "Crystalline/db"
	. "Crystalline/udp"
	"fmt"
	"net/http"
)

func main() {

	fmt.Println("system online")
	go Udp_port()
	go Watch_dog() //监视设备状态
	Httpserver()

}
func Httpserver() {

	http.HandleFunc("/", Login) //设置访问的路由
	http.HandleFunc("/add", Add_device)
	http.HandleFunc("/del", Delete_device)
	http.ListenAndServe(":8001", nil) //设置监听的端口
	fmt.Println("系统启动失败，可能是http端口被占用")

}
