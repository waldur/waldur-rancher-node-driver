module github.com/waldur/waldur-rancher-node-driver

go 1.25.0

replace github.com/docker/docker => github.com/moby/moby v1.4.2-0.20170731201646-1009e6a40b29

require (
	github.com/google/uuid v1.6.0
	github.com/rancher/machine v0.16.2
	github.com/waldur/go-client v0.0.0-20260401113022-40a829c05f95
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/docker/docker v28.5.2+incompatible // indirect
	github.com/docker/machine v0.16.2 // indirect
	github.com/oapi-codegen/runtime v1.3.1 // indirect
	github.com/sirupsen/logrus v1.9.4 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/term v0.41.0 // indirect
)
