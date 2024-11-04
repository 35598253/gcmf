package gcmf

import (
	"admin/hack/gcmf/Ip2Info"
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"

	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
)

// Ip 获取
var Ip = new(ipActive)

type ipActive struct {
}

// OutInfo 信息
type OutInfo struct {
	Status   string
	Info     string
	InfoCode string
	Country  string
	Province string
	City     string
	District string
	Isp      string
	Location string
	Ip       string
}

// GIpInfo 自定义IP格式
type GIpInfo struct {
	CountryCode string `json:"country_code"`
	PhoneCode   string `json:"phone_code"`
	Country     string `json:"country"`
	Region      string `json:"region"`
	City        string `json:"city"`
	Isp         string `json:"isp"`
}

// GetIP 获取Ip
func (i *ipActive) GetIP(r *ghttp.Request) string {
	ip := r.Header.Get("X-Real-IP")

	if net.ParseIP(ip) != nil {
		return ip
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i
		}
	}

	ip, _, _ = net.SplitHostPort(r.RemoteAddr)

	if net.ParseIP(ip) != nil {
		return ip
	}

	return "127.0.0.1"
}

// GetIpInfo 获取IpInfo
func (i *ipActive) GetIpInfo(Ctx context.Context, Ip string, ip6 ...bool) (string, error) {

	if Ip == "127.0.0.1" {
		return "内网IP", nil
	}
	cfg, err := GetConfig(Ctx)
	if err != nil {
		return "", err
	}
	var Ip6 bool
	var ipInfo string
	if len(ip6) > 0 {
		Ip6 = ip6[0]
	}
	ipCfg := cfg.Ip
	//本地Ip
	if ipCfg.Type == 0 {
		if Ip6 {
			errs := errors.New("not support IP6")
			return "", errs
		}

		ipInfo, err = i.IpInfo(Ip)
		if err != nil {
			return "", err
		}
	}
	if ipCfg.Type == 1 {
		key := ipCfg.AmapKey
		ipInfo = i.AmapInfo(Ctx, key, Ip, Ip6)
	}

	return ipInfo, nil
}

// AmapInfo 获取Ip信息
func (i *ipActive) AmapInfo(Ctx context.Context, Key, Ip string, ip6 ...bool) string {

	var ipInfo *OutInfo
	var _type = 4
	if len(ip6) > 0 && ip6[0] {
		_type = 6
	}
	var url = fmt.Sprintf("https://restapi.amap.com/v5/ip?key=%s&ip=%s&type=%d", Key, Ip, _type)
	// 发起请求
	_bytes := g.Client().Timeout(16*time.Second).GetBytes(Ctx, url)
	info := "--"
	_ = gconv.Struct(_bytes, &ipInfo)
	if ipInfo.Info == "OK" {
		info = fmt.Sprintf("%s-%s-%s-%s", ipInfo.Country, ipInfo.Province, ipInfo.City, ipInfo.Isp)
	}

	return info
}

// IpInfo 地址
func (i *ipActive) IpInfo(Ip string) (string, error) {

	ipInfo, err := Ip2Info.GetInfo(Ip)

	if err != nil {
		return "", err
	}
	info := fmt.Sprintf("%s-%s-%s-%s", ipInfo.Country, ipInfo.Region, ipInfo.City, ipInfo.Isp)
	return info, nil
}
