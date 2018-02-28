package udp

import (
	"bytes"
	//	"encoding/binary"
	"encoding/hex"
	"fmt"
	//	"io/ioutil"
	"net"
	. "platform_ota/alert"
	. "platform_ota/conf"
	. "platform_ota/db"
	. "platform_ota/ota"
	. "platform_ota/tools"

	"strings"
	"time"

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
	fmt.Println("get data")
	//	fmt.Println(data_str)
	switch flag {
	case Mask_data:
		//		fmt.Println(data_str)
		data_str = strings.TrimLeft(data_str, Mask_data) //用于统一端口接受数据类型前置码
		Writedown_data(data_str)
	case Mask_heartbeat:
		//		fmt.Println(data_str)
		data_str = strings.TrimLeft(data_str, Mask_heartbeat) //用于统一端口接受数据类型前置码
		//		fmt.Println(data_str)
		Setdownip(data_str, address)
	case Mask_reboot:
		data_str = strings.TrimLeft(data_str, Mask_reboot)
		Confirm(data_str, address) //用于设备验证连接

	case Mask_ota:
		data_str = strings.TrimLeft(data_str, Mask_ota)
		Ota(data_str, address)

	}
}

func Confirm(data string, ip *net.UDPAddr) {

	c := fmt.Sprintln(time.Now().UTC())

	fmt.Println("校准时间", c)

	senddata := Mask_reboot + data

	hex_senddata, _ := hex.DecodeString(senddata)

	fmt.Println(hex_senddata)

	result := [][]byte{hex_senddata, []byte(c)}
	d := bytes.Join(result, []byte(""))

	_, Err_recv_heart = Socket_recv_heart.WriteToUDP(d, ip)

	//	_, Err_recv_heart = Socket_recv_heart.WriteToUDP(hex_senddata, udpip)
	if Err_recv_heart != nil {
		Alert_to_weichat("确认信号发送失败！")
		fmt.Println("确认信号发送失败!", Err_recv_heart)
		return
	}
	return

}

//首先去db0中查找经过注册设备的分组和ID号
func Findgroup(group string, randkey string, app string) {
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
		Sendcmd(Client1, i, randkey, app)
	}
	err := Client0.ZRemRangeByScore(randkey, "0", "0").Err() //用后删除randkey
	if err != nil {
		fmt.Println("delete failed:", err)
		return
	}
	return
}

//func Sendcmd(client1 *redis.Client, address string, randkey string, app_str string) {
//	ota_pack_time = make(map[string]int)
//	var app byte
//	fmt.Println("dizhi", address)

//	uid := Substr_last(address, 0, 4)

//	ota_pack_time[uid] = 0

//	var name string
//	switch app_str {

//	case "1":
//		app = byte(1)
//		name = "APP1.bin"

//	case "2":
//		app = byte(2)
//		name = "APP2.bin"

//	}

//	ota_pack = make(map[int][]byte)

//	b, err := ioutil.ReadFile(name)
//	if err != nil {
//		fmt.Print(err)
//	}
//	var buf bytes.Buffer
//	len_b, _ := buf.Write(b)
//	fmt.Println(len_b)
//	times := len_b / 512
//	if (len_b % 512) > 0 {
//		times = times + 1
//	}

//	fmt.Println("总包数", times)
//	pack_num = times
//	cmd_pack := []byte{0x23, 0x2a, 0x31, 0x00, 0x03, app, byte(times - 1)}
//	var crc_cmd byte
//	for _, num := range cmd_pack {
//		crc_cmd = crc_cmd + num
//	}
//	crc_cmd_byte := []byte{crc_cmd}
//	cmdpack_pre := [][]byte{cmd_pack, crc_cmd_byte}
//	cmdpack := bytes.Join(cmdpack_pre, []byte(""))
//	ota_pack[0] = cmdpack
//	for a := 0; a < times; a++ {
//		var crc byte
//		var length = []byte{0x00, 0x00}
//		data := buf.Next(512)
//		binary.BigEndian.PutUint16(length, 512+2)
//		//fmt.Println(length)
//		datapack := []byte{0x23, 0x2a, 0x32}
//		number := byte(a + 1)
//		number_byte := []byte{number}
//		datapack_pre := [][]byte{datapack, length}
//		datapack = bytes.Join(datapack_pre, []byte(""))
//		datapack_pre = [][]byte{datapack, number_byte}
//		datapack = bytes.Join(datapack_pre, []byte(""))
//		datapack_pre = [][]byte{datapack, data}
//		datapack = bytes.Join(datapack_pre, []byte(""))
//		for _, num := range datapack {
//			crc = crc + num
//		}
//		crc_byte := []byte{crc}
//		datapack_pre = [][]byte{datapack, crc_byte}
//		datapack = bytes.Join(datapack_pre, []byte(""))
//		//fmt.Println(datapack)
//		ota_pack[a+1] = datapack
//	}

//	value := client1.Get(address).Val()
//	udpip, _ := net.ResolveUDPAddr("udp4", value)
//	_, Err_recv_heart = Socket_recv_heart.WriteToUDP(ota_pack[0], udpip)

//	if Err_recv_heart != nil {

//		fmt.Println("发送数据失败!", Err_recv_heart)
//		return
//	}
//	return
//}
