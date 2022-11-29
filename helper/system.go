package helper

import (
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

var appName string

func init() {
	rand.Seed(time.Now().UnixNano())
	appName = GetEnvDefault("APP_NAME", "")
}

func AppName() string {
	return appName
}

// 取环境变量
func GetEnvDefault(key, defVal string) string {
	val, ex := os.LookupEnv(key)
	if !ex {
		return defVal
	}
	return val
}

// 获取本机mac地址
func GetMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return macAddrs
	}
	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		macAddrs = append(macAddrs, macAddr)
	}
	return macAddrs
}

// 获取本机所有ip地址
func GetLocalIPs() (ips []string) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}
	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ips = append(ips, ipNet.IP.String())
			}
		}
	}
	return ips
}

// 获取本机ip地址，默认获取对外的ip地址
func GetLocalIP(address ...string) string {
	var err error
	address = append(address, []string{
		"223.5.5.5:53", "8.8.8.8:53",
	}...)
	for _, addr := range address {
		ip := net.ParseIP(addr)
		if ip != nil {
			switch IpType(ip.String()) {
			case 4:
				addr = (addr + ":80")
			case 6:
				addr = ("[" + addr + "]:80")
			}
		}
		conn, err := net.Dial("udp", addr)
		if err != nil {
			continue
		}
		defer conn.Close()
		localAddr := conn.LocalAddr().(*net.UDPAddr)
		return localAddr.IP.String()
	}
	log.Println("GetLocalIP error：", err)
	return ""
}
