package main

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/backups"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/spf13/viper"
	"log"
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
	listOpts := backups.ListOpts{}

	pager := backups.List(authOpenStack(), listOpts)

	pager.EachPage(func(page pagination.Page) (bool, error) {
		bList, err := backups.ExtractBackups(page)
		if err != nil {
			return false, err
		}

		for _, b := range bList {
			backup, err := backups.Get(authOpenStack(), b.ID).Extract()
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%v, %v \n", backup.ID, backup.VolumeID)

		}

		return true, nil
	})
}
