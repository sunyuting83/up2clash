# 盒装Clash局域网透明代理
---
##### clash in docker 进行路由转发实现全局透明代理
---
### 食用方法
#### 使用有线网卡Docker创建网络
1. 有线网卡开启混杂模式
```
sudo ip link set eth0 promisc on
```

2. 有线网卡docker创建网络,注意将网段改为你自己的
```
sudo docker network create -d macvlan --subnet=192.168.0.0/24 --gateway=192.168.0.1 -o parent=eth0 macnet
```

#### 使用无线网卡Docker创建网络
1. 无线网卡开启混杂模式
```
sudo ip link set wlp4s0 promisc on
```
2. 创建IPVlan虚拟网卡，重启失效，可以增加到/etc/rc.local达到开机启动目的。,注意将IP改为你自己的
```
sudo ip link add ipvlan link wlp4s0 type ipvlan mode l2
sudo ifconfig ipvlan up
sudo ip addr add dev ipvlan 192.168.0.118
```
3. 无线网卡docker创建网络,注意将网段改为你自己的
```
sudo docker network create -d ipvlan \
--subnet=192.168.0.0/24 --gateway=192.168.0.1 \
-o parent=wlp4s0 \
-o macvlan_mode=bridge \
ipnet
```

#### 新建resolv.conf文件
- 文件内容
```
nameserver 127.0.0.1
#把127.0.0.1改成你的网关地址
```

#### 运行容器
> 将resolv.conf修改成你的网关ip
- 食用有线网卡
```
sudo docker run -it --name clash_tp -d -v /your/path/resolv.conf:/etc/resolv.conf --sysctl net.ipv4.ip_forward=1 --network macnet --ip 192.168.0.119 --mac-address 00:50:56:00:60:42 --privileged --restart=always v2rss/clash_transparent_proxy
```
- 食用无线网卡
```
sudo docker run -it --name clash_tp -d --sysctl net.ipv4.ip_forward=1 --network ipnet --ip 192.168.0.119 -v /your/path/resolv.conf:/etc/resolv.conf --privileged --restart=always v2rss/clash_transparent_prox
```

#### 容器参数说明
| 参数  | 说明 | 传入值 |
| ------------ | ------------ | ------------ |
| -v | 将resolv.conf修改成你的网关ip | 新建的resolv.conf文件位置 |
| --sysctl | 开启流量转发 | 必须 |
| --ip | 修改成你需要固定的IP地址 | 同一网段不被占用的IP |

#### 局域网设备设置
将路由/手机/电脑等客户端 网关设置为容器ip,如192.168.0.119 ,dns也设置成这个

> 如果路由器不支持设置网关，请期待下一版本加入dnsmasq提供局域网DHCP服务。路由关闭DHCP由容器负责DHCP服务。

#### 更新Clash订阅
局域网直接访问容器IP:7893
```
http://192.168.0.119:7893
```

#### 订阅更新服务可用参数
| 参数  | 说明 |
| ------------ | ------------ |
| geturl | clash在线转换服务地址 |
| suburl | 你的订阅地址 |

> 例： http://192.168.0.119:7893/?suburl=https://yoursub.app

#### 构建方法
下载对应平台最新的Clash内核 Country.mmdb config.yaml
> 其他平台请自行修改Dockerfile和各种.sh中的文件名称

打包对应平台app.go 名称： up2c
```
#arm平台例：
CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o up2c app.go
```

全部放置在Dockerfile同一目录运行：
```
sudo docker build -t v2rss/clash_transparent_proxy .
```

#### 未来计划
- 加入dnsmasq实现局域网DHCP服务和去广告服务
- 手动更新clash内核和Country.mmdb
- 手动更新dnsmasq去广告规则