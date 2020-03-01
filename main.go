package main

import (
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	log "github.com/sirupsen/logrus"
)

func init() {
	f, err := os.OpenFile("terraform-provider-pipeline.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
	} else {
		log.SetOutput(f)
	}

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(f)
	log.SetLevel(log.DebugLevel)
}

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return Provider()
		},
	})
}
