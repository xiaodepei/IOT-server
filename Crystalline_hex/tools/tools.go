package tools

import (
	//	"encoding/hex"
	. "Crystalline_hex/conf"
	//	"bytes"
	"encoding/json"
	"fmt"
	//	"hash/crc32"
	"strconv"
	"strings"
)

//type datatype struct {
//	Randkey  string
//	Data     string
//	Crcvalue uint32
//}

type Item struct {
	Item []Datatype
}

type Datatype struct {
	Gateway_id string
	Randkey    string
	Data       []Format1
}

type Format1 struct {
	Id        string
	Tempature int64
	Power     int64
}

func AnalyzeMessage(buff []byte, len int) []string {
	analMsg := make([]string, 0)
	strNow := ""
	for i := 0; i < len; i++ {
		if string(buff[i:i+1]) == ":" {
			analMsg = append(analMsg, strNow)
			strNow = ""
		} else {
			strNow += string(buff[i : i+1])
		}
	}
	analMsg = append(analMsg, strNow)
	return analMsg
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

//func Crccal(data string) (result uint32) {
//	//	fmt.Println("crccal")
//	var crc32key = crc32.MakeTable(0xD5828281)
//	databyte := []byte(data)
//	result1 := crc32.Checksum(databyte, crc32key)

//	return result1

//}

//func Code_json_data(data []byte, randkey string) (result []byte) {
//	//	fmt.Println("codejson")
//	pack := &datatype{
//		Randkey: randkey,
//		Data:    data,
//	}

//	result, err := json.Marshal(pack)
//	if err != nil {
//		fmt.Println("error:", err)
//	}
//	return result
//}

//下面是网关发送数据包的内容
func Code_format(data []string, randkey string) (data_format []byte, url string) {
	datastr := strings.Join(data, "")

	devicecode := Substr(datastr, 0, Devicecode_length)
	//	device_id := Substr(data, Devicecode_length, Device_id_length)
	//	randkey := Substr(data, Device_id_length+Devicecode_length, Randkey_length)
	//data = Substr(datastr, Device_id_length+Devicecode_length+Randkey_length, len(data)-Device_id_length+Devicecode_length+Randkey_length)
	switch devicecode {
	case Devicecode1:
		data_format = Format_1(data, randkey)
		if randkey != Transmit_randkey {
			url = Url_1
		} else {
			url = Url_1_auto
		}
	}

	return data_format, url

}

//下面是网关接收到的子设备的设备ID 以及数据
func Format_1(data []string, randkey string) (result []byte) {
	var I Item
	for _, data_str := range data {
		var D Datatype
		gateway_id := Substr(data_str, Devicecode_length, Device_id_length)
		data_tag := Substr(data_str, Devicecode_length+Randkey_length+Device_id_length, len(data_str)-Devicecode_length+Randkey_length-Device_id_length)
		for a := 0; a < len(data_tag)/24; a++ {
			data1 := Substr(data_tag, 24*a, 24)
			id := Substr(data1, 0, 8)
			tempature := Substr(data1, 8, 4)
			tempature_int := Tempature_cover(tempature)
			power := Substr(data1, 12, 4)
			power_format := Power_cover(power)
			D.Data = append(D.Data, Format1{Id: id, Tempature: tempature_int, Power: power_format})

		}

		I.Item = append(I.Item, Datatype{Gateway_id: gateway_id, Randkey: randkey, Data: D.Data})

	}

	result, err := json.Marshal(I)
	if err != nil {
		fmt.Println("error:", err)
	}
	return result

}

func Tempature_cover(tempature string) (result int64) {
	temp, _ := strconv.ParseInt(tempature, 16, 16)
	result = temp / 10
	return result
}

func Power_cover(power string) (result int64) {
	power_format, _ := strconv.ParseInt(power, 16, 16)
	return power_format
}
