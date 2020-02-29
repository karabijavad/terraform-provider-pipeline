package main

import (
	"context"
	"log"
	"os"

	"github.com/antihax/optional"
	pipeline "github.com/banzaicloud/banzai-cli/.gen/pipeline"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBanzaiCloudAwsSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceBanzaiCloudAwsSecretCreate,
		Read:   resourceBanzaiCloudAwsSecretRead,
		Update: resourceBanzaiCloudAwsSecretUpdate,
		Delete: resourceBanzaiCloudAwsSecretDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"access_key_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"secret_access_key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"validate": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceBanzaiCloudAwsSecretCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(banzaiCloudProvider).client
	organizationID := m.(banzaiCloudProvider).organizationID

	out := &pipeline.CreateSecretRequest{
		Name: d.Get("name").(string),
		Type: "amazon",
	}
	out.Values = map[string]interface{}{}
	out.Values["AWS_ACCESS_KEY_ID"] = d.Get("access_key_id").(string)
	out.Values["AWS_SECRET_ACCESS_KEY"] = d.Get("secret_access_key").(string)

	response, _, err := client.SecretsApi.AddSecrets(
		context.Background(),
		organizationID,
		*out,
		&pipeline.AddSecretsOpts{
			Validate: optional.NewBool(d.Get("validate").(bool)),
		},
	)
	if err != nil {
		panic(err)
	}
	d.SetId(response.Id)

	return resourceBanzaiCloudAwsSecretRead(d, m)
}

type MissingSecretResponse struct {
	code    int
	message string
	error   string
}

func resourceBanzaiCloudAwsSecretRead(d *schema.ResourceData, m interface{}) error {
	f, err := os.OpenFile("terraform-provider-banzaicloud.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	client := m.(banzaiCloudProvider).client
	organizationID := m.(banzaiCloudProvider).organizationID
	secret, response, err := client.SecretsApi.GetSecret(context.Background(), organizationID, d.Id())
	if err != nil {
		if response.StatusCode == 400 {
			d.SetId("")
			return nil
		} else {
			panic(err)
		}
	}

	d.Set("name", secret.Name)
	d.Set("access_key_id", secret.Values["AWS_ACCESS_KEY_ID"])
	d.Set("secret_access_key", secret.Values["AWS_SECRET_ACCESS_KEY"])
	return nil
}

func resourceBanzaiCloudAwsSecretUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(banzaiCloudProvider).client
	organizationID := m.(banzaiCloudProvider).organizationID
	out := &pipeline.CreateSecretRequest{
		Name: d.Get("name").(string),
		Type: "amazon",
	}
	out.Values = map[string]interface{}{}
	out.Values["AWS_ACCESS_KEY_ID"] = d.Get("access_key_id").(string)
	out.Values["AWS_SECRET_ACCESS_KEY"] = d.Get("secret_access_key").(string)

	client.SecretsApi.UpdateSecrets(
		context.Background(),
		organizationID,
		d.Id(),
		*out,
		&pipeline.UpdateSecretsOpts{
			Validate: optional.NewBool(d.Get("validate").(bool)),
		},
	)

	return nil
}

func resourceBanzaiCloudAwsSecretDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(banzaiCloudProvider).client
	organizationID := m.(banzaiCloudProvider).organizationID
	_, err := client.SecretsApi.DeleteSecrets(context.Background(), organizationID, d.Id())
	if err != nil {
		panic(err)
	}
	return nil
}
