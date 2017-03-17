package tools

import (
	"encoding/json"

	"fmt"
	"hash/crc32"
)

type datatype struct {
	Randkey  string
	Data     string
	Crcvalue uint32
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

func Crccal(data string) (result uint32) {
	fmt.Println("crccal")
	var crc32key = crc32.MakeTable(0xD5828281)
	databyte := []byte(data)
	result1 := crc32.Checksum(databyte, crc32key)

	return result1

}

func Code_json(datastr string, randkey string, crcvalue uint32) (result []byte) {
	fmt.Println("codejson")
	pack := &datatype{
		Randkey:  randkey,
		Data:     datastr,
		Crcvalue: crcvalue,
	}

	result, err := json.Marshal(pack)
	if err != nil {
		fmt.Println("error:", err)
	}
	return result
}
