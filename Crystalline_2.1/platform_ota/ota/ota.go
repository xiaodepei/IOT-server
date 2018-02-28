package ota

//ota使用

import (
	"bytes"
	"encoding/binary"
	//	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	. "platform_ota/tools"

	"gopkg.in/redis.v5"

	. "platform_ota/conf"
)

var ota_pack map[int][]byte

var ota_pack_time map[string]int

var pack_num int

func Sendcmd(client1 *redis.Client, address string, randkey string, app_str string) {
	ota_pack_time = make(map[string]int)
	var app byte
	fmt.Println("dizhi", address)

	uid := Substr_last(address, 0, 4)

	ota_pack_time[uid] = 0

	var name string
	switch app_str {

	case "1":
		app = byte(1)
		name = "APP1.bin"

	case "2":
		app = byte(2)
		name = "APP2.bin"

	}

	ota_pack = make(map[int][]byte)

	b, err := ioutil.ReadFile(name)
	if err != nil {
		fmt.Print(err)
	}
	var buf bytes.Buffer
	len_b, _ := buf.Write(b)
	fmt.Println(len_b)
	times := len_b / 512
	if (len_b % 512) > 0 {
		times = times + 1
	}

	fmt.Println("总包数", times)
	pack_num = times
	cmd_pack := []byte{0x23, 0x2a, 0x31, 0x00, 0x03, app, byte(times - 1)}
	var crc_cmd byte
	for _, num := range cmd_pack {
		crc_cmd = crc_cmd + num
	}
	crc_cmd_byte := []byte{crc_cmd}
	cmdpack_pre := [][]byte{cmd_pack, crc_cmd_byte}
	cmdpack := bytes.Join(cmdpack_pre, []byte(""))
	ota_pack[0] = cmdpack
	for a := 0; a < times; a++ {
		var crc byte
		var length = []byte{0x00, 0x00}
		data := buf.Next(512)
		binary.BigEndian.PutUint16(length, 512+2)
		//fmt.Println(length)
		datapack := []byte{0x23, 0x2a, 0x32}
		number := byte(a + 1)
		number_byte := []byte{number}
		datapack_pre := [][]byte{datapack, length}
		datapack = bytes.Join(datapack_pre, []byte(""))
		datapack_pre = [][]byte{datapack, number_byte}
		datapack = bytes.Join(datapack_pre, []byte(""))
		datapack_pre = [][]byte{datapack, data}
		datapack = bytes.Join(datapack_pre, []byte(""))
		for _, num := range datapack {
			crc = crc + num
		}
		crc_byte := []byte{crc}
		datapack_pre = [][]byte{datapack, crc_byte}
		datapack = bytes.Join(datapack_pre, []byte(""))
		//fmt.Println(datapack)
		ota_pack[a+1] = datapack
	}

	value := client1.Get(address).Val()
	udpip, _ := net.ResolveUDPAddr("udp4", value)
	_, Err_recv_heart = Socket_recv_heart.WriteToUDP(ota_pack[0], udpip)

	if Err_recv_heart != nil {

		fmt.Println("发送数据失败!", Err_recv_heart)
		return
	}
	return
}

func Ota(data string, address *net.UDPAddr) {
	fmt.Println("ota", data)

	uid_code := Substr_last(data, 0, 4)
	code := Substr(uid_code, 0, 2)
	uid := Substr_last(uid_code, 0, 2)
	fmt.Println("uid", uid)
	fmt.Println("code", code)

	Send_ota(code, uid, address)

}

func Send_ota(code string, uid string, address *net.UDPAddr) {
	var name string
	now_num := ota_pack_time[uid]

	switch code {

	case "00":
		name = "正确"
		Socket_recv_heart.WriteToUDP(ota_pack[now_num+1], address)
		ota_pack_time[uid] = now_num + 1

	case "01":
		name = "完成更新"

	case "02":
		name = "app指定错误"

	case "03":
		name = "程序中指定地址错误"

	case "04":
		name = "没有更新指令"

	case "05":
		name = "长度错误"

	case "06":
		name = "未知错误"

	case "07":
		name = "命令错误"
	case "08":
		name = "校验码错误"

	case "09":
		name = "包序号错误"
		Socket_recv_heart.WriteToUDP(ota_pack[now_num-1], address)
		ota_pack_time[uid] = now_num - 1

	case "0a":
		name = "app地址重复，无法更新"

	}
	a := fmt.Sprintf("%.2f", (float64(now_num)/float64(pack_num))*100)
	fmt.Println("更新进度:", "设备:", uid, "状态：", name, "进度:", a)

}
