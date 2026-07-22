// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package vendor

import "sync"

// VendorProfile 厂商参数预设。
type VendorProfile struct {
	Make     string
	Model    string
	Defaults map[string]any
}

// Registry 厂商注册表，线程安全。
type Registry struct {
	mu     sync.RWMutex
	byName map[string]*VendorProfile
}

// NewRegistry 创建注册表。
func NewRegistry() *Registry {
	return &Registry{byName: make(map[string]*VendorProfile)}
}

// Register 注册一个厂商预设。
func (r *Registry) Register(p VendorProfile) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byName[p.Make+"/"+p.Model] = &p
}

// Find 查找厂商预设。
func (r *Registry) Find(make, model string) (*VendorProfile, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.byName[make+"/"+model]
	return p, ok
}
