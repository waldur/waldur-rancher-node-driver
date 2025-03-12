package driver

import (
	"context"
	"errors"
	"fmt"

	"net/http"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/state"
	waldurclient "github.com/waldur/go-client"
)

const (
	driverName = "waldur"
)

type Driver struct {
	*drivers.BaseDriver

	ApiUrl               string
	ApiToken             string
	ProjectUuid          string
	OfferingUuid         string
	PlanUuid             string
	FlavorUuid           string
	ImageUuid            string
	SystemVolumeSize     int
	SystemVolumeTypeUuid string
	DataVolumeTypeUuid   string
	Subnets              []string
	SecurityGroupUuid    string
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

	d.Subnets = flags.StringSlice("waldur-subnets")

	return nil
}

// Create creates a host in Waldur using the driver's config
func (d *Driver) Create() error {
	log.Infof("Creating instance for %s...", d.MachineName)

	projectUri := fmt.Sprintf("%s/api/projects/%s/", d.ApiUrl, d.ProjectUuid)
	offeringUri := fmt.Sprintf("%s/api/marketplace-public-offerings/%s/", d.ApiUrl, d.OfferingUuid)
	flavorUri := fmt.Sprintf("%s/api/openstack-flavors/%s/", d.ApiUrl, d.FlavorUuid)
	imageUri := fmt.Sprintf("%s/api/openstack-images/%s/", d.ApiUrl, d.ImageUuid)
	systemVolumeTypeUri := fmt.Sprintf("%s/api/openstack-volume-types/%s/", d.ApiUrl, d.SystemVolumeTypeUuid)
	dataVolumeTypeUri := fmt.Sprintf("%s/api/openstack-volume-types/%s/", d.ApiUrl, d.DataVolumeTypeUuid)
	subnets := make([]map[string]string, len(d.Subnets))
	defaultSecGroupUri := fmt.Sprintf("%s/api/openstack-security-groups/%s/", d.ApiUrl, d.SecurityGroupUuid)
	securityGroups := make([]map[string]string, 1)
	securityGroups[0] = map[string]string{
		"url": defaultSecGroupUri,
	}

	for i, subnet := range d.Subnets {
		subnetUri := fmt.Sprintf("%s/api/openstack-subnets/%s/", d.ApiUrl, subnet)
		subnets[i] = map[string]string{
			"subnet": subnetUri,
		}
	}
	var attributes interface{} = map[string]interface{}{
		"name":               d.GetMachineName(),
		"flavor":             flavorUri,
		"image":              imageUri,
		"system_volume_size": d.SystemVolumeSize * 1024,
		"system_volume_type": systemVolumeTypeUri,
		"data_volume_type":   dataVolumeTypeUri,
		"ports":              subnets,
		"security_groups":    securityGroups,
		// TODO: add floating_ips
		// "floating_ips": floating_ips,
	}

	acceptingTermsOfService := true

	limits := map[string]int{}

	hc := http.Client{}
	auth, err := waldurclient.NewTokenAuth(d.ApiToken)
	if err != nil {
		log.Errorf("Error while creating token auth %s", err)
		return err
	}

	client, err := waldurclient.NewClientWithResponses(d.ApiUrl, waldurclient.WithHTTPClient(&hc), waldurclient.WithRequestEditorFn(auth.Intercept))
	if err != nil {
		log.Errorf("Error creating Waldur client %s", err)
		return err
	}
	requestType := waldurclient.RequestTypesCreate

	payload := waldurclient.MarketplaceOrdersCreateJSONRequestBody{
		AcceptingTermsOfService: &acceptingTermsOfService,
		Attributes:              &attributes,
		Limits:                  &limits,
		Offering:                offeringUri,
		Project:                 projectUri,
		Type:                    &requestType,
		// Plan:                    &planUri,
	}

	ctx := context.Background()
	resp, err := client.MarketplaceOrdersCreateWithResponse(ctx, payload)

	if err != nil {
		log.Errorf("Error calling API: %v", err)
		return err
	}

	if resp.StatusCode() != 201 {
		responseBody := string(resp.Body[:])
		log.Errorf("Unable to create an instance %s, code %d, details", d.MachineName, resp.StatusCode(), responseBody)
		msg := fmt.Sprintf("Unable to create an instance %s, code %d", d.MachineName, resp.StatusCode())
		return errors.New(msg)
	}

	log.Infof("Successfully created instance %s", d.MachineName)
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
