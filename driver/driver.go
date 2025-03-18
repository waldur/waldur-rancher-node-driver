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
	"github.com/google/uuid"
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
	FlavorUuid           string
	ImageUuid            string
	SystemVolumeSize     int
	SystemVolumeTypeUuid string
	DataVolumeTypeUuid   string
	SubnetUuids          []string
	SecurityGroupUuid    string
	ResourceUuid         string
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
		mcnflag.StringFlag{
			EnvVar: "WALDUR_OFFERING_UUID",
			Name:   "waldur-offering-uuid",
			Usage:  "UUID of the VM offering in Waldur",
		},
		mcnflag.StringFlag{
			EnvVar: "WALDUR_FLAVOR_UUID",
			Name:   "waldur-flavor-uuid",
			Usage:  "UUID of the VM flavor in Waldur",
		},
		mcnflag.StringFlag{
			EnvVar: "WALDUR_IMAGE_UUID",
			Name:   "waldur-image-uuid",
			Usage:  "UUID of the VM image in Waldur",
		},
		mcnflag.IntFlag{
			EnvVar: "WALDUR_SYS_VOLUME_SIZE",
			Name:   "waldur-sys-volume-size",
			Usage:  "System volume size for Waldur VM (GB)",
		},
		mcnflag.StringFlag{
			EnvVar: "WALDUR_SYS_VOLUME_TYPE_UUID",
			Name:   "waldur-sys-volume-type-uuid",
			Usage:  "UUID of the system volume type in Waldur",
		},
		mcnflag.StringFlag{
			EnvVar: "WALDUR_DATA_VOLUME_TYPE_UUID",
			Name:   "waldur-data-volume-type-uuid",
			Usage:  "UUID of the data volume type in Waldur",
		},
		mcnflag.StringFlag{
			EnvVar: "WALDUR_SEC_GROUP_UUID",
			Name:   "waldur-sec-group-uuid",
			Usage:  "UUID of the security group in Waldur",
		},
		mcnflag.StringSliceFlag{
			EnvVar: "WALDUR_SUBNET_UUIDS",
			Name:   "waldur-subnet-uuids",
			Usage:  "List of UUIDs of subnets in Waldur",
		},
	}
}

// SetConfigFromFlags configures the driver with the command line arguments
func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.ApiUrl = flags.String("waldur-api-url")
	d.ApiToken = flags.String("waldur-api-token")
	d.ProjectUuid = flags.String("waldur-proj-uuid")
	d.OfferingUuid = flags.String("waldur-offering-uuid")
	d.FlavorUuid = flags.String("waldur-flavor-uuid")
	d.ImageUuid = flags.String("waldur-image-uuid")
	d.SystemVolumeSize = flags.Int("waldur-sys-volume-size")
	d.SystemVolumeTypeUuid = flags.String("waldur-sys-volume-type-uuid")
	d.DataVolumeTypeUuid = flags.String("waldur-data-volume-type-uuid")
	d.SecurityGroupUuid = flags.String("waldur-sec-group-uuid")
	d.SubnetUuids = flags.StringSlice("waldur-subnet-uuids")

	// Validation
	if d.ApiUrl == "" {
		return fmt.Errorf("Waldur requires the --waldur-api-url option")
	}

	if d.ApiToken == "" {
		return fmt.Errorf("Waldur requires the --waldur-api-token option")
	}

	if d.ProjectUuid == "" {
		return fmt.Errorf("Waldur requires the --waldur-proj-uuid option")
	}
	if d.OfferingUuid == "" {
		return fmt.Errorf("Waldur requires the --waldur-offering-uuid option")
	}
	if d.FlavorUuid == "" {
		return fmt.Errorf("Waldur requires the --waldur-flavor-uuid option")
	}
	if d.ImageUuid == "" {
		return fmt.Errorf("Waldur requires the --waldur-image-uuid option")
	}
	if d.SystemVolumeSize == 0 {
		return fmt.Errorf("Waldur requires the --waldur-sys-volume-size to be greater than 5 GB")
	}
	if d.SystemVolumeTypeUuid == "" {
		return fmt.Errorf("Waldur requires the --waldur-sys-volume-type-uuid option")
	}
	if d.DataVolumeTypeUuid == "" {
		return fmt.Errorf("Waldur requires the --waldur-data-volume-type-uuid option")
	}
	if d.SecurityGroupUuid == "" {
		return fmt.Errorf("Waldur requires the --waldur-sec-group-uuid option")
	}
	if d.SubnetUuids == nil {
		d.SubnetUuids = []string{}
	}

	return nil
}

