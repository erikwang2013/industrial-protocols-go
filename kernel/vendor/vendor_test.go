package vendor

import "testing"

func TestRegistryRegisterAndFind(t *testing.T) {
	r := NewRegistry()
	RegisterSiemens(r)
	RegisterRockwell(r)

	p, ok := r.Find("Siemens", "S7-1200")
	if !ok {
		t.Fatal("S7-1200 not found")
	}
	if p.Make != "Siemens" {
		t.Errorf("expected Siemens, got %s", p.Make)
	}
	if p.Defaults["protocol"] != "modbus" {
		t.Errorf("expected modbus protocol, got %v", p.Defaults["protocol"])
	}

	_, ok = r.Find("Unknown", "X")
	if ok {
		t.Error("should not find unknown vendor")
	}
}

func TestRegisterAllDefaults(t *testing.T) {
	r := NewRegistry()
	RegisterSiemens(r)
	RegisterRockwell(r)

	p1, _ := r.Find("Siemens", "S7-1500")
	if p1.Defaults["protocol"] != "profinet" {
		t.Error("S7-1500 should default to profinet")
	}

	p2, _ := r.Find("Rockwell", "ControlLogix")
	if p2.Defaults["port"].(int) != 44818 {
		t.Error("ControlLogix should default to port 44818")
	}
}
