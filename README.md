#Crystalline
###by-小德培
Crystalline IOT-server is the server such as a routing which can manage IOT devices and can control devices with simple Http command in actice in any network environment ! you can make up your own IOT-platform！ sorry i haven't write the device code yet! if you r interested in ,you can join in ! 
 有想加入的小伙伴吗？
 快联系我！
 现在代码正在重构，，，，因为写在一起了，可读性太差不能看                                    

本平台可将以UDP/TCP为基础通讯的硬件设备进行群体调回（远程控制以及数据采集），支持高并发！开发者仅需要在通讯层按照本协议进行软件设计即可在地球上任何一个能接入互联网的地方对设备进行远程控制，不需要做任何网络/路由适配转发等等。并且对设备的控制仅需http协议即可（太TM无脑了）。



v1.6版本更改了通讯协议，即增加了randkey这一设备回应追踪机制，可以追踪所请求的硬件是否做出回应，以及当设备全部回应时通知请求方去取数据。


##功能以及原理
-------



- ####功能：本平台的主要功能是①连接并监管远端设备②根据请求下发相应指令给相应远端设备③接受远端设备返回的数据并返回给数据请求方。

- ####原理：本平台使用go语言构造，具备并支持高并发基础，使用高速TCP/IP协议与硬件进行通讯，使用通用http协议接收请求，并配合高速memcache型数据库redis做数据以及状态缓存,从而达到对多设备的远程集群调取。




##通讯协议
------
###硬件-server协议：
- ####device：

         UDP4心跳包：|   *  (单一符号)   | deviceID （5位str）|，外部移动网络发送间隔不得大于2min，端口暂定为54321
         
         udp4数据包：|   #  (单一符号)   | deviceID（5位str） | randkey（19位str）| data（不得大于65535 str）

- ####server：

         发送udp4数据：|deviceID (5位str) | randkey（19位str）| cmd （指令格式待定）
    	（后期加入ack返回包）







###数据库-server协议：
- ####redis：

		redis：自带驱动tcp

###请求方-server协议：


- ####http指令发送：
		请求数据：
                    post:xxx.xxx.xx.xxx:8001 （后期版本改为url）
                          
		携带参数：
                    body：admin  （简单验证机制，后续版本改善）	
				    key:    1 （这里是在redis——db0中注册编组的，想要控制那一组就写那一组）
		设备添加




- ####http数据返回：
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
				
				

- ####校验：
	`CRC32`：

	***注：进行crc运算的部分仅是json包里Data数据（字符串格式）***
			多项式：x³²+ x³¹+ x²⁴+ x²²+ x¹⁶+ x¹⁴+ x⁸+ x⁷+ x⁵+ x³+ x¹+ x⁰
			得到crcvalue中的值后进行crc32运算，然后将一同接收到的crcvalue进行比对（应先将crcvalue转为16进制），如果数值一致则数据完整性正常。




- ####设备添加：
        地址：
                    post:xxx.xxx.xx.xxx:8001/add

        携带参数：
                    body:add
                    
                    body value: 设备ID+分组
                            例：abcd2（将abcd分到第2中）
                            
        数据返回：
                    "add ok"


- ####设备删除
        地址：
                    post:xxx.xxx.xx.xxx:8001/del

        携带参数：
                    body:del
                    
                    body value: 设备ID
                            例：abcd
                            
        数据返回：
                    "del"













*此平台正在开发阶段，使用者只需根据通讯协议进行硬件开发，即可接入平台授权使用。*


##平台部署
-----
###数据库：
#####***redis（对接入设备心跳解析后的address进行记录和devices的授权分组，高并发必备）***

- ####分表：（数据仅供参考，格式才是精髓）
    
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
				    	——   row   |       value      |   score
				    	      1	        deviceID+data    randkey
					
				    		......
**不用再使用数据库进行data数据传输，同时mysql的godriver有问题，有超时连接的bug，高速数据写入以后不要用。**

	
##更新信息：
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


##版本改进计划：
    ①加入ACK信号，方便硬件确认数据完整性
    ②改为可视化控制台，自由配置通讯数据格式
    ③进一步优化性能，提高并发量
    
##原作者：小德培
##Wechat:Deathkingdom
##QQ:990834049@qq.com
