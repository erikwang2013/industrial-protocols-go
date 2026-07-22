// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package vendor

import "sync"

type VendorProfile struct {
	Make     string
	Model    string
	Defaults map[string]any
}

type Registry struct {
	mu     sync.RWMutex
	byName map[string]*VendorProfile
}

func NewRegistry() *Registry {
	return &Registry{byName: make(map[string]*VendorProfile)}
}

func (r *Registry) Register(p VendorProfile) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.byName[p.Make+"/"+p.Model] = &p
}

func (r *Registry) Find(make, model string) (*VendorProfile, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.byName[make+"/"+model]
	return p, ok
}
