package raxutils

import (
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

type AuthConfig struct {
	IdentityEndpoint string
	Username         string
	Password         string
	TenantID         string
}

func CreateToken(c *AuthConfig) (*tokens.Token, error) {
	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: c.IdentityEndpoint,
		Username:         c.Username,
		Password:         c.Password,
		TenantID:         c.TenantID,
	}

	// authenticate with provider
	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		return nil, err
	}

	// Create a new client to the provider
	client := openstack.NewIdentityV2(provider)

	// Create a new token to the provider
	opts := tokens.WrapOptions(authOpts)
	token, err := tokens.Create(client, opts).ExtractToken()
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
