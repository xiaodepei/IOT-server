package login

import (
	. "Crystalline/conf"
	"fmt"
	"net/http"
)

///////////////////////////////////////////////////
//验证是否存在该用户
//存在该用户则返回pass
//不存在则返回no user
//////////////////////////////////////////////////
func Getuser(name string, pwd string) (state string) {

	b := Client4.Get(name).Val()
	if b == "" {
		state = "no user"
	} else if b == pwd {
		state = "pass"
	}
	return state

}

///////////////////////////////////////////////////
//验证服务器是否处于联通状态
//联通返回：online
//出现问题：offline
func Getwebalive(URL string) (state string) {
	resp, err := http.Get(URL)
	if err != nil {
		fmt.Println("WEB端服务离线：", err)
		state = "offline"
	} else {
		resp.Body.Close()
		state = "online"
	}
	return state
}
