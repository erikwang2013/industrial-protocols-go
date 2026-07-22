package vendor

func RegisterSiemens(r *Registry) {
	r.Register(VendorProfile{Make: "Siemens", Model: "S7-1200", Defaults: map[string]any{
		"protocol": "modbus", "port": 502, "endian": "big",
	}})
	r.Register(VendorProfile{Make: "Siemens", Model: "S7-1500", Defaults: map[string]any{
		"protocol": "profinet", "port": 34964,
	}})
	r.Register(VendorProfile{Make: "Siemens", Model: "S7-300", Defaults: map[string]any{
		"protocol": "profibus",
	}})
}
