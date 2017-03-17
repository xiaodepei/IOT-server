package login

import (
	. "Crystalline/conf"
	//	"fmt"
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
