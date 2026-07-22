// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package config

import (
	"os"
	"time"

	"github.com/erikwang2013/industrial-protocols-go/kernel/connection"
	"gopkg.in/yaml.v3"
)

type configFile struct {
	Devices map[string]deviceYAML `yaml:"devices"`
}

type poolYAML struct {
	MaxSize     int    `yaml:"max_size"`
	IdleTimeout string `yaml:"idle_timeout"`
}

type deviceYAML struct {
	Name    string    `yaml:"name"`
	Network string    `yaml:"network"`
	Addr    string    `yaml:"addr"`
	Timeout string    `yaml:"timeout"`
	Pool    *poolYAML `yaml:"pool,omitempty"`
}

type ConfigRepository struct {
	devices map[string]*connection.DeviceConfig
}

func Load(path string) (*ConfigRepository, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cf configFile
	if err := yaml.Unmarshal(data, &cf); err != nil {
		return nil, err
	}
	repo := &ConfigRepository{devices: make(map[string]*connection.DeviceConfig)}
	for key, dy := range cf.Devices {
		dc := &connection.DeviceConfig{
			Name: dy.Name, Network: dy.Network, Addr: dy.Addr,
		}
		if dy.Timeout != "" {
			d, err := time.ParseDuration(dy.Timeout)
			if err != nil {
				return nil, err
			}
			dc.Timeout = d
		}
		if dy.Pool != nil {
			idleTimeout := 60 * time.Second
			if dy.Pool.IdleTimeout != "" {
				idleTimeout, _ = time.ParseDuration(dy.Pool.IdleTimeout)
			}
			dc.Pool = &connection.PoolConfig{
				MaxSize: dy.Pool.MaxSize, IdleTimeout: idleTimeout,
			}
		}
		repo.devices[key] = dc
	}
	return repo, nil
}

func (r *ConfigRepository) Get(name string) (*connection.DeviceConfig, bool) {
	cfg, ok := r.devices[name]
	return cfg, ok
}

func (r *ConfigRepository) All() map[string]*connection.DeviceConfig {
	result := make(map[string]*connection.DeviceConfig, len(r.devices))
	for k, v := range r.devices {
		result[k] = v
	}
	return result
}
