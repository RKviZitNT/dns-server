package acl

import (
	"net"
	"sync"
)

type ACL struct {
	allowIPs       map[string]bool // cписок разрешённых IP-адресов
	allowCIDRs     []*net.IPNet    // cписок разрешённых CIDR-диапазонов
	allowRecursion map[string]bool // cписок IP-адресов, которым разрешены рекурсивные запросы
	mu             sync.Mutex
}

func NewACL() *ACL {
	return &ACL{
		allowIPs:       make(map[string]bool),
		allowCIDRs:     make([]*net.IPNet, 0),
		allowRecursion: make(map[string]bool),
	}
}
