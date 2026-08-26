// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestLoadYAML(t *testing.T) {
	path := filepath.Join("testdata", "example.yaml")
	repo, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	cfg, ok := repo.Get("plc-001")
	if !ok {
		t.Fatal("plc-001 not found")
	}
	if cfg.Addr != "192.168.1.10:502" {
		t.Errorf("expected 192.168.1.10:502, got %s", cfg.Addr)
	}
	if cfg.Pool == nil || cfg.Pool.MaxSize != 5 {
		t.Error("expected pool config")
	}
	_, ok = repo.Get("nonexistent")
	if ok {
		t.Error("should not find nonexistent device")
	}
}

func TestLoadMissingFile(t *testing.T) {
	_, err := Load("/nonexistent.yaml")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestTempYAML(t *testing.T) {
	content := []byte("devices:\n  test:\n    name: test\n    network: tcp\n    addr: \"127.0.0.1:502\"\n")
	tmpfile := filepath.Join(t.TempDir(), "config.yaml")
	os.WriteFile(tmpfile, content, 0644)
	repo, err := Load(tmpfile)
	if err != nil {
		t.Fatal(err)
	}
	cfg, ok := repo.Get("test")
	if !ok {
		t.Fatal("test device not found")
	}
	if cfg.Addr != "127.0.0.1:502" {
		t.Errorf("wrong addr: %s", cfg.Addr)
	}
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	tmpfile := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(tmpfile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return tmpfile
}

func TestLoad_InvalidYAML(t *testing.T) {
	if _, err := Load(writeConfig(t, "devices: [unclosed")); err == nil {
		t.Error("expected YAML parse error")
	}
}

func TestLoad_InvalidDeviceTimeout(t *testing.T) {
	if _, err := Load(writeConfig(t, "devices:\n  d:\n    name: d\n    timeout: \"abc\"\n")); err == nil {
		t.Error("expected duration parse error")
	}
}

func TestLoad_InvalidPoolIdleTimeout(t *testing.T) {
	_, err := Load(writeConfig(t, "devices:\n  d:\n    name: d\n    pool:\n      max_size: 2\n      idle_timeout: \"xyz\"\n"))
	if err == nil || !strings.Contains(err.Error(), "invalid idle_timeout") {
		t.Errorf("expected invalid idle_timeout error, got %v", err)
	}
}

func TestLoad_Defaults(t *testing.T) {
	repo, err := Load(writeConfig(t, "devices:\n  d:\n    name: d\n    addr: \"1.2.3.4:502\"\n"))
	if err != nil {
		t.Fatal(err)
	}
	cfg, ok := repo.Get("d")
	if !ok {
		t.Fatal("device not found")
	}
	if cfg.Addr != "1.2.3.4:502" {
		t.Errorf("wrong addr: %s", cfg.Addr)
	}
	if cfg.Timeout != 0 {
		t.Errorf("expected zero timeout, got %v", cfg.Timeout)
	}
	if cfg.Pool != nil {
		t.Error("expected nil pool")
	}
}

func TestLoad_PoolDefaults(t *testing.T) {
	repo, err := Load(writeConfig(t, "devices:\n  d:\n    name: d\n    pool:\n      max_size: 3\n"))
	if err != nil {
		t.Fatal(err)
	}
	cfg, _ := repo.Get("d")
	if cfg.Pool == nil {
		t.Fatal("expected pool config")
	}
	if cfg.Pool.MaxSize != 3 {
		t.Errorf("expected max_size 3, got %d", cfg.Pool.MaxSize)
	}
	if cfg.Pool.IdleTimeout != 60*time.Second {
		t.Errorf("expected default idle_timeout 60s, got %v", cfg.Pool.IdleTimeout)
	}
}

func TestLoad_AllReturnsCopy(t *testing.T) {
	repo, err := Load(writeConfig(t, "devices:\n  a:\n    name: a\n    addr: \"1:2\"\n  b:\n    name: b\n    addr: \"1:2\"\n"))
	if err != nil {
		t.Fatal(err)
	}
	all := repo.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(all))
	}
	delete(all, "a")
	if _, ok := repo.Get("a"); !ok {
		t.Error("All() should return a copy, repo must be unaffected")
	}
	if len(repo.All()) != 2 {
		t.Error("repo should still hold both devices")
	}
}
