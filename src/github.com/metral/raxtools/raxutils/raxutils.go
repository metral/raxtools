package raxutils

import (
	"errors"
	"flag"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tokens"
)

var (
	IdentityEndpoint = flag.String("authURL", "", "Openstack Identity Endpoint")
	Username         = flag.String("username", "", "Openstack Username")
	Password         = flag.String("password", "", "Openstack Password")
	TenantID         = flag.String("tenantID", "", "Openstack Tenant ID")
)

func NewIdentityClient(a gophercloud.AuthOptions) (*gophercloud.ServiceClient, error) {
	// authenticate with provider
	provider, err := openstack.AuthenticatedClient(a)
	if err != nil {
		return nil, err
	}

	// Create a new client to the provider
	client := openstack.NewIdentityV2(provider)
	if client == nil {
		return nil, errors.New("Could not create new identity client")
	}

	return client, nil
}

func NewObjectStorageClient(a gophercloud.AuthOptions, region string) (*gophercloud.ServiceClient, error) {
	// authenticate with provider
	provider, err := openstack.AuthenticatedClient(a)
	if err != nil {
		return nil, err
	}

	// Create a new client to the provider
	client, err := openstack.NewObjectStorageV1(provider, gophercloud.EndpointOpts{
		Region: region,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func CreateToken(a gophercloud.AuthOptions, c *gophercloud.ServiceClient) (*tokens.Token, error) {
	opts := tokens.WrapOptions(a)

	// Create a new token
	token, err := tokens.Create(c, opts).ExtractToken()
	if err != nil {
		return nil, err
	}

	return token, nil
}

func FlagsSet() bool {
	if *IdentityEndpoint == "" ||
		*Username == "" ||
		*Password == "" ||
		*TenantID == "" {
		return false
	}

	return true
}
