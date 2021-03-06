package db

import (
	"bytes"
	//	"time"

	//	. "platform_ota/alert"
	. "platform_ota/conf"
	. "platform_ota/login"
	. "platform_ota/tools"

	"fmt"
	"net"
	"net/http"
	//	"net/url"
	"strconv"
	//	"strings"

	"gopkg.in/redis.v5"
)

func Setdownip(data string, address *net.UDPAddr) {

	ipaddr := AnalyzeMessage([]byte(address.String()), len(address.String()))
	ipinfo := string(ipaddr[0]) + ":" + string(ipaddr[1])

	//旧代码
	//	data = Substr(data, 0, Device_id_length+Devicecode_length)
	//新代码
	devicecode := Substr(data, 0, Devicecode_length)
	device_id_length := Substr(data, Devicecode_length, 1)
	int_length, _ := strconv.ParseInt(device_id_length, 16, 0)

	id := Substr(data, Devicecode_length+2, int(int_length)*2)
	data = devicecode + id
	//新代码
	Client1.Set(data, ipinfo, 0).Err()
	Client0.ZIncrBy("1", 1, data) //在db0的“1”中开一个全注册设备表，并对接收到心跳的进行标记+1s
	return
}

func Writedown_data(data string) {
	//旧代码
	//	id := Substr(data, 0, Device_id_length+Devicecode_length) //分离出deviceid
	//新代码
	devicecode := Substr(data, 0, Devicecode_length)
	id_length := Substr(data, Devicecode_length, 1)

	int_id_length, _ := strconv.ParseInt(id_length, 16, 0)
	int_id_length = (int_id_length) * 2
	id := Substr(data, Devicecode_length+2, int(int_id_length))
	id = devicecode + id
	//新代码
	//旧代码
	//	randkey := Substr(data, Devicecode_length+Device_id_length, 18) //分离出randkey
	//新代码
	randkey_length := Substr(data, Devicecode_length+2+int(int_id_length), 1)

	randkey_length_int, _ := strconv.ParseInt(randkey_length, 16, 0)
	randkey_length_int = (randkey_length_int) * 2
	randkey := Substr(data, Devicecode_length+2+int(int_id_length)+2, int(randkey_length_int))
	//新代码
	fmt.Println(randkey)
	if randkey != Transmit_randkey {
		groupnum := Client0.ZScore("0", randkey).Val() //获取randkey对应的设备数量

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

	} else if randkey == Transmit_randkey {
		//		Tempnum = Tempnum + 1
		go Transmit(data)

	}

}

func senddata(randkeyfloat float64, groupnum float64, randkey string) {
	//	data_format := []byte{}
	//	result := []byte{}
	//	var API_SEND_SERVER string

	count := Client3.ZCount("0", randkey, randkey).Val()
	data := Client3.ZRangeByScore("0", redis.ZRangeBy{randkey, randkey, 0, count}).Val()
	//datastr := strings.Join(data, "") //此处取消分割

	//	for _, i := range data {

	////		result, API_SEND_SERVER = Code_format(i)
	////		a := [][]byte{result, data_format}
	////		data_format = bytes.Join(a, []byte(""))
	//	}

	pack, url := Code_format(data, randkey)
	fmt.Println("cccccccccccccccccc")
	body := bytes.NewBuffer(pack)
	resp, err := http.Post(url, "value", body)

	//resp, err := http.PostForm(API_SEND_SERVER, url.Values{"data": {pack}})
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
		fmt.Fprintf(w, "no_user") //用户名密码错误
	} else if "pass" == flag {
		id := r.PostFormValue("id")
		if 0 == Client0.ZScore("2", id).Val() {
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
			fmt.Fprintf(w, "1") //设备添加成功
		} else {
			fmt.Fprintf(w, "0") //该设备已注册
		}
	}
	return
}

func Delete_device(w http.ResponseWriter, r *http.Request) {

	var flag string
	r.ParseForm() //解析参数
	name := r.PostFormValue("name")
	key := r.PostFormValue("key")
	flag = Getuser(name, key)
	if "no user" == flag {
		fmt.Fprintf(w, "no_user") //用户名密码错误
	} else if "pass" == flag {
		id := r.PostFormValue("id")
		//		fmt.Println(id)
		//		a := Client0.ZRank("2", id).Val()
		//		fmt.Println(a)
		//		if 0 != a {
		if 0 != Client0.ZScore("2", id).Val() {
			err_del := Client0.ZRem("2", id).Err()
			if err_del != nil {
				fmt.Println("del failed:", err_del)
			}
			err_del2 := Client0.ZRem("1", id).Err()
			if err_del2 != nil {
				fmt.Println("del2 failed:", err_del2)
			}
			fmt.Fprintf(w, "1") //删除成功

		} else {
			fmt.Fprintf(w, "0") //该设备不存在
		}

	}
	return
}

func Transmit(data string) {
	//	fmt.Println("Transmit", Num)
	//	var API_AUTO_SERVER string
	device_code := Substr(data, 0, Devicecode_length)
	data_format := []string{data}
	//	fmt.Println(Getwebalive_flag)
	website_state := Client3.HGet("website_state", device_code).Val()
	if website_state == "online" {
		//	fmt.Println("123123133")

		result, url_ := Code_format(data_format, Transmit_randkey)
		//		fmt.Println(url)
		//		a := [][]byte{result, data_format}
		//		data_format = bytes.Join(a, []byte(""))
		//		pack := Code_json_data(data_format, Transmit_randkey)

		//		fmt.Println(string(result))
		body := bytes.NewBuffer(result)
		resp, err := http.Post(url_, "Auto", body)

		if err != nil {
			//			Alert_to_weichat("Transmit 发送失败")
			//			Client3.ZAdd("log", redis.Z{1, time.Now().String()})
			//			Num = Num + 1
			fmt.Println("Auto_post failed", err)
			fmt.Println("save")
			Client3.ZAdd(device_code, redis.Z{1, data})
		} else {
			//			resp.Close = true
			resp.Body.Close()
		}

	} else {
		fmt.Println("save")
		Client3.ZAdd(device_code, redis.Z{1, data})
	}
}

//消息转发功能
//func Transmit(Tempnum int, data string) {
//	fmt.Println("Transmit")
//	Client3.ZAdd("Transmit_flag", redis.Z{float64(1), float64(Tempnum)})
//	Client3.ZAdd("Transmit", redis.Z{float64(Tempnum), data})
//	//利用数据唯一码（Tempnum）标识已发送数据和因连接中断未能即使发送的数据
//	if Getwebalive(API_AUTO_SERVER) == "online" {
//		count := Client3.ZCount("Transmit_flag", "1", "1").Val()
//		list := Client3.ZRangeByScore("Transmit_flag", redis.ZRangeBy{"1", "1", 0, count}).Val()
//		for _, i := range list {
//			dat := Client3.ZRangeByScore("Transmit", redis.ZRangeBy{i, i, 0, 1}).Val()
//			//下面历遍dat仅仅是为了把数组转为string
//			for _, k := range dat {
//				crcvalue := Crccal(k)
//				packeddata := Code_json(k, Transmit_randkey, crcvalue)
//				pack := string(packeddata)
//				resp, err := http.PostForm(API_AUTO_SERVER, url.Values{"data": {pack}})
//				if err != nil {
//					fmt.Println("Auto_post failed", err)
//				} else {
//					//发送成功后将Transmit_flag中的flagscore置0表示可以drop
//					Client3.ZAdd("Transmit_flag", redis.Z{float64(0), i})
//					resp.Body.Close()
//				}
//			}
//		}
//		//将Transmit_flag自身置零的以及Transmit中的已发送数据一并drop
//		count_0 := Client3.ZCount("Transmit_flag", "0", "0").Val()
//		list_0 := Client3.ZRangeByScore("Transmit_flag", redis.ZRangeBy{"0", "0", 0, count_0}).Val()
//		Client3.ZRemRangeByScore("Transmit_flag", "0", "0")
//		for _, i := range list_0 {
//			Client3.ZRemRangeByScore("Transmit", i, i)
//		}

//	} else {
//		return
//	}

//}

//func Ota(data string, address *net.UDPAddr) {
//	fmt.Println("ota", data)

//	uid_code := Substr_last(data, 0, 4)
//	code := Substr(uid_code, 0, 2)
//	uid := Substr_last(uid_code, 0, 2)
//	fmt.Println("uid", uid)
//	fmt.Println("code", code)

//	Send_ota(code, uid, address)

//}

//func Send_ota(code string, uid string, address *net.UDPAddr) {
//	var name string
//	now_num := ota_pack_time[uid]

//	switch code {

//	case "00":
//		name = "正确"
//		Socket_recv_heart.WriteToUDP(ota_pack[now_num+1], address)
//		ota_pack_time[uid] = now_num + 1

//	case "01":
//		name = "完成更新"

//	case "02":
//		name = "app指定错误"

//	case "03":
//		name = "程序中指定地址错误"

//	case "04":
//		name = "没有更新指令"

//	case "05":
//		name = "长度错误"

//	case "06":
//		name = "未知错误"

//	case "07":
//		name = "命令错误"
//	case "08":
//		name = "校验码错误"

//	case "09":
//		name = "包序号错误"
//		Socket_recv_heart.WriteToUDP(ota_pack[now_num-1], address)
//		ota_pack_time[uid] = now_num - 1

//	case "0a":
//		name = "app地址重复，无法更新"

//	}

//	fmt.Println("更新进度:", "设备:", uid, "状态：", name)

//}
