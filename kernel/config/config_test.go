// Copyright (c) 2026 erik <erik@erik.xyz> — https://erik.xyz

package config

import (
	"os"
	"path/filepath"
	"testing"
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
