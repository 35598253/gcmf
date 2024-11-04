package Ip2Info

import (
	"fmt"
	"strings"
)

// GIpInfo IP格式
type GIpInfo struct {
	CountryCode string `json:"country_code"`
	PhoneCode   string `json:"phone_code"`
	Country     string `json:"country"`
	Region      string `json:"region"`
	City        string `json:"city"`
	Isp         string `json:"isp"`
}

func GetInfo(Ip string) (*GIpInfo, error) {

	Ip2Path := "./resource/ip2region/ip2region.xdb"
	cBuff, err := LoadContentFromFile(Ip2Path)
	if err != nil {

		return nil, fmt.Errorf("failed to create searcher: %s\n", err.Error())
	}
	searcher, err := NewWithBuffer(cBuff)
	if err != nil {

		return nil, fmt.Errorf("failed to create searcher with content: %s\n", err)
	}
	region, _ := searcher.SearchByStr(Ip)
	str := strings.Split(region, "|")

	if len(str) < 6 {
		return nil, fmt.Errorf("failed to Ip Search Data")
	}
	var countryCode string
	var phoneCode string
	var country string
	var _region string
	var city string
	var isp string
	if str[0] != "0" {
		countryCode = str[0]
	}
	if str[1] != "0" {
		phoneCode = str[1]
	}
	if str[2] != "0" {
		country = str[2]
	}
	if str[3] != "0" {
		_region = str[3]
	}
	if str[4] != "0" {
		city = str[4]
	}
	if str[5] != "0" {
		isp = str[5]
	}
	return &GIpInfo{
		CountryCode: countryCode,
		PhoneCode:   phoneCode,
		Country:     country,
		Region:      _region,
		City:        city,
		Isp:         isp,
	}, nil

}
