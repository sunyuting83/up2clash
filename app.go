package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	var (
		port string
	)
	flag.StringVar(&port, "p", "7893", "Default Port 7893")
	flag.Parse()
	app := gin.New()
	app.GET("/", func(c *gin.Context) {
		var geturl string = c.DefaultQuery("geturl", "https://pub-api-1.bianyuan.xyz")
		var suburl string = c.DefaultQuery("suburl", "https://jiang.netlify.app")
		suburl = url.QueryEscape(suburl)
		var exclude string = c.DefaultQuery("exclude", "%e7%be%8e%e5%9b%bd%7c%e7%88%b1%e6%b2%99%e5%b0%bc%e4%ba%9a%7c%e7%bd%97%e9%a9%ac%e5%b0%bc%e4%ba%9a%7c%e5%8a%a0%e6%8b%bf%e5%a4%a7%7c%e5%8c%97%e7%be%8e%7c%e4%bf%84%e7%bd%97%e6%96%af%7c%e5%85%8b%e7%bd%97%e5%9c%b0%e4%ba%9a%7c%e4%b9%8c%e5%85%8b%e5%85%b0%7c%e6%ac%a7%e7%9b%9f%7c%e6%91%a9%e5%b0%94%e5%a4%9a%e7%93%a6%7cBuLink%7c%e5%9f%83%e5%8f%8a%7c%e8%91%a1%e8%90%84%e7%89%99%7c%e6%8b%89%e8%84%b1%e7%bb%b4%e4%ba%9a%7c%e8%8b%b1%e5%9b%bd%7c%e8%a5%bf%e7%8f%ad%e7%89%99%7c%e7%91%9e%e5%85%b8%7c%e5%8d%a2%e6%a3%ae%e5%a0%a1%7c%e5%be%b7%e5%9b%bd%7c%e7%88%b1%e5%b0%94%e5%85%b0%7c%e6%b3%a2%e5%85%b0%7c%e6%8d%b7%e5%85%8b%7c%e6%84%8f%e5%a4%a7%e5%88%a9%7cisx.yt%7c%e8%8d%b7%e5%85%b0%7cNL%7cUS%7cDE%7cFR%7cAU%7cRU%7cRO%7cZZ%7cGB%7c%e5%8d%b0%e5%ba%a6%7c%e4%b8%b9%e9%ba%a6")
		var ConfigPath string = c.DefaultQuery("config", "/etc/clash/config.yaml")
		var command string = c.DefaultQuery("command", "/etc/clash/clash.sh restart")
		var data string = GetData(geturl, suburl, exclude, ConfigPath, command)
		c.String(200, data)
	})
	app.Run(strings.Join([]string{":", port}, ""))
}

