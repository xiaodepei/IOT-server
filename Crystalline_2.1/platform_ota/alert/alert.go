package alert

import (
	"fmt"
	"net"
	. "platform_ota/conf"
	. "platform_ota/tools"
)

func Alert_to_weichat(data string) {

	enaes_data := EnAES(data)
	udpip, _ := net.ResolveUDPAddr("udp4", Alert_ip)
	//	senddata := []byte(address + randkey)
	_, Err_socket_alert = Socket_alert.WriteToUDP([]byte(enaes_data), udpip)
	fmt.Println("send alert message")
	if Err_socket_alert != nil {
		fmt.Println("发送数据失败!", Err_socket_alert)
		return
	}
}
