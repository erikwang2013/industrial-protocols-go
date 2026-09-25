module github.com/erikwang2013/industrial-protocols-go/protocols/iot/hartip

go 1.24.1

require (
	github.com/erikwang2013/industrial-protocols-go/kernel v1.1.4
	github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/hart v1.1.4
)

replace (
	github.com/erikwang2013/industrial-protocols-go/kernel => ../../../kernel
	github.com/erikwang2013/industrial-protocols-go/protocols/fieldbus/hart => ../../fieldbus/hart
)
