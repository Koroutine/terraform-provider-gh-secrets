package main

import (
	"context"

	"github.com/google/go-github/v39/github"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/oauth2"
)

func Provider() *schema.Provider {
	return &schema.Provider{

		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("GITHUB_TOKEN", nil),
				Description: "The OAuth token used to connect to GitHub. Anonymous mode is enabled if both `token` and `app_auth` are not set.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"gh-secrets_value": resourceSecret(),
		},
		ConfigureContextFunc: config,
	}

}

func config(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	token := d.Get("token").(string)

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &Meta{
		client: client,
	}, nil

}
