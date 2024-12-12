package acl

// разрешить рекурсивные запросы для IP-адреса
func (a *ACL) AllowRecursion(ip string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.allowRecursion[ip] = true
}

// запретить рекурсивные запросы для IP-адреса
func (a *ACL) DenyRecursion(ip string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	delete(a.allowRecursion, ip)
}

// проверить, разрешены ли рекурсивные запросы для IP-адреса
func (a *ACL) IsRecursionAllowed(ip string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	return a.allowRecursion[ip]
}
