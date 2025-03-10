module github.com/waldur/waldur-rancher-node-driver

go 1.23.7

replace github.com/docker/docker => github.com/moby/moby v1.4.2-0.20170731201646-1009e6a40b29

require github.com/docker/machine v0.16.2

require (
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/docker/docker v0.0.0-00010101000000-000000000000 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
)
