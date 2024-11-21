package blockchain

import (
	"fmt"
	"log"
	"net"
)

type Node struct {
	IP         string
	ListenPort int
	Validator  string
}

func (n *Node) NewNode(validator string, listenPort int) *Node {
	ip, err := GetPreferredIP()
	if err != nil {
		log.Fatalf("Failed to get IP: %v", err)
	}

	return &Node{
		IP:         ip,
		ListenPort: listenPort,
		Validator:  validator,
	}
}

func (n *Node) GetIP() string {
	return fmt.Sprintf("%s:%d", n.IP, n.ListenPort)
}

func GetPreferredIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	//OS 레벨에서 MAC 정보
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip != nil && !ip.IsLoopback() && ip.IsGlobalUnicast() {
				return ip.String(), nil
			}
		}
	}

	return "", fmt.Errorf("no valid IP address found")
}

func (n *Node) IsPublicIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil || !parsedIP.IsGlobalUnicast() {
		return false
	}

	// 사설 ip대역
	privateIPRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}
	for _, cidr := range privateIPRanges {
		_, subnet, _ := net.ParseCIDR(cidr)
		if subnet.Contains(parsedIP) {
			return false
		}
	}

	return true
}

func IsValidator(validator string) bool {
	return true
}
