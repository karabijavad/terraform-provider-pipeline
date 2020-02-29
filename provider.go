package main

import (
	pipeline "github.com/banzaicloud/banzai-cli/.gen/pipeline"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"golang.org/x/oauth2"
)

type banzaiCloudProvider struct {
	client         *pipeline.APIClient
	organizationID int32
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	httpClient := oauth2.NewClient(nil, oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: d.Get("access_token").(string),
		},
	))
	client := pipeline.NewAPIClient(&pipeline.Configuration{
		BasePath:      d.Get("api_url").(string),
		DefaultHeader: make(map[string]string),
		UserAgent:     "terraform-provider-banzaicloud",
		Debug:         true,
		HTTPClient:    httpClient,
	})
	return banzaiCloudProvider{
		client:         client,
		organizationID: int32(d.Get("organization_id").(int)),
	}, nil
}

// Provider main
func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureFunc: configureProvider,
		Schema: map[string]*schema.Schema{
			"access_token": {
				Type:     schema.TypeString,
				Required: true,
			},
			"api_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"banzaicloud_aws_secret": resourceBanzaiCloudAwsSecret(),
		},
	}
}
