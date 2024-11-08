module github.com/liuxd6825/jsonschema/cmd/jv

go 1.21.1

require (
	github.com/liuxd6825/jsonschema/v6 v6.0.1
	github.com/spf13/pflag v1.0.5
	gopkg.in/yaml.v3 v3.0.1
)

require golang.org/x/text v0.14.0 // indirect

replace github.com/liuxd6825/jsonschema/v6 v6.0.1 => ../..
