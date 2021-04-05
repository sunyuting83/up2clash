#!/bin/bash

set -e


reset_iptables(){
    iptables -P INPUT ACCEPT
    iptables -P FORWARD ACCEPT
    iptables -P OUTPUT ACCEPT
    iptables -t nat -F
    iptables -t mangle -F
    iptables -F
    iptables -X
}

set_clash_iptables(){
    # 在 nat 表中创建新链
    iptables -t nat -N CLASHRULE

    iptables -t nat -A CLASHRULE -p tcp --dport 1905 -j RETURN

    iptables -t nat -A CLASHRULE -d 0.0.0.0/8 -j RETURN
    iptables -t nat -A CLASHRULE -d 10.0.0.0/8 -j RETURN
    iptables -t nat -A CLASHRULE -d 127.0.0.0/8 -j RETURN
    iptables -t nat -A CLASHRULE -d 169.254.0.0/16 -j RETURN
    iptables -t nat -A CLASHRULE -d 172.16.0.0/12 -j RETURN
    iptables -t nat -A CLASHRULE -d 192.168.0.0/16 -j RETURN
    iptables -t nat -A CLASHRULE -d 224.0.0.0/4 -j RETURN
    iptables -t nat -A CLASHRULE -d 240.0.0.0/4 -j RETURN
    iptables -t nat -A CLASHRULE -p tcp -j REDIRECT --to-ports 7892

    #拦截 dns 请求并且转发!
    iptables -t nat -A PREROUTING -p udp --dport 53 -j REDIRECT --to-ports 53
    iptables -t nat -A PREROUTING -p tcp --dport 53 -j REDIRECT --to-ports 53

    # 在 PREROUTING 链前插入 CLASHRULE 链,使其生效
    iptables -t nat -A PREROUTING -p tcp -j CLASHRULE
}

reset_iptables
set_clash_iptables

#开启转发# curl ip.sb
# sudo ip link set ovs_eth1 promisc on
# sudo ip link set ovs_eth1 promisc on docker
# sudo docker network create -d macvlan --subnet=192.168.2.0/24 --gateway=192.168.2.1 -o parent=ovs_eth1 macnet

# sudo docker run --name clash_tp -d -v /var/services/homes/admin/clash:/clash_config --sysctl net.ipv4.ip_forward=1 --network macnet --ip 192.168.2.118 --mac-address 00:50:56:00:60:42 --privileged zhangyi2018/clash_transparent_proxy:1.0.7


#echo "1" > /proc/sys/net/ipv4/ip_forward

/bin/sh /etc/clash/clash.sh start

exec "$@"