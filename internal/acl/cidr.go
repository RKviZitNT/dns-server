package acl

import "net"

// добавить CIDR-диапазон в список разрешённых
func (a *ACL) AllowCIDR(cidr string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}
	a.allowCIDRs = append(a.allowCIDRs, ipNet)
	return nil
}

// удалить CIDR-диапазон из списка разрешённых
func (a *ACL) DenyCIDR(cidr string) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	// удаляем CIDR из списка
	for i, existingCIDR := range a.allowCIDRs {
		if existingCIDR.String() == ipNet.String() {
			a.allowCIDRs = append(a.allowCIDRs[:i], a.allowCIDRs[i+1:]...)
			break
		}
	}
	return nil
}
