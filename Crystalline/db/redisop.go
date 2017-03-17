package db

import (
	. "Crystalline/conf"
	. "Crystalline/login"
	. "Crystalline/tools"

	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"gopkg.in/redis.v5"
)

func Setdownip(data string, address *net.UDPAddr) {

	ipaddr := AnalyzeMessage([]byte(address.String()), len(address.String()))
	ipinfo := string(ipaddr[0]) + ":" + string(ipaddr[1])
	data = Substr(data, 0, Device_id_length)
	Client1.Set(data, ipinfo, 0).Err()
	Client0.ZIncrBy("1", 1, data) //在db0的“1”中开一个全注册设备表，并对接收到心跳的进行标记+1s
	return
}

func Writedown_data(data string) {
	//	fmt.Println("writedown")
	id := Substr(data, 0, Device_id_length)        //分离出deviceid
	randkey := Substr(data, Device_id_length, 19)  //分离出randkey
	groupnum := Client0.ZScore("0", randkey).Val() //获取randkey对应的设备数量
	//	fmt.Println(groupnum)
	check := Client1.Exists(id)
	if check.Val() == true {
		//		go confirm(data)
		randkeyfloat, _ := strconv.ParseFloat(randkey, 64)
		fmt.Println(randkeyfloat)
		Client3.ZAdd("0", redis.Z{randkeyfloat, data})
		num := Client0.ZIncrBy("0", -1, randkey).Val() //收到一个设备的数据结果，则对其randkey减一/
		fmt.Println(num)
		if num == 0 {
			Client0.ZRem("0", randkey) //删除randkey
			//			fmt.Println(groupnum)//groupnum不能传至此处，此时已经改变，不再是原先设备数量，
			go senddata(randkeyfloat, groupnum, randkey)

		}
	}

	return
}

func senddata(randkeyfloat float64, groupnum float64, randkey string) {

	count := Client3.ZCount("0", randkey, randkey).Val()

	data := Client3.ZRangeByScore("0", redis.ZRangeBy{randkey, randkey, 0, count}).Val()
	datastr := strings.Join(data, " ")
	crcvalue := Crccal(datastr)

	packeddata := Code_json(datastr, randkey, crcvalue)
	pack := string(packeddata)
	resp, err := http.PostForm(API_SEND_SERVER,
		url.Values{"data": {pack}})
	if err != nil {
		fmt.Println("post failed", err)

	} else {
		resp.Body.Close()
	}

	/////////////////////////////////////
	err_rem := Client3.ZRemRangeByScore("0", randkey, randkey).Err()
	if err_rem != nil {
		fmt.Println("删除缓存failed", err_rem)
		//		client0.ZRem("0", randkey) //删除randkey
		return
	}

	fmt.Println("data sending over")

	return

}

func Add_device(w http.ResponseWriter, r *http.Request) {

	var flag string
	r.ParseForm() //解析参数
	name := r.PostFormValue("name")
	key := r.PostFormValue("key")
	flag = Getuser(name, key)
	if "no user" == flag {
		fmt.Fprintf(w, " err ")
	} else if "pass" == flag {
		id := r.PostFormValue("id")
		if 0 == Client0.ZRank("2", id).Val() {
			groupnum := r.PostFormValue("groupnum")
			groupnumfloat, _ := strconv.ParseFloat(groupnum, 64)
			err_add := Client0.ZAdd("1", redis.Z{0, id}).Err()
			if err_add != nil {
				fmt.Println("add failed:", err_add)
			}
			err_add2 := Client0.ZAdd("2", redis.Z{groupnumfloat, id}).Err()
			if err_add2 != nil {
				fmt.Println("add2 failed:", err_add2)
			}
			fmt.Fprintf(w, "设备增加成功 ")

		} else {
			fmt.Fprintf(w, "该设备已注册")
		}

	}

}
func Delete_device(w http.ResponseWriter, r *http.Request) {

	var flag string
	r.ParseForm() //解析参数
	name := r.PostFormValue("name")
	key := r.PostFormValue("key")
	flag = Getuser(name, key)
	if "no user" == flag {
		fmt.Fprintf(w, " err ")
	} else if "pass" == flag {
		id := r.PostFormValue("id")
		if 0 != Client0.ZRank("2", id).Val() {
			err_del := Client0.ZRem("2", id).Err()
			if err_del != nil {
				fmt.Println("del failed:", err_del)
			}
			err_del2 := Client0.ZRem("1", id).Err()
			if err_del2 != nil {
				fmt.Println("del2 failed:", err_del2)
			}
			fmt.Fprintf(w, " 已删除 ")

		} else {
			fmt.Fprintf(w, "未发现该设备")
		}

	}

}
