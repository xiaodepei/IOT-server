package login

import (
	//	"fmt"
	"net/http"
	//	"net/url"

	. "Crystalline_hex/conf"
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
//func Getwebalive(URL string) (state string) {
//	resp, err := http.Get(URL)
//	if err != nil {
//		fmt.Println("WEBserver——offline")
//		state = "offline"
//	} else {
//		resp.Body.Close()
//		state = "online"
//	}
//	return state
//}

//func Getwebalive(URL string) (state string) {
//	resp, _ := http.Get(URL)
//	a := resp.StatusCode

//	//	defer resp.Body.Close()
//	if a == 200 {
//		state = "online"
//		resp.Body.Close()
//	} else {
//		fmt.Println("WEBserver——offline")
//		state = "offline"
//	}
//	return state
//}

//func Getwebalive(URL_ string) (state string) {
//	resp, err := http.Get(URL_)
//	fmt.Println(resp)
//	if err != nil {
//		fmt.Println(err)
//		state = "offline"

//	} else {
//		a := resp.StatusCode
//		if a == 200 {
//			state = "online"
//			resp.Body.Close()
//		}

//	}
//	return state
//}

func Getwebalive(URL_ string) (state string) {
	resp, _ := http.Get(URL_)
	a := resp.StatusCode
	if a == 200 {
		state = "online"
		resp.Body.Close()
	} else {
		state = "offline"
		resp.Body.Close()
	}
	return state
}
