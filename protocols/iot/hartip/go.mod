module github.com/erikwang2013/industrial-protocols-go/protocols/iot/hartip

go 1.24.1

require (
	github.com/erikwang2013/industrial-protocols-go/kernel v0.0.0
	github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/hart v0.0.0
)

replace (
	github.com/erikwang2013/industrial-protocols-go/kernel => ../../../kernel
	github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/hart => ../../fieldbus/hart
)
