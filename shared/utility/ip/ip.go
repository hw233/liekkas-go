package ip

import (
	"fmt"
	"net"

	"shared/utility/errors"
)

func LocalAddr() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return "", nil
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ip, ok := address.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String(), nil
			}
		}
	}

	return "", errors.New("not found addr")
}
