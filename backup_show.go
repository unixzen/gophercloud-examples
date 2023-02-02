package main

import (
	"flag"
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/backups"
	"log"
	"github.com/spf13/viper"
	"time"
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

	backupid := flag.String("backupid", "", "Backup ID")
	flag.Parse()
	backup, err := backups.Get(authOpenStack(), *backupid).Extract()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Backup ID: ", backup.ID)
	fmt.Println("Backup Name: ", backup.Name)
	fmt.Println("Backup Volume: ", backup.VolumeID)
	fmt.Println("Backup Status: ", backup.Status)
	fmt.Println("Backup CreatedAt: ", backup.CreatedAt)
	fmt.Println("=====================================")
	timeNow := time.Now()
	periodCalculate := timeNow.Sub(backup.CreatedAt)
	periodStore, _ := time.ParseDuration("168h")
	if periodCalculate > periodStore {
		fmt.Println(periodCalculate)
		fmt.Println(periodStore)
		fmt.Println("Need to backup volume")
	}

}
