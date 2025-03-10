package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	"github.com/waldur/waldur-rancher-node-driver/driver"
)

func main() {
	plugin.RegisterDriver(NewDriver("", ""))
}
