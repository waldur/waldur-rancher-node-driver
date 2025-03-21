package main

import (
	"github.com/rancher/machine/libmachine/drivers/plugin"
	"github.com/waldur/waldur-rancher-node-driver/driver"
)

func main() {
	plugin.RegisterDriver(driver.NewDriver("", ""))
}
