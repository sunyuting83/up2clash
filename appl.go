package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

func main() {
	var (
		geturl  string
		suburl  string
		sub     string
		url     string
		command string = "/etc/config/sh/clash.sh restart"
		local   bool
		exclude string = "%e7%be%8e%e5%9b%bd%7c%e7%88%b1%e6%b2%99%e5%b0%bc%e4%ba%9a%7c%e7%bd%97%e9%a9%ac%e5%b0%bc%e4%ba%9a%7c%e5%8a%a0%e6%8b%bf%e5%a4%a7%7c%e5%8c%97%e7%be%8e%7c%e4%bf%84%e7%bd%97%e6%96%af%7c%e5%85%8b%e7%bd%97%e5%9c%b0%e4%ba%9a%7c%e4%b9%8c%e5%85%8b%e5%85%b0%7c%e6%ac%a7%e7%9b%9f%7c%e6%91%a9%e5%b0%94%e5%a4%9a%e7%93%a6%7cBuLink%7c%e5%9f%83%e5%8f%8a%7c%e8%91%a1%e8%90%84%e7%89%99%7c%e6%8b%89%e8%84%b1%e7%bb%b4%e4%ba%9a%7c%e8%8b%b1%e5%9b%bd%7c%e8%a5%bf%e7%8f%ad%e7%89%99%7c%e7%91%9e%e5%85%b8%7c%e5%8d%a2%e6%a3%ae%e5%a0%a1%7c%e5%be%b7%e5%9b%bd%7c%e7%88%b1%e5%b0%94%e5%85%b0%7c%e6%b3%a2%e5%85%b0%7c%e6%8d%b7%e5%85%8b%7c%e6%84%8f%e5%a4%a7%e5%88%a9%7cisx.yt%7c%e8%8d%b7%e5%85%b0%7cNL%7cUS%7cDE%7cFR%7cAU%7cRU%7cRO%7cZZ%7cGB%7c%e5%8d%b0%e5%ba%a6%7c%e4%b8%b9%e9%ba%a6"
		emoji   string = "true"
		list    string = "false"
		udp     string = "true"
		tfo     string = "false"
		scv     string = "false"
		fdn     string = "false"
		sort    string = "false"
	)
	flag.StringVar(&geturl, "geturl", "http://127.0.0.1:25500", "请求网址")
	flag.StringVar(&sub, "sub", "clash", "请求类型")
	flag.StringVar(&suburl, "suburl", "http%3a%2f%2f127.0.0.1%3a5550%2f%3fw%3d1%26i%3d2", "订阅网址")
	flag.BoolVar(&local, "local", false, "本地代理")
	flag.Parse()
	url = strings.Join([]string{geturl, "/sub?target=", sub, "&url=", suburl, "&exclude=", exclude, "&emoji=", emoji, "&list=", list, "&udp=", udp, "&tfo=", tfo, "&scv=", scv, "&fdn=", fdn, "&sort=", sort}, "")
	res, errers := http.Get(url)
	if errers != nil {
		fmt.Println("get error")
		return
	}
	if res != nil {
		defer res.Body.Close()
		if res.StatusCode != 200 {
			// log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
			fmt.Println("status error")
			return
		}
		s, _ := ioutil.ReadAll(res.Body)
		if local {
			ioutil.WriteFile("/etc/config/conf/clash.yaml", s, 0644)
		} else {
			tran, err := ioutil.ReadFile("./transparent")
			if err != nil {
				fmt.Println("read error")
				return
			}
			// fmt.Println(tran)
			d := string(s)
			l := len(d)
			e := strings.Index(d, "proxies:")
			t := string(tran)
			d = d[e : l-1]
			dd := strings.Join([]string{t, d}, "\n")
			var strtobyte []byte = []byte(dd)
			ioutil.WriteFile("/etc/config/conf/clash.yaml", strtobyte, 0644)
		}
		// fmt.Println(command)
		cmd := exec.Command("/bin/bash", "-c", command)
		if err := cmd.Run(); err != nil {
			fmt.Println("run error")
			return
		}
		return
	}
}
