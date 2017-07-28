package control

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	//	"strings"
	"time"

	. "Crystalline_hex/alert"
	. "Crystalline_hex/conf"
	. "Crystalline_hex/login"
	. "Crystalline_hex/tools"
	. "Crystalline_hex/udp"

	"gopkg.in/redis.v5"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var flag string
	var mesage string
	var groupnum int64 //用来
	r.ParseForm()      //解析参数
	name := r.PostFormValue("name")
	key := r.PostFormValue("key")
	flag = Getuser(name, key)
	if "no user" == flag {
		fmt.Fprintf(w, " err ")
	} else if "pass" == flag {
		mesage = r.PostFormValue("groupnum")
		a := time.Now().UnixNano()
		randkey := strconv.FormatInt(a, 10)
		randkey = Substr(randkey, 1, 18)
		member := Client0.ZRangeByScore("2", redis.ZRangeBy{mesage, mesage, 0, 0}).Val() //拿出指令中需要的全部设备ID，后续版本加入同一条指令查询多组的功能

		for num, i := range member { //num是通过数组下标来计数member的，因为从0开始所以需要+1，i是每个deviceid
			Client0.ZAdd(randkey, redis.Z{0, i})
			groupnum = int64(num)
		}
		group_with_offline_num := Client0.ZInterStore(randkey, redis.ZStore{Weight, "sum"}, randkey, "Offline").Val() //将求交集之后的可用device储存在randkey命名的表中，且同时得到可用设备数量
		groupnum = groupnum + 1 - group_with_offline_num                                                              //因为是下标，所以需要+1，然后再减去离线的设备
		Client0.ZAdd("0", redis.Z{float64(groupnum), randkey})                                                        //将randkey为srt写入db0中，记录对应randkey的次数，接受函数需要使用到
		go Findgroup(mesage, randkey)                                                                                 //放在这里比放在上面快多了，groupnum其实和message一样，就是查询一下当前分组有多少个设备
		fmt.Fprintf(w, randkey)                                                                                       //输出到客户端

	}

}

func Watch_dog() {
	for _ = range Ticker.C {
		Alert_to_weichat("开始")
		Client0.ZRemRangeByScore("Offline", "0", "0") //开始之前首先删除上次的离线设备数据
		result := Client0.ZRevRangeByScore("1", redis.ZRangeBy{Min_times, Min_times, 0, 0}).Val()

		for _, i := range result {
			Client0.ZAdd("Offline", redis.Z{0, i}) //将离线设备进行记录
		}

		for _, i := range result {
			//resultstr := strings.Join(result, " ")
			uRl, deviceid := Show_url(i)
			//			fmt.Println(i)
			resp, err := http.PostForm(uRl, url.Values{"offline": {deviceid}})
			//resp, err := http.Post(url, "offlinedevice", deviceid)
			if err != nil {
				fmt.Println("离线数据发送失败", err)
				Alert_to_weichat("离线数据发送失败")
			} else {
				resp.Body.Close()
			}
		}
		fmt.Println("11111")
		devicenumber := Client0.ZCard("1").Val()
		num := Client0.ZRange("1", 0, devicenumber).Val()
		//		fmt.Println(num)
		for _, mumber := range num {
			Client0.ZAdd("1", redis.Z{0, mumber})

		}
		url_group := Client3.HKeys("website").Val()
		fmt.Println(url_group)
		for _, i := range url_group {
			fmt.Println("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
			Getwebalive_flag := Getwebalive(i) //获取远程端口状态
			fmt.Println(Getwebalive_flag)
			code := Client3.HGet("website", i).Val()
			fmt.Println("33333")
			fmt.Println(code)
			time.Sleep(10000)
			fmt.Println("44444")
			err := Client3.HSet("website_state", code, Getwebalive_flag).Err()
			fmt.Println("55555")
			if err != nil {
				fmt.Println("ahahahahahahahah", err)
			}
			if Getwebalive_flag == "online" {
				fmt.Println("66666")
				predata := Client3.ZCard(code).Val()
				fmt.Println("77777")
				fmt.Println(predata)
				if predata > 0 {
					go Send_store(predata, code)
				}
			} else {
				fmt.Println("888888")
				Alert_to_weichat(i + "offline")
				fmt.Println(i, "offline")
				fmt.Println("999999")
			}
			fmt.Println("10101010")

		}

		fmt.Printf("ticked at %v ", time.Now())
	}
}

func Send_store(datanum int64, device_code string) {
	var a int64
	//fmt.Println(999999 + datanum)
	for a = 0; a < datanum; a++ {
		fmt.Println("1")
		list := Client3.ZRangeByScore(device_code, redis.ZRangeBy{"1", "1", 0, 1}).Val()
		//		fmt.Println(list)
		for _, i := range list {
			fmt.Println("2")
			result := Client3.ZRem(device_code, i).Val()
			fmt.Println(i)
			if result == 1 {
				fmt.Println("send store!!!!!")
				fmt.Println("3")
				pack, url := Code_format(list, Transmit_randkey)
				fmt.Println("4")
				//				fmt.Println(string(pack))
				body := bytes.NewBuffer(pack)
				resp, err := http.Post(url, "send_store", body)

				if err != nil {
					fmt.Println("post failed", err)
				} else {
					resp.Body.Close()
				}

			}
		}
	}

}
