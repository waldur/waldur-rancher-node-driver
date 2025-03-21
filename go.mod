module github.com/waldur/waldur-rancher-node-driver

go 1.24.1

replace github.com/docker/docker => github.com/moby/moby v1.4.2-0.20170731201646-1009e6a40b29

require (
	github.com/google/uuid v1.6.0
	github.com/rancher/machine v0.15.0-rancher126
	github.com/waldur/go-client v0.0.0-20250315211045-7b98813ba8f8
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20250102033503-faa5f7b0171c // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/docker/docker v28.0.1+incompatible // indirect
	github.com/oapi-codegen/runtime v1.1.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
)
