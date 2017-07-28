package udp

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"

	. "Crystalline_hex/alert"
	. "Crystalline_hex/conf"
	. "Crystalline_hex/db"
	. "Crystalline_hex/tools"

	"gopkg.in/redis.v5"
)

func Udp_port() {

	if Err_recv_heart != nil {
		fmt.Println("监听失败", Err_recv_heart)
	}
	defer Socket_recv_heart.Close()

	for {
		data := make([]byte, 1024)
		read, address, err := Socket_recv_heart.ReadFromUDP(data)
		//		fmt.Println(read)
		if err != nil {
			fmt.Println("读取数据失败", err)
			continue
		}
		go selection(data[:read], address)
	}
}

func selection(data []byte, address *net.UDPAddr) {
	data_str := hex.EncodeToString(data)
	flag := Substr(data_str, 0, 2)
	//	fmt.Println(data_str)
	switch flag {
	case Mask_data:
		data_str = strings.TrimLeft(data_str, Mask_data) //用于统一端口接受数据类型前置码

		Writedown_data(data_str)
	case Mask_heartbeat:
		//		fmt.Println(data_str)
		data_str = strings.TrimLeft(data_str, Mask_heartbeat) //用于统一端口接受数据类型前置码
		//		fmt.Println(data_str)
		Setdownip(data_str, address)
	case Mask_reboot:
		data_str = strings.TrimLeft(data_str, Mask_reboot)
		Confirm(data_str) //用于设备验证连接
	}

}

func Confirm(data string) {
	fmt.Println(data)
	//	data = Substr(data, 8, Device_id_length)
	value := Client1.Get(data).Val()

	fmt.Println(value)
	udpip, _ := net.ResolveUDPAddr("udp4", value)
	senddata := Mask_reboot + data
	hex_senddata, _ := hex.DecodeString(senddata)
	_, Err_recv_heart = Socket_recv_heart.WriteToUDP(hex_senddata, udpip)
	if Err_recv_heart != nil {
		Alert_to_weichat("确认信号发送失败！")
		fmt.Println("确认信号发送失败!", Err_recv_heart)
		return
	}
	return

}

//首先去db0中查找经过注册设备的分组和ID号
func Findgroup(group string, randkey string) {
	fmt.Println("findgroup")

	data := Client0.ZRevRangeByScore("2", redis.ZRangeBy{group, group, 0, 0}).Val() //把设备表改为了在一个有序数列下valus——score关系
	offline := Client0.ZRevRangeByScore(randkey, redis.ZRangeBy{"0", "0", 0, 0}).Val()
	fmt.Println(data)

	for _, i := range offline {
		for num, address := range data {
			if i == address {
				data[num] = "nil"
			}
		}
	}

	for _, i := range data {
		if i == "nil" {
			continue
		}
		Sendcmd(Client1, i, randkey)
	}
	err := Client0.ZRemRangeByScore(randkey, "0", "0").Err() //用后删除randkey
	if err != nil {
		fmt.Println("delete failed:", err)
		return
	}
	return
}
func Sendcmd(client1 *redis.Client, address string, randkey string) {
	data := address + randkey
	senddata, _ := hex.DecodeString(data)
	value := client1.Get(address).Val()
	udpip, _ := net.ResolveUDPAddr("udp4", value)

	//	senddata := []byte(address + randkey)
	_, Err_recv_heart = Socket_recv_heart.WriteToUDP(senddata, udpip)
	if Err_recv_heart != nil {
		Alert_to_weichat("发送UDP失败（sendcmd）！")
		fmt.Println("发送数据失败!", Err_recv_heart)
		return
	}
	return
}
