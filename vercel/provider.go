package vercel

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vercel "github.com/sigmadigitalza/go-vercel-client/v2"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: configureContext,
		ResourcesMap: map[string]*schema.Resource{
			"vercel_project":        resourceProject(),
			"vercel_project_env":    resourceProjectEnv(),
			"vercel_project_domain": resourceProjectDomain(),
		},
	}
}

func configureContext(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	client, err := vercel.New()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}