func (d *Driver) getWaldurClient() (*waldurclient.ClientWithResponses, error) {
	hc := http.Client{}
	auth, err := waldurclient.NewTokenAuth(d.ApiToken)
	if err != nil {
		log.Errorf("Error while creating token auth %s", err)
		return nil, err
	}

	client, err := waldurclient.NewClientWithResponses(d.ApiUrl, waldurclient.WithHTTPClient(&hc), waldurclient.WithRequestEditorFn(auth.Intercept))
	if err != nil {
		log.Errorf("Error creating Waldur client %s", err)
		return nil, err
	}

	return client, nil
}

func (d *Driver) getWaldurResource(client waldurclient.ClientWithResponses) (*waldurclient.Resource, error) {
	ctx := context.Background()
	resourceUuid, err := uuid.Parse(d.ResourceUuid)
	if err != nil {
		log.Errorf("Error converting resource UUID string to UUID object: %s", err)
		return nil, err
	}
	resp, err := client.MarketplaceResourcesRetrieveWithResponse(ctx, resourceUuid, &waldurclient.MarketplaceResourcesRetrieveParams{})

	if err != nil {
		log.Errorf("Error calling instance retrieval API: %v", err)
		return nil, err
	}

	if resp.StatusCode() != 200 {
		responseBody := string(resp.Body[:])
		log.Errorf("Unable to fetch the instance %s (%s), code %d, details", d.GetMachineName(), d.ResourceUuid, resp.StatusCode(), responseBody)
		msg := fmt.Sprintf("Unable to fetch the instance %s (%s), code %d", d.GetMachineName(), d.ResourceUuid, resp.StatusCode())
		return nil, errors.New(msg)
	}

	return resp.JSON200, nil
}

// Create creates a host in Waldur using the driver's config
func (d *Driver) Create() error {
	log.Infof("Creating instance for %s...", d.GetMachineName())

	projectUri := fmt.Sprintf("%s/api/projects/%s/", d.ApiUrl, d.ProjectUuid)
	offeringUri := fmt.Sprintf("%s/api/marketplace-public-offerings/%s/", d.ApiUrl, d.OfferingUuid)
	flavorUri := fmt.Sprintf("%s/api/openstack-flavors/%s/", d.ApiUrl, d.FlavorUuid)
	imageUri := fmt.Sprintf("%s/api/openstack-images/%s/", d.ApiUrl, d.ImageUuid)
	systemVolumeTypeUri := fmt.Sprintf("%s/api/openstack-volume-types/%s/", d.ApiUrl, d.SystemVolumeTypeUuid)
	dataVolumeTypeUri := fmt.Sprintf("%s/api/openstack-volume-types/%s/", d.ApiUrl, d.DataVolumeTypeUuid)
	subnets := make([]map[string]string, len(d.SubnetUuids))
	defaultSecGroupUri := fmt.Sprintf("%s/api/openstack-security-groups/%s/", d.ApiUrl, d.SecurityGroupUuid)
	securityGroups := make([]map[string]string, 1)
	securityGroups[0] = map[string]string{
		"url": defaultSecGroupUri,
	}

	for i, subnet := range d.SubnetUuids {
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

	client, err := d.getWaldurClient()
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
	}

	ctx := context.Background()
	resp, err := client.MarketplaceOrdersCreateWithResponse(ctx, payload)

	if err != nil {
		log.Errorf("Error calling API for instance creation: %v", err)
		return err
	}

	if resp.StatusCode() != 201 {
		responseBody := string(resp.Body[:])
		log.Errorf("Unable to create an instance %s, code %d, details", d.GetMachineName(), resp.StatusCode(), responseBody)
		msg := fmt.Sprintf("Unable to create an instance %s, code %d", d.GetMachineName(), resp.StatusCode())
		return errors.New(msg)
	}

	log.Infof("Successfully created instance %s", d.GetMachineName())
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
	url := fmt.Sprintf("%s/api/marketplace-resources/%s/", d.ApiUrl, d.ResourceUuid)
	return url, nil
}

