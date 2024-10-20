package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

func InitGeoDB() (*geoip2.Reader, error) {
	return GetGeoData()
}

func checkIPGeo(geoDB *geoip2.Reader, ip net.IP) (string, error) {
	record, err := geoDB.Country(ip)
	if err != nil {
		return "CDN", err
	}
	return record.Country.IsoCode, nil
}

// 处理加密DNS地址，
// 示例返回值
// 208.67.220.123,208.67.220.123,US
// https://doh.familyshield.opendns.com/dns-query,146.112.41.3,US
// tls://familyshield.opendns.com,208.67.222.123,US
// https://freedns.controld.com/p3,...
// https://dns.bebasid.com/unfiltered,...
// 2620:119:53::53,...
// https://doh.cleanbrowsing.org/doh/family-filter/,...
func CheckGeo(geoDB *geoip2.Reader, _server string, preferIPv4 bool) (string, string, error) {
	server := strings.TrimSpace(_server)
	server = strings.TrimSuffix(server, "/")
	if server == "" {
		return "0.0.0.0", "PRIVATE", fmt.Errorf("服务器地址为空")
	}
	var ip net.IP
	if strings.Contains(server, "://") {
		// URL
		server = strings.TrimPrefix(server, "https://")
		server = strings.TrimPrefix(server, "tls://")
		server = strings.TrimPrefix(server, "quic://")
		server = strings.TrimPrefix(server, "http://")

		if strings.Contains(server, "/") {
			// 带路径
			parts := strings.SplitN(server, "/", 2)
			server = parts[0]
		}
		if strings.Contains(server, "[") && strings.Contains(server, "]") {
			// IPv6 网址
			server = strings.SplitN(server, "]", 2)[0]
			server = strings.TrimPrefix(server, "[")
		} else if strings.Contains(server, ":") {
			// 普通 URL 带端口
			parts := strings.SplitN(server, ":", 2)
			server = parts[0]
		}
		// 解析成 IP
		ips, err := net.LookupIP(server)
		ipc := len(ips)
		if err != nil || ipc == 0 {
			// 无法解析IP地址
			return "0.0.0.0", "PRIVATE", fmt.Errorf("无法解析IP地址")
		}
		if ipc == 1 {
			// 只有一个IP地址
			ip = ips[0]
		} else {
			// 多个 IP 地址
			if preferIPv4 {
				for _, _ip := range ips {
					if _ip.To4() != nil {
						ip = _ip
						break
					}
				}
				if ip == nil {
					ip = ips[0]
				}
			} else {
				ip = ips[0]
			}
		}
	} else {
		// IP
		colonCount := strings.Count(server, ":")
		if colonCount == 1 && strings.Contains(server, ".") {
			// IPv4 带端口
			server = strings.SplitN(server, ":", 2)[0]
		}
		ip = net.ParseIP(server)
		if ip == nil {
			ip = net.IPv4zero
		}
	}
	if ip.IsPrivate() || ip.IsUnspecified() {
		return ip.String(), "PRIVATE", fmt.Errorf("IP地址为私有地址")
	}
	geoCode, err := checkIPGeo(geoDB, ip)
	return ip.String(), geoCode, err
}