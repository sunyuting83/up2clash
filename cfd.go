package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

type Success struct {
	Result     []Result      `json:"result"`
	Success    bool          `json:"success"`
	Errors     []interface{} `json:"errors"`
	Messages   []interface{} `json:"messages"`
	ResultInfo ResultInfo    `json:"result_info"`
}
type Meta struct {
	AutoAdded           bool   `json:"auto_added"`
	ManagedByApps       bool   `json:"managed_by_apps"`
	ManagedByArgoTunnel bool   `json:"managed_by_argo_tunnel"`
	Source              string `json:"source"`
}
type Result struct {
	ID         string    `json:"id"`
	ZoneID     string    `json:"zone_id"`
	ZoneName   string    `json:"zone_name"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Content    string    `json:"content"`
	Proxiable  bool      `json:"proxiable"`
	Proxied    bool      `json:"proxied"`
	TTL        int       `json:"ttl"`
	Locked     bool      `json:"locked"`
	Meta       Meta      `json:"meta"`
	CreatedOn  time.Time `json:"created_on"`
	ModifiedOn time.Time `json:"modified_on"`
}
type ResultInfo struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Count      int `json:"count"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

type PostData struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Proxied bool   `json:"proxied"`
}

func main() {
	var (
		CFApi      string = "https://api.cloudflare.com/client/v4/zones/"
		proxy      bool   = false
		username   string = os.Args[1]
		password   string = os.Args[2]
		hostname   string = os.Args[3]
		ipAddr     string
		recordType string = "A"
	)
	ipAddr = GetIpAddr()
	ip, a := ParseIP(ipAddr)
	if ip != nil {
		if a != 4 {
			recordType = "AAAA"
		}
	}
	var CurrentUrl string = strings.Join([]string{CFApi, username, "/dns_records?type=", recordType, "&name=", hostname}, "")
	recordId, recordIp, resSuccess := CloudFlareApi(CurrentUrl, "GET", password, []byte(""), true)
	if resSuccess {
		if recordIp == ipAddr {
			fmt.Println("nochg")
			return
		}
		data := MakePostData(proxy, ipAddr, hostname, recordType)
		if recordId == "null" {
			var createDnsApi string = strings.Join([]string{CFApi, username, "/dns_records"}, "")
			_, _, success := CloudFlareApi(createDnsApi, "POST", password, data, false)
			if success {
				fmt.Println("good")
				return
			} else {
				fmt.Println("badauth")
				return
			}
		} else {
			var updateDnsApi string = strings.Join([]string{CFApi, username, "/dns_records/", recordId}, "")
			_, _, success := CloudFlareApi(updateDnsApi, "PUT", password, data, false)
			if success {
				fmt.Println("good")
				return
			} else {
				fmt.Println("badauth")
				return
			}
		}
	}
	fmt.Println("badauth")
	return
}

// ParseIP Parse IP Type
func ParseIP(s string) (net.IP, int) {
	ip := net.ParseIP(s)
	if ip == nil {
		return nil, 0
	}
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return ip, 4
		case ':':
			return ip, 6
		}
	}
	return nil, 0
}

// GetIpAddr get ip addr
func GetIpAddr() (i string) {
	ip, err := getData("https://www.taobao.com/help/getip.php", "GET", []byte(""), "")
	if err == nil {
		ips := string(ip)
		length := len(ips)
		start := strings.Index(ips, `ip:"`)
		a := ips[start+4 : length]
		end := strings.Index(a, `"}`)
		i = a[0:end]
	}
	return
}

// getData get data
func getData(url string, types string, data []byte, password string) (s []byte, err error) {
	client := &http.Client{}
	reqest, err := http.NewRequest(types, url, bytes.NewBuffer(data))

	if len(password) > 0 {
		reqest.Header.Add("Authorization", strings.Join([]string{"Bearer", password}, " "))
	}

	if err != nil {
		return []byte(""), err
	}
	response, err := client.Do(reqest)
	if err != nil {
		return []byte(""), err
	}
	defer response.Body.Close()
	d, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return []byte(""), err
	}
	return d, nil
}

// CloudFlareApi Cloud Flare Api
func CloudFlareApi(url string, types string, password string, data []byte, getIP bool) (recordIp string, recordId string, resSuccess bool) {
	d, err := getData(url, types, data, password)
	if err != nil {
		return "", "", false
	}
	var p *Success
	//3.json解析到结构体
	if err := json.Unmarshal(d, &p); err != nil {
		return "", "", p.Success
	}
	if p.Success {
		if getIP {
			return p.Result[0].ID, p.Result[0].Content, p.Success
		}
		return "", "", p.Success
	}
	return "", "", p.Success
}

// MakePostData Make post data
func MakePostData(proxy bool, ipAddr string, hostname string, recordType string) (bd []byte) {
	var d *PostData
	d = &PostData{
		Type:    recordType,
		Name:    hostname,
		Content: ipAddr,
		Proxied: proxy,
	}
	b, _ := json.Marshal(d)
	return b
}
