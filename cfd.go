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

func main() {
	var (
		proxy      string = "true"
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
	recordIp, recordId, resSuccess := GetCurrentIP(username, password, hostname, recordType)
	if resSuccess {
		if recordIp == ipAddr {
			fmt.Println("nochg")
			return
		}
		if recordId == "null" {
			fmt.Println(recordId, proxy)
			return
		} else {
			fmt.Println(recordId)
			return
		}
	}
	fmt.Println("badauth")
	return
	// fmt.Println(proxy, username, password, hostname, ipAddr, ip, recordType)
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
	ip, err := getData("https://www.taobao.com/help/getip.php")
	if err == nil {
		length := len(ip)
		start := strings.Index(ip, `ip:"`)
		a := ip[start+4 : length]
		end := strings.Index(a, `"}`)
		i = a[0:end]
	}
	return
}

// getData get data
func getData(u string) (s string, err error) {
	res, err := http.Get(u)
	if err != nil {
		return "error", err
	}
	if res != nil {
		defer res.Body.Close()
		if res.StatusCode != 200 {
			return "error", err
		}
		d, _ := ioutil.ReadAll(res.Body)
		s = string(d)
	}
	return s, nil
}

// GetCurrentIP Get Current IP
func GetCurrentIP(username string, password string, hostname string, recordType string) (recordIp string, recordId string, resSuccess bool) {
	client := &http.Client{}
	url := strings.Join([]string{"https://api.cloudflare.com/client/v4/zones/", username, "/dns_records?type=", recordType, "&name=", hostname}, "")
	reqest, err := http.NewRequest("GET", url, nil)

	reqest.Header.Add("Authorization", strings.Join([]string{"Bearer", password}, " "))

	if err != nil {
		fmt.Println(err)
	}
	response, err := client.Do(reqest)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	d, _ := ioutil.ReadAll(response.Body)
	var p *Success
	//3.json解析到结构体
	if err := json.Unmarshal(d, &p); err != nil {
		fmt.Println(err)
	}
	if p.Success {
		return p.Result[0].ID, p.Result[0].Content, p.Success
	}
	return "", "", p.Success
}

// postData
func postData(url string, proxy string, ipAddr string, hostname string, recordType string) (bd []byte, f bool) {
	data := make(map[string]string)
	data["type"] = recordType
	data["name"] = hostname
	data["content"] = ipAddr
	data["proxied"] = proxy
	b, _ := json.Marshal(data)

	resp, err := http.Post(url,
		"application/json",
		bytes.NewBuffer(b))
	if err != nil {
		return []byte(""), false
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, true
}
