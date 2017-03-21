package conf

import (
	"net"
	"time"

	"gopkg.in/redis.v5"
)

////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////

var Name = "admin"
var key = "000000"
var API_SEND_SERVER = "http://119.29.142.168:81/main/int/"
var Add_name = "add"
var Del_name = "del"
var Weight = []float64{0, 0}                                 //用于筛选有效指令求库交集使用设置权重
var Min_times = "0"                                          //用于watchdog中的心跳筛选，最小值。
var flag = "0"                                               //用于在watch——dog进行清零动作的时候防止setip写入
var sec int = 120                                            //方便在配置页面进行设置做的转换
var Ticker = time.NewTicker(time.Duration(sec) * 1000000000) //用于设置watch——dog的检测时间间隔

var num = 0
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
var Device_id_length = 16
var Mask_heartbeat = "*"
var Mask_data = "#"
var Mask_reboot = "$"
var UDP_PORT_RECV_heart = 54321

var Socket_recv_heart, Err_recv_heart = net.ListenUDP("udp4", &net.UDPAddr{
	IP:   net.IPv4(0, 0, 0, 0),
	Port: UDP_PORT_RECV_heart,
})
var Client4 = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "xiaodepe", //  password set
	DB:       6,          // use default DB
})
