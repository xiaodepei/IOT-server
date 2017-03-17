package udp

import (
	. "Crystalline/conf"
	. "Crystalline/db"
	. "Crystalline/tools"
	"fmt"
	"net"
	"strings"

	"gopkg.in/redis.v5"
)

func Udp_port() {

	if Err_recv_heart != nil {
		fmt.Println("监听失败", Err_recv_heart)
	}
	defer Socket_recv_heart.Close()

	for {
		data := make([]byte, 2048)
		read, address, err := Socket_recv_heart.ReadFromUDP(data)

		if err != nil {
			fmt.Println("读取数据失败", err)
			continue
		}
		go selection(string(data[:read]), address)

	}

}

func selection(data string, address *net.UDPAddr) {
	flag := Substr(data, 0, 1)
	switch flag {
	case Mask_data:
		data = strings.TrimLeft(data, "#") //用于统一端口接受数据类型前置码

		Writedown_data(data)
	case Mask_heartbeat:
		data = strings.TrimLeft(data, "*") //用于统一端口接受数据类型前置码

		Setdownip(data, address)
	case Mask_reboot:
		data = strings.TrimLeft(data, "$")
		Confirm(data) //用于设备验证连接
	}

}
func Confirm(data string) {

	data = Substr(data, 0, Device_id_length)

	value := Client1.Get(data).Val()
	udpip, _ := net.ResolveUDPAddr("udp4", value)
	senddata := []byte("$" + data)
	_, Err_recv_heart = Socket_recv_heart.WriteToUDP(senddata, udpip)
	if Err_recv_heart != nil {
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

	value := client1.Get(address).Val()
	udpip, _ := net.ResolveUDPAddr("udp4", value)
	senddata := []byte(address + randkey)
	_, Err_recv_heart = Socket_recv_heart.WriteToUDP(senddata, udpip)

	if Err_recv_heart != nil {
		fmt.Println("发送数据失败!", Err_recv_heart)
		return
	}
	return
}
