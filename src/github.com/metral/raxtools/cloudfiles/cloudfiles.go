package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/metral/raxtools/raxutils"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
)

var (
	push        = flag.Bool("push", false, "Push directory to CloudFiles")
	pull        = flag.Bool("pull", false, "Pull object from CloudFiles")
	contentPath = flag.String("contentPath", "", "Path for content to push to CloudFiles")
	container   = flag.String("container", "", "Container name")
	object      = flag.String("object", "", "Object name to pull from CloudFiles")
)

func checkFlags() {
	if !raxutils.FlagsSet() || !FlagsSet() {
		flag.Usage()
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func FlagsSet() bool {
	if *push &&
		(*contentPath == "" ||
			*container == "") {
		return false
	}
	if *pull &&
		(*container == "" ||
			*object == "") {
		return false
	}

	return true
}

func doPush(client *gophercloud.ServiceClient,
	contentPath, container, object string) (http.Header, error) {
	content, err := os.Open(contentPath)
	if err != nil {
		return nil, err
	}

	objectName := path.Base(contentPath)
	result := objects.Create(client, container, objectName, content, nil)
	hdr, err := result.ExtractHeader()
	if err != nil {
		return nil, err
	}

	return hdr, nil
}

func doPull(client *gophercloud.ServiceClient,
	container, object string) ([]byte, error) {
	result := objects.Download(client, container, object, nil)

	bytes, err := result.ExtractContent()
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func main() {
	flag.Parse()
	flag.Usage = usage
	checkFlags()

	// setup auth options
	authOpts := gophercloud.AuthOptions{
		IdentityEndpoint: *raxutils.IdentityEndpoint,
		Username:         *raxutils.Username,
		Password:         *raxutils.Password,
		TenantID:         *raxutils.TenantID,
	}

	// create client
	client, err := raxutils.NewObjectStorageClient(authOpts, "DFW")
	if err != nil {
		panic(err)
	}

	if *push {
		_, err := doPush(client, *contentPath, *container, *object)
		if err != nil {
			log.Printf("%v", err)
			os.Exit(1)
		}
	} else if *pull {
		bytes, err := doPull(client, *container, *object)
		if err != nil {
			log.Printf("%v", err)
			os.Exit(1)
		}
		ioutil.WriteFile(*object, bytes, 0644)
	}
}