// GetState returns the state of the host
func (d *Driver) GetState() (state.State, error) {
	// Here you would implement the API call to check the instance state
	client, err := d.getWaldurClient()
	if err != nil {
		log.Errorf("Error creating Waldur client %s", err)
		return state.None, err
	}

	resource, err := d.getWaldurResource(*client)
	if err != nil {
		return state.None, err
	}
	resourceStateStr := *resource.BackendMetadata.RuntimeState
	if resourceStateStr == "" {
		return state.None, nil
	}

	resourceState := waldurclient.CoreStates(resourceStateStr)

	resourceStateMap := map[waldurclient.CoreStates]state.State {
		"ACTIVE": state.Running,
		"BUILDING": state.Starting,
		"DELETED": state.Stopped,
		"SOFT_DELETED": state.Stopped,
		"ERROR": state.Error,
		"UNKNOWN": state.None,
		"HARD_REBOOT": state.Starting,
		"REBOOT": state.Starting,
		"REBUILD": state.Starting,
		"PAUSED": state.Paused,
		"SHUTOFF": state.Stopped,
		"STOPPED": state.Stopped,
		"SUSPENDED": state.Paused,
	}

	log.Infof("Successfully fetched instance, state %s", resourceState)
	return resourceStateMap[resourceState], nil
}

// Start starts the host
func (d *Driver) Start() error {
	log.Infof("Starting instance %s", d.GetMachineName())
	client, err := d.getWaldurClient()
	if err != nil {
		log.Errorf("Error creating Waldur client %s", err)
		return err
	}

	resource, err := d.getWaldurResource(*client)
	if err != nil {
		return err
	}

	ctx := context.Background()
	instanceResp, err := client.OpenstackInstancesStartWithResponse(ctx, *resource.ResourceUuid)

	if err != nil {
		log.Errorf("Error calling instance starting API: %v", err)
		return err
	}

	if instanceResp.StatusCode() != 202 {
		responseBody := string(instanceResp.Body[:])
		log.Errorf("Unable to start the instance %s (%s), code %d, details", d.GetMachineName(), d.ResourceUuid, instanceResp.StatusCode(), responseBody)
		msg := fmt.Sprintf("Unable to start the instance %s (%s), code %d", d.GetMachineName(), d.ResourceUuid, instanceResp.StatusCode())
		return errors.New(msg)
	}

	log.Infof("Successfully started the instance %s", d.GetMachineName())

	return nil
}

// Stop stops the host
func (d *Driver) Stop() error {
	log.Infof("Stopping instance %s", d.GetMachineName())
	client, err := d.getWaldurClient()
	if err != nil {
		log.Errorf("Error creating Waldur client %s", err)
		return err
	}

	resource, err := d.getWaldurResource(*client)
	if err != nil {
		return err
	}

	ctx := context.Background()
	instanceResp, err := client.OpenstackInstancesStopWithResponse(ctx, *resource.ResourceUuid)
	log.Infof("Req: %s", instanceResp.HTTPResponse.Request.URL)

	if err != nil {
		log.Errorf("Error calling instance stopping API: %v", err)
		return err
	}

	if instanceResp.StatusCode() != 202 {
		responseBody := string(instanceResp.Body[:])
		log.Errorf("Unable to stop the instance %s (%s), code %d, details", d.GetMachineName(), d.ResourceUuid, instanceResp.StatusCode(), responseBody)
		msg := fmt.Sprintf("Unable to stop the instance %s (%s), code %d", d.GetMachineName(), d.ResourceUuid, instanceResp.StatusCode())
		return errors.New(msg)
	}

	log.Infof("Successfully stopped the instance %s", d.GetMachineName())

	return nil
}

// Restart restarts the host
func (d *Driver) Restart() error {
	log.Infof("Restarting instance %s", d.GetMachineName())
	client, err := d.getWaldurClient()
	if err != nil {
		log.Errorf("Error creating Waldur client %s", err)
		return err
	}

	resource, err := d.getWaldurResource(*client)
	if err != nil {
		return err
	}

	ctx := context.Background()
	instanceResp, err := client.OpenstackInstancesRestartWithResponse(ctx, *resource.ResourceUuid)

	if err != nil {
		log.Errorf("Error calling instance restarting API: %v", err)
		return err
	}

	if instanceResp.StatusCode() != 202 {
		responseBody := string(instanceResp.Body[:])
		log.Errorf("Unable to restart the instance %s (%s), code %d, details", d.GetMachineName(), d.ResourceUuid, instanceResp.StatusCode(), responseBody)
		msg := fmt.Sprintf("Unable to restart the instance %s (%s), code %d", d.GetMachineName(), d.ResourceUuid, instanceResp.StatusCode())
		return errors.New(msg)
	}

	log.Infof("Successfully restarted the instance %s", d.GetMachineName())

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
