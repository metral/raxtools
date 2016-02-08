package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/identity/v2/tokens"
)

var (
	identityEndpoint = flag.String("authURL", "", "Openstack Identity Endpoint")
	username         = flag.String("username", "", "Openstack Username")
	password         = flag.String("password", "", "Openstack Password")
	tenantID         = flag.String("tenantID", "", "Openstack Tenant ID")
)

type AuthConfig struct {
	OSIdentityEndpoint string
	OSUsername         string
	OSPassword         string
	OSTenantID         string
}

func CreateToken(c *AuthConfig) (*tokens.Token, error) {
	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: c.OSIdentityEndpoint,
		Username:         c.OSUsername,
		Password:         c.OSPassword,
		TenantID:         c.OSTenantID,
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

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func checkFlags() {
	if *identityEndpoint == "" ||
		*username == "" ||
		*password == "" ||
		*tenantID == "" {
		flag.Usage()
	}
}

func main() {
	flag.Usage = usage
	flag.Parse()
	checkFlags()

	config := AuthConfig{
		OSIdentityEndpoint: *identityEndpoint,
		OSUsername:         *username,
		OSPassword:         *password,
		OSTenantID:         *tenantID,
	}

	token, err := CreateToken(&config)
	if err != nil {
		panic(err)
	}

	log.Printf("Token: %s\n", token.ID)
}
