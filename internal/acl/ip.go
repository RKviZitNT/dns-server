package acl

import "net"

// добавить IP-адрес в список разрешённых
func (a *ACL) AllowIP(ip string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.allowIPs[ip] = true
}

// удалить IP-адрес из списка разрешённых
func (a *ACL) DenyIP(ip string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	delete(a.allowIPs, ip)
}

// проверить доступ для IP-адреса
func (a *ACL) IsAllowed(ip string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	// проверяем, есть ли IP в списке разрешённых
	if a.allowIPs[ip] {
		return true
	}

	// проверяем, попадает ли IP в один из CIDR-диапазонов
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return false
	}

	for _, cidr := range a.allowCIDRs {
		if cidr.Contains(parsedIP) {
			return true
		}
	}

	return false
}
