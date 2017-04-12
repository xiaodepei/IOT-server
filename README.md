

 Crystalline
---

本平台可将以UDP/TCP为基础通讯的硬件设备进行群体调回（远程控制以及数据采集），支持高并发！开发者仅需要在通讯层按照本协议进行软件设计即可在地球上任何一个能接入互联网的地方对设备进行远程控制，不需要做任何网络/路由适配转发等等。并且对设备的控制仅需http协议即可（太TM无脑了）。



v1.6版本更改了通讯协议，即增加了randkey这一设备回应追踪机制，可以追踪所请求的硬件是否做出回应，以及当设备全部回应时通知请求方去取数据。


功能以及原理
-------



- 功能：本平台的主要功能是①连接并监管远端设备②根据请求下发相应指令给相应远端设备③接受远端设备返回的数据并返回给数据请求方。

- 原理：本平台使用go语言构造，具备并支持高并发基础，使用高速TCP/IP协议与硬件进行通讯，使用通用http协议接收请求，并配合高速memcache型数据库redis做数据以及状态缓存,从而达到对多设备的远程集群调取。




通讯协议
------
硬件-server协议：

- device：

         UDP4心跳包：|   *  (单一符号)   | deviceID | 外部移动网络发送间隔不得大于2min，端口暂定为54321
         
         UDP4数据包：|   #  (单一符号)   | deviceID | randkey（19位str）| data（不得大于65535 str）
         
         UDP4ACK包： |   $  (单一符号)   | deviceID | 

- server：

         发送udp4数据：|deviceID  | randkey（19位str）| cmd （指令格式待定）
        
         返回UDP4ACK包：| $ (单一符号)  | deviceID 

















数据库-server协议：

- redis：

    	redis：自带驱动tcp

请求方-server协议：


- http指令发送：
		请求数据：
                    post：XXX.XXX.XXX.XXX:8001 （后期版本改为url）
                          
		携带参数：
                    name:   XXXXX（用户名需要在redis中注册）	
				    key:    XXXXX（用户名对应密码）
                    groupnum: XX (硬件的逻辑分组，需要注册)
		设备添加




- http数据返回：
		数据返回post（指定服务器非调用者）：
			JSON格式：
				结构体：type datatype struct {
						Randkey  string
						Data     string
						Crcvalue uint32
					}
					
				json包：pack := &datatype{
						Randkey:  randkey,
						Data:     datastr,
						Crcvalue: crcvalue,
					}
				
				

- 校验：
	`CRC32`：

	***注：进行crc运算的部分仅是json包里Data数据（字符串格式）***
			多项式：x³²+ x³¹+ x²⁴+ x²²+ x¹⁶+ x¹⁴+ x⁸+ x⁷+ x⁵+ x³+ x¹+ x⁰
			得到crcvalue中的值后进行crc32运算，然后将一同接收到的crcvalue进行比对（应先将crcvalue转为16进制），如果数值一致则数据完整性正常。




- 设备添加：
        地址：
                    post:XXX.XXX.XXX.XXX:8001/add

        携带参数：
                    name:         xiaodepei(注册用户名)
                    key:          XXXXX(密码)
                    id:           XXXXXX（设备的唯一识别码或者就是上面的deviceid）
                    groupnum:     1(分组)
                                      
        数据返回：
                    "add ok"


- 设备删除
        地址：
                    post:XXX.XXX.XXX.XXX:8001/del

        携带参数：
                    name:         xiaodepei(注册用户名)
                    key:          123456(密码)
                    id:           XXXXXX（设备的唯一识别码或者就是上面的deviceid）                         
       
       数据返回：
                    "del ok"













*此平台正在开发阶段，使用者只需根据通讯协议进行硬件开发，即可接入平台授权使用。*


平台部署
-----
数据库：
***redis（对接入设备心跳解析后的address进行记录和devices的授权分组，高并发必备）***

- 分表：（数据仅供参考，格式才是精髓）
    
			xiaodepei（库）
				——db0（用于设备授权ID以及分组，分组用于集群远控调用，需要手动建表，有序集合）
					——2（分组）（此表数据应同步一份给调用者，或动态关联，或作为从机）
						—— row  |  value  |  score
						    1	   12345       0
						    2      54321       0
						......
					
					——1（分组）此表为所有已注册设备列表，主要作为设备统计以及watch——dog检测设备在线状态使用
						—— row  |  value  |  score
						    1	   12345       0
						    2      54321       0
                            
						score在这里作为计数watch——dog监控周期中接受对应deviceID的心跳包数量
						value作为deviceID

					——0（分组）此分组是保留分组，作为randkey的缓存分组，不使用（数据不保留）
						——  row   |   value   |   score
						   1         randkey	  num 
                                       
						randkey为生成的19位随机str，每个请求对应唯一，请求周期结束后，程序自动删除，
						num为计数判定，其实就是上面分组里面的设备数，当每个设备都把数据返回后，此值将自减为0同时自动删除。
                        
                    ——Offline 将根据设置的心跳监测间隔将离线的设备丢到里面(不需要创建，没有离线设备时，不会有该表)
                    
                        ——row     |   value     | score 
                            1         12345         0

				    ——db1（用于server自动生成并更新从device解析的心跳数据，主要是设备address，不需要手动建立，字符串key-value）内部数据保留但随时根据根据IP更改
				    	——  key   |   value
					       12345    192.168.0.1:8080
					      54321    127.0.0.1:8000

			    ——db2（用于做接收采集数据的缓冲池，不需要手动建立，有序集合）内部数据不保留
				    ——0分组（用于主动采集数据缓冲）
                        ——   row   |       value      |   score
				    	      1	        deviceID+data    randkey
					——Transmit分组（用于设备主动发送数据的缓存）
                    
			            ——  row   |       value      |   score
                            1            data           Tempnum                        
                        
                    ——Transmit_flag分组（用于对硬件主动发送的数据做状态标记，未发送给服务器的数据标记为1,已发送的数据标记为0等待系统清除）
                        ——  row   |       value      |   score
                            1            Tempnum           1或2
                              
                ——db3（用于平台用户注册）
                        ——  key   |   value  
                            name      psd
                            


	
更新信息：
    接着更：

	继续更：12-06：今天添加了watchdog来监测设备状态，及时通知服务器设备离线情况。

	12-26：
		①更改了randkey生成机制，19位，测试结果100万次做到不重复。
		②增加算法，保证离线数据不会收到指令而产生报错
		③与硬件的通讯端口减为一个口（54321）

	12-30： 
		①server-2.3加入了confirm的硬件软重启/时间校正指令        
    	1-11：
		①完善文档
	
	2-21:
		新增模块——设备添加与删除模块
	3-18:
		修改了设备增加与删除部分的权限验证机制
        
        
    4-1:已加入debug模块，独立于主程序外


版本改进计划：
    ①加入ACK信号，方便硬件确认数据完整性，（已完成）
    ②改为可视化控制台，自由配置通讯数据格式
    ③进一步优化性能，提高并发量
    ④创建转发fifo，对接收到的数据自动进行转发，先入先出原则，（已完成）
    ⑤指令扩展，增加平台对设备的指令多样化。


by-小德培
Wechat:Deathkingdom
QQ:990834049@qq.com
















