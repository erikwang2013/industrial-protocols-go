package vendor

func RegisterRockwell(r *Registry) {
	r.Register(VendorProfile{Make: "Rockwell", Model: "CompactLogix", Defaults: map[string]any{
		"protocol": "ethernetip", "port": 44818,
	}})
	r.Register(VendorProfile{Make: "Rockwell", Model: "ControlLogix", Defaults: map[string]any{
		"protocol": "ethernetip", "port": 44818,
	}})
	r.Register(VendorProfile{Make: "Rockwell", Model: "MicroLogix", Defaults: map[string]any{
		"protocol": "modbus", "port": 502,
	}})
}
