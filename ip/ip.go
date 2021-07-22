package ip

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	// IndexLength 定义
	IndexLength = 12
)

// Region Ip转
type Region struct {
	// db file handler
	dbFileHandler *os.File

	// super block index info
	firstIndexPtr int64
	lastIndexPtr  int64
	totalBlocks   int64

	// for memory mode only
	// the original db binary string

	dbBinStr []byte
	dbFile   string
}

// Info 信息
type Info struct {
	Country  string
	Province string
	City     string
	Region   string
	ISP      string
}

// getIPInfo 获取Ip信息
func getIPInfo(line []byte) Info {
	lineSlice := strings.Split(string(line), "|")
	ipInfo := Info{}
	length := len(lineSlice)
	if length < 5 {
		for i := 0; i <= 5-length; i++ {
			lineSlice = append(lineSlice, "")
		}
	}

	ipInfo.Country = lineSlice[0]
	ipInfo.Region = lineSlice[1]
	ipInfo.Province = lineSlice[2]
	ipInfo.City = lineSlice[3]
	ipInfo.ISP = lineSlice[4]
	return ipInfo
}

// New 初始化IP服务
func New(path string) (*Region, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	file.Close()
	return &Region{
		dbFile:        path,
		dbFileHandler: file,
	}, nil

}

// Search 查询方式这里只是进行内存查询
func (c *Region) Search(ipStr string) (ipInfo Info, err error) {
	ipInfo = Info{}

	if c.totalBlocks == 0 {
		c.dbBinStr, err = ioutil.ReadFile(c.dbFile)

		if err != nil {

			return ipInfo, err
		}

		c.firstIndexPtr = getLong(c.dbBinStr, 0)
		c.lastIndexPtr = getLong(c.dbBinStr, 4)
		c.totalBlocks = (c.lastIndexPtr-c.firstIndexPtr)/IndexLength + 1
	}

	ip, err := ip2long(ipStr)
	if err != nil {
		return ipInfo, err
	}

	h := c.totalBlocks
	var dataPtr, l int64
	for l <= h {

		m := (l + h) >> 1
		p := c.firstIndexPtr + m*IndexLength
		sip := getLong(c.dbBinStr, p)
		if ip < sip {
			h = m - 1
		} else {
			eip := getLong(c.dbBinStr, p+4)
			if ip > eip {
				l = m + 1
			} else {
				dataPtr = getLong(c.dbBinStr, p+8)
				break
			}
		}
	}
	if dataPtr == 0 {
		return ipInfo, errors.New("not found")
	}

	dataLen := (dataPtr >> 24) & 0xFF
	dataPtr = dataPtr & 0x00FFFFFF
	ipInfo = getIPInfo(c.dbBinStr[(dataPtr)+4 : dataPtr+dataLen])
	return ipInfo, nil

}
func getLong(b []byte, offset int64) int64 {

	val := int64(b[offset]) | int64(b[offset+1])<<8 | int64(b[offset+2])<<16 | int64(b[offset+3])<<24

	return val

}

func ip2long(IPStr string) (int64, error) {
	bits := strings.Split(IPStr, ".")
	if len(bits) != 4 {
		return 0, errors.New("ip format error")
	}

	var sum int64
	for i, n := range bits {
		bit, _ := strconv.ParseInt(n, 10, 64)
		sum += bit << uint(24-8*i)
	}

	return sum, nil
}
