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

nohup /etc/up2c/up2c > /dev/null 2>&1 &
nohup /etc/clash/clash-linux-amd64 -d /etc/clash > /dev/null 2>&1
exec "$@"