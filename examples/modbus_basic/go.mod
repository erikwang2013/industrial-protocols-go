module github.com/erikwang2013/industrial-protocols-go/examples/modbus_basic

go 1.24.1

require (
	github.com/erikwang2013/industrial-protocols-go/kernel v0.0.0
	github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/modbus v0.0.0
)

replace (
	github.com/erikwang2013/industrial-protocols-go/kernel => ../../kernel
	github.com/erikwang2013/industrial-protocols-go/protocols/ethernet/modbus => ../../protocols/ethernet/modbus
)
