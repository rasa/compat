module github.com/rasa/compat

// tinygo requires 1.24
go 1.24.2

toolchain go1.24.5

require (
	github.com/OneOfOne/xxhash v1.2.8
	github.com/adrg/xdg v0.5.3
	github.com/capnspacehook/go-acl v0.0.0-20191017210724-28a28d1b4c77
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/mattn/go-colorable v0.1.14
	github.com/sergi/go-diff v1.4.0
	github.com/shirou/gopsutil/v4 v4.25.8
	golang.org/x/sys v0.36.0
)

require (
	github.com/ebitengine/purego v0.8.4 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/hectane/go-acl v0.0.0-20230122075934-ca0b05cb1adb // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	golang.org/x/mod v0.26.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/telemetry v0.0.0-20250710130107-8d8967aff50b // indirect
	golang.org/x/tools v0.35.1-0.20250728180453-01a3475a31bc // indirect
	golang.org/x/tools/gopls v0.20.0 // indirect
)

// replace github.com/capnspacehook/go-acl v0.0.0-20191017210724-28a28d1b4c77 => ../go-acl
