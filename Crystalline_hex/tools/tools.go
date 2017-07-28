package tools

import (
	//	"encoding/hex"
	. "Crystalline_hex/conf"
	//	"bytes"
	"bytes"
	"encoding/json"
	"fmt"
	//	"hash/crc32"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"strconv"
	"strings"
)

//type datatype struct {
//	Randkey  string
//	Data     string
//	Crcvalue uint32
//}

type Alert struct {
	Item []Alert_data
}

type Alert_data struct {
	Time string
	Info string
}

type Item struct {
	Item []Datatype
}
type Item2 struct {
	Item []Datatype2
}
type Datatype struct {
	Gateway_id string
	Randkey    string
	Data       []Format1
}

type Datatype2 struct {
	Gateway_id string
	Randkey    string
	Data       []Format2
}

type Format1 struct {
	Id        string
	Tempature int64
	Power     int64
}

type Format2 struct {
	Tag_id               string
	Tempature_device     int64
	Tempature_enviroment int64
	Power                int64
	Time                 string
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
		data_format = Format_2(data, randkey)
		if randkey != Transmit_randkey {
			url = Url_1
		} else {
			url = Url_1_auto
		}

	case Devicecode2:
		data_format = Format_2(data, randkey)
		if randkey != Transmit_randkey {
			url = Url_2
		} else {
			url = Url_2_auto
		}

	}

	return data_format, url

}

//下面是网关接收到的子设备的设备ID 以及数据
//format1是耳标测温的数据
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

//format——2是新版本gis(正版)
//func Format_2(data []string, randkey string) (result []byte) {
//	var I Item2
//	for _, data_str := range data {
//		var D Datatype2
//		gateway_id := Substr(data_str, Devicecode_length, Device_id_length)
//		data_tag := Substr(data_str, Devicecode_length+Randkey_length+Device_id_length, len(data_str)-Devicecode_length+Randkey_length-Device_id_length)
//		for a := 0; a < len(data_tag)/36; a++ {
//			data1 := Substr(data_tag, 36*a, 36)
//			id := Substr(data1, 0, 16)
//			tempature_device := Substr(data1, 16, 4)
//			tempature_device_int := Tempature_cover(tempature_device)
//			tempature_enviroment := Substr(data1, 20, 4)
//			tempature_enviroment_int := Tempature_cover(tempature_enviroment)
//			power := Substr(data1, 24, 4)
//			power_format := Power_cover(power)
//			time_str := Substr(data1, 28, 8)
//			D.Data = append(D.Data, Format2{Tag_id: id, Tempature_device: tempature_device_int, Tempature_enviroment: tempature_enviroment_int, Power: power_format, Time: time_str})
//		}
//		I.Item = append(I.Item, Datatype2{Gateway_id: gateway_id, Randkey: randkey, Data: D.Data})
//	}

//	result, err := json.Marshal(I)
//	if err != nil {
//		fmt.Println("error:", err)
//	}
//	return result

//}

//format——2是新版本gis（临时版本）
func Format_2(data []string, randkey string) (result []byte) {
	var I Item2
	for _, data_str := range data {
		var D Datatype2
		gateway_id := Substr(data_str, Devicecode_length, Device_id_length)
		data_tag := Substr(data_str, Devicecode_length+Randkey_length+Device_id_length, len(data_str)-Devicecode_length+Randkey_length-Device_id_length)
		for a := 0; a < len(data_tag)/36; a++ {
			data1 := Substr(data_tag, 36*a, 36)
			id := Substr(data1, 4, 16)
			tempature_device := Substr(data1, 20, 4)
			tempature_device_int := Tempature_cover(tempature_device)
			tempature_enviroment := Substr(data1, 24, 4)
			tempature_enviroment_int := Tempature_cover(tempature_enviroment)
			power := Substr(data1, 0, 4)
			power_format := Power_cover(power)
			time_str := Substr(data1, 28, 8)
			D.Data = append(D.Data, Format2{Tag_id: id, Tempature_device: tempature_device_int, Tempature_enviroment: tempature_enviroment_int, Power: power_format, Time: time_str})
		}
		I.Item = append(I.Item, Datatype2{Gateway_id: gateway_id, Randkey: randkey, Data: D.Data})
	}

	result, err := json.Marshal(I)
	if err != nil {
		fmt.Println("error:", err)
	}
	return result

}

func Show_url(device_id string) (url string, deviceid string) {

	devicecode := Substr(device_id, 0, Devicecode_length)
	deviceid = Substr(device_id, Devicecode_length, Device_id_length)
	switch devicecode {

	case Devicecode1:
		url = Url_1
	case Devicecode2:
		url = Url_2
	}

	return url, deviceid

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

///////////////////////////////////////////////////////////////////////////
//AES加密//
///////////////////////////////////////////////////////////////////////////

var key = []byte("1234567890123456")

//////////////////
//加密密钥
//////////////////

func EnAES(data string) (ened string) {

	result, err := AesEncrypt([]byte(data), key)
	if err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(result)

}

func DeAES(data string) (deed string) {

	origData, err := AesDecrypt([]byte(data), key)
	if err != nil {
		panic(err)
	}
	return string(origData)

}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
