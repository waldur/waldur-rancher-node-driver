package driver

import (
	"fmt"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/state"
)

const (
	driverName = "waldur"
)

type Driver struct {
	*drivers.BaseDriver

	ApiUrl       string
	ApiToken       string
}

// NewDriver creates and returns a new instance of Waldur driver
func NewDriver(hostName, storePath string) *Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{
			MachineName: hostName,
			StorePath:   storePath,
		},
	}
}

// GetCreateFlags registers the flags this driver adds to
// "docker hosts create"
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "WALDUR_API_URL",
			Name:   "waldur-api-url",
			Usage:  "Waldur API URL",
		},
		mcnflag.StringFlag{
			EnvVar: "WALDUR_API_TOKEN",
			Name:   "waldur-api-token",
			Usage:  "Waldur API URL",
		},
		mcnflag.StringFlag{
			EnvVar: "WALDUR_PROJ_UUID",
			Name:   "waldur-proj-uuid",
			Usage:  "UUID of the project in Waldur",
		},
	}
}

// SetConfigFromFlags configures the driver with the command line arguments
func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.ApiUrl = flags.String("waldur-api-url")
	d.ApiToken = flags.String("waldur-api-token")

	// Validation
	if d.ApiUrl == "" {
		return fmt.Errorf("Waldur requires the --waldur-api-url option")
	}

	if d.ApiToken == "" {
		return fmt.Errorf("Waldur requires the --waldur-api-token option")
	}

	return nil
}

// Create creates a host in Waldur using the driver's config
func (d *Driver) Create() error {
	// Implement your provider's logic to create a node
	log.Infof("Creating instance for %s...", d.MachineName)

	// TODO

	return nil
}

// PreCreateCheck validates parameters and checks if creation is possible
func (d *Driver) PreCreateCheck() error {
	// Implement any pre-creation checks
	// TODO

	return nil
}

// GetURL returns the URL of the docker daemon on the host
func (d *Driver) GetURL() (string, error) {
	// TODO
	return "", nil
}

// GetState returns the state of the host
func (d *Driver) GetState() (state.State, error) {
	// Here you would implement the API call to check the instance state
	// TODO

	return state.Running, nil
}

// Start starts the host
func (d *Driver) Start() error {
	// TODO: implement the API call to start the instance
	log.Infof("Starting instance %s")
	return nil
}

// Stop stops the host
func (d *Driver) Stop() error {
	// TODO: implement the API call to stop the instance
	log.Infof("Stopping instance %s", "")
	return nil
}

// Restart restarts the host
func (d *Driver) Restart() error {
	// TODO: implement the API call to restart the instance
	log.Infof("Restarting instance %s", "")
	return nil
}

// Kill forcefully stops the host
func (d *Driver) Kill() error {
	// TODO: implement the API call to force stop the instance
	log.Infof("Force stopping instance %s", "")
	return nil
}

// Remove removes the host
func (d *Driver) Remove() error {
	// TODO: implement the API call to delete the instance
	log.Infof("Removing instance %s", "")
	return nil
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	return driverName
}
