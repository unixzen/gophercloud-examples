package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/backups"
	"github.com/spf13/viper"
)

func readConfig(filename string) (string, string, string, string, string) {
	viper.SetConfigName(filename)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Config file not found")
	}
	IdentityEndpoint := viper.GetString("IdentityEndpoint")
	Username := viper.GetString("Username")
	Password := viper.GetString("Password")
	TenantID := viper.GetString("TenantID")
	DomainID := viper.GetString("DomainID")
	return IdentityEndpoint, Username, Password, TenantID, DomainID
}

func authOpenStack() (client *gophercloud.ServiceClient) {
	configFile := "backup"
	identityEndpoint, username, password, tenantID, domainID := readConfig(configFile)
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: identityEndpoint,
		Username:         username,
		Password:         password,
		TenantID:         tenantID,
		DomainID:         domainID,
	}

	provider, err := openstack.AuthenticatedClient(opts)

	if err != nil {
		log.Fatal(err)
	}

	client, err = openstack.NewBlockStorageV2(provider, gophercloud.EndpointOpts{
		Region: "regionOne",
	})

	if err != nil {
		log.Fatal(err)
	}

	return client
}

func main() {
	backupname := flag.String("backupname", "", "Backup name")
	volumeid := flag.String("volumeid", "", "Volume ID which will be backed up")
	incremental := flag.Bool("incremental", false, "Incremental backup: true or false")
	flag.Parse()

	createOpts := backups.CreateOpts{Name: *backupname, VolumeID: *volumeid, Force: true, Incremental: *incremental}

	backup, err := backups.Create(authOpenStack(), createOpts).Extract()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(backup.ID)
	fmt.Println(backup.Name)

}
