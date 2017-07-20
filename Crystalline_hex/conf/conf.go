package conf

import (
	"fmt"
	//	"fmt"
	"net"
	"time"

	"github.com/astaxie/beego/config"
	"gopkg.in/redis.v5"
)

////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////

//var Tempnum int = 0

//var Getwebalive_flag = "online"

//var API_AUTO_SERVER string
//var API_SEND_SERVER string
var Url_1 string
var Url_1_auto string
var Url_2 string
var Url_2_auto string
var Alert_ip string

//var API_AUTO_SERVER = "http://www.crystoneiot.com:81/main/auto_receive/"
//var API_SEND_SERVER = "http://www.crystoneiot.com:81/main/int/"

//GIS备用服务器版本使用真实IP地址
//var API_AUTO_SERVER = "http://112.74.179.11:81/main/auto_receive/"
//var API_SEND_SERVER = "http://112.74.179.11:81/main/int/"

//var API_AUTO_SERVER = "http://119.29.142.168:81/main/auto_receive/"
//var API_SEND_SERVER = "http://119.29.142.168:81/main/int/"
var Add_name string
var Del_name string
var Devicecode1 string
var Devicecode2 string

var Weight = []float64{0, 0} //用于筛选有效指令求库交集使用设置权重
var Min_times = "0"          //用于watchdog中的心跳筛选，最小值。
var flag = "0"               //用于在watch——dog进行清零动作的时候防止setip写入
var Sec int                  //方便在配置页面进行设置做的转换
//var Ticker = time.NewTicker(time.Duration(Sec) * 1000000000) //用于设置watch——dog的检测时间间隔

var Ticker *time.Ticker

var Num = 0
var num2 = 0
var Client0 = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "xiaodepe", //  password set
	DB:       3,          // use default DB
})
var Client1 = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "xiaodepe", //  password set
	DB:       4,          // use default DB
})

//client1是服务器更新/存储设备address的分表

//var client2, err = sql.Open("mysql", "mdp:mdptest@tcp(45.124.65.33:3306)/mdptest")

var Client3 = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "xiaodepe", //  password set
	DB:       5,          // use default DB
})
var Device_id_length int
var Devicecode_length int
var Randkey_length int

var Mask_heartbeat string //"*"
var Mask_data string      //"#"
var Mask_reboot string    //"$"
var UDP_PORT_RECV_heart int

var Socket_recv_heart, Err_recv_heart = net.ListenUDP("udp4", &net.UDPAddr{
	IP:   net.IPv4(0, 0, 0, 0),
	Port: 54321,
})

var Socket_alert, Err_socket_alert = net.ListenUDP("udp4", &net.UDPAddr{
	IP: net.IPv4(0, 0, 0, 0),
	//	Port: 52100,
})

var Client4 = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "xiaodepe", //  password set
	DB:       6,          // use default DB
})

var Transmit_randkey string
var Getwebalive_flag string

func init() {

	configdata, err := config.NewConfig("ini", "config.conf")

	if err != nil {
		fmt.Println(err)
	}

	Transmit_randkey = configdata.String("Transmit_randkey") //触发自动转发的randkey

	Getwebalive_flag = configdata.String("Getwebalive_flag") //判断网站后端是否在线

	//	API_AUTO_SERVER = configdata.String("API_AUTO_SERVER") //自动转发的网站地址

	//	API_SEND_SERVER = configdata.String("API_SEND_SERVER") //主动获取数据的目标网站地址

	Add_name = configdata.String("Add_name") //添加设备的指令

	Del_name = configdata.String("Del_name") //删除设备的指令

	Sec, _ = configdata.Int("sec") //心跳维护时间间隔即在sec内收到至少一次心跳，则判定设备在线
	Ticker = time.NewTicker(time.Duration(Sec) * 1000000000)

	Randkey_length, _ = configdata.Int("Randkey_length") //randkey的长度

	Devicecode_length, _ = configdata.Int("Devicecode_length") //设备类别码的长度

	Device_id_length, _ = configdata.Int("Device_id_length") //设备id的长度

	Mask_heartbeat = configdata.String("Mask_heartbeat") //心跳前置掩码

	//	Mask_heartbeat = "2A"

	Devicecode1 = configdata.String("Devicecode1") //数据类别前置码1

	Devicecode2 = configdata.String("Devicecode2") //数据类别前置码1

	Mask_data = configdata.String("Mask_data") //数据前置掩码

	Url_1 = configdata.String("Url_1") //获取第一目标地址

	Url_1_auto = configdata.String("Url_1_auto") //获取第一目标的自动转发地址

	Url_2 = configdata.String("Url_2") //获取第二目标地址

	Url_2_auto = configdata.String("Url_2_auto") //获取第二目标的自动转发地址

	Alert_ip = configdata.String("Alert_ip")

	Mask_data = configdata.String("Mask_data") //数据前置掩码

	Mask_reboot = configdata.String("Mask_reboot") //此处用于设备的时间校准或者自身的相关响应操作

	UDP_PORT_RECV_heart, _ = configdata.Int("UDP_PORT_RECV_heart") //UDP与设备通讯的开放端口
	//	UDP_PORT_RECV_heart = 54321

	Client3.HSet("website", Url_1_auto, Devicecode1)
	Client3.HSet("website", Url_2_auto, Devicecode2)

}
