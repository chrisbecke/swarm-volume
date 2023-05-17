module docker-volume-plugin

go 1.19

require (
	driver v0.0.0-00010101000000-000000000000
	github.com/docker/go-plugins-helpers v0.0.0-20211224144127-6eecb7beb651
)

replace driver => ../driver

require (
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	golang.org/x/mod v0.8.0 // indirect
	golang.org/x/net v0.6.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
)