func GetData(geturl string, suburl string, exclude string, ConfigPath string, command string) (b string) {
	var (
		url      string
		emoji    string = "true"
		list     string = "false"
		udp      string = "true"
		tfo      string = "false"
		scv      string = "false"
		fdn      string = "false"
		sort     string = "false"
		template string = `port: 7890
socks-port: 7891
allow-lan: true
bind-address: "*"
ipv6: false
mode: rule
log-level: silent
external-controller: 0.0.0.0:9090
redir-port: 7892
secret: "123456"
external-ui: "/etc/clash/ui"
dns:
  enable: true
  ipv6: false
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16
  listen: 0.0.0.0:53
  fake-ip-filter:
##Custom fake-ip-filter##
  - '*.lan'
  - 'time.windows.com'
  - 'time.nist.gov'
  - 'time.apple.com'
  - 'time.asia.apple.com'
  - '*.ntp.org.cn'
  - '*.openwrt.pool.ntp.org'
  - 'time1.cloud.tencent.com'
  - 'time.ustc.edu.cn'
  - 'pool.ntp.org'
  - 'ntp.ubuntu.com'
  - 'ntp.aliyun.com'
  - 'ntp1.aliyun.com'
  - 'ntp2.aliyun.com'
  - 'ntp3.aliyun.com'
  - 'ntp4.aliyun.com'
  - 'ntp5.aliyun.com'
  - 'ntp6.aliyun.com'
  - 'ntp7.aliyun.com'
  - 'time1.aliyun.com'
  - 'time2.aliyun.com'
  - 'time3.aliyun.com'
  - 'time4.aliyun.com'
  - 'time5.aliyun.com'
  - 'time6.aliyun.com'
  - 'time7.aliyun.com'
  - '*.time.edu.cn'
  - 'time1.apple.com'
  - 'time2.apple.com'
  - 'time3.apple.com'
  - 'time4.apple.com'
  - 'time5.apple.com'
  - 'time6.apple.com'
  - 'time7.apple.com'
  - 'time1.google.com'
  - 'time2.google.com'
  - 'time3.google.com'
  - 'time4.google.com'
  - 'music.163.com'
  - '*.music.163.com'
  - '*.126.net'
  - 'musicapi.taihe.com'
  - 'music.taihe.com'
  - 'songsearch.kugou.com'
  - 'trackercdn.kugou.com'
  - '*.kuwo.cn'
  - 'api-jooxtt.sanook.com'
  - 'api.joox.com'
  - 'joox.com'
  - 'y.qq.com'
  - '*.y.qq.com'
  - 'streamoc.music.tc.qq.com'
  - 'mobileoc.music.tc.qq.com'
  - 'isure.stream.qqmusic.qq.com'
  - 'dl.stream.qqmusic.qq.com'
  - 'aqqmusic.tc.qq.com'
  - 'amobile.music.tc.qq.com'
  - '*.xiami.com'
  - '*.music.migu.cn'
  - 'music.migu.cn'
  - '*.msftconnecttest.com'
  - '*.msftncsi.com'
  - 'localhost.ptlogin2.qq.com'
  - '*.*.*.srv.nintendo.net'
  - '*.*.stun.playstation.net'
  - 'xbox.*.*.microsoft.com'
  - '*.*.xboxlive.com'
  - 'proxy.golang.org'
##Custom fake-ip-filter END##
  nameserver:
  - 223.5.5.5
  - 8.8.8.8
  fallback:
  - https://cloudflare-dns.com/dns-query
  - https://dns.google/dns-query
  - https://1.1.1.1/dns-query
  - tls://8.8.8.8:853
  fallback-filter:
    geoip: true
    ipcidr:
      - 0.0.0.0/8
      - 10.0.0.0/8
      - 100.64.0.0/10
      - 127.0.0.0/8
      - 169.254.0.0/16
      - 172.16.0.0/12
      - 192.0.0.0/24
      - 192.0.2.0/24
      - 192.88.99.0/24
      - 192.168.0.0/16
      - 198.18.0.0/15
      - 198.51.100.0/24
      - 203.0.113.0/24
      - 224.0.0.0/4
      - 240.0.0.0/4
      - 255.255.255.255/32
#===================== Clash-General-Settings =====================#`
	)
	url = strings.Join([]string{geturl, "/sub?target=clash", "&url=", suburl, "&exclude=", exclude, "&emoji=", emoji, "&list=", list, "&udp=", udp, "&tfo=", tfo, "&scv=", scv, "&fdn=", fdn, "&sort=", sort}, "")
	res, errers := http.Get(url)
	if errers != nil {
		return "get error"
	}
	if res != nil {
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return "status error"
		}
		s, _ := ioutil.ReadAll(res.Body)
		d := string(s)
		l := len(d)
		e := strings.Index(d, "proxies:")
		d = d[e : l-1]
		dd := strings.Join([]string{template, d}, "\n")
		var strtobyte []byte = []byte(dd)
		serr := ioutil.WriteFile(ConfigPath, strtobyte, 0644)
		if serr != nil {
			return "bad path"
		}
		cmd := exec.Command("/bin/sh", "-c", command)
		if err := cmd.Run(); err != nil {
			return "run error"
		}
		return "OK"
	}
	return "OK"
}
