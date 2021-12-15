package vercel

import (
	"context"
	"errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vercel "github.com/sigmadigitalza/go-vercel-client/v2"
	"strings"
)

var (
	ProjectDomainNotFoundError = errors.New("project domain not found")
	InvalidDomainIdError       = errors.New("invalid domain ID specified")
)

func importStateProjectDomainContext(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	// for importing of domains, specify the id in the following format: "vercel-project-name:domain-name"
	if strings.Contains(id, ":") {
		values := strings.Split(id, ":")

		if len(values) != 2 {
			return nil, InvalidDomainIdError
		}

		err := d.Set("name", values[0])
		if err != nil {
			return nil, err
		}

		err = d.Set("domain", values[1])
		if err != nil {
			return nil, err
		}

		d.SetId(values[1])
	}

	return []*schema.ResourceData{d}, nil
}

func resourceProjectDomain() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectDomainCreate,
		ReadContext:   resourceProjectDomainRead,
		UpdateContext: resourceProjectDomainUpdate,
		DeleteContext: resourceProjectDomainDestroy,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"redirect": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: importStateProjectDomainContext,
		},
	}
}

func resourceProjectDomainCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*vercel.Client)

	name := d.Get("name").(string)
	domain := d.Get("domain").(string)
	redirect := d.Get("redirect").(string)

	_, err := client.Project.AddDomain(ctx, name, domain, redirect)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(domain)

	resourceProjectDomainRead(ctx, d, m)

	return diag.Diagnostics{}
}

func resourceProjectDomainRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*vercel.Client)

	name := d.Get("name").(string)
	domain := d.Get("domain").(string)

	project, err := client.Project.GetProject(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	projectDomain, err := findProjectDomain(domain, project.Alias)
	if err == ProjectDomainNotFoundError {
		d.SetId("")
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return hydrateProjectDomain(diag.Diagnostics{}, projectDomain, d)
}

func resourceProjectDomainUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*vercel.Client)

	name := d.Get("name").(string)
	domain := d.Get("domain").(string)
	redirect := d.Get("redirect").(string)

	_, err := client.Project.UpdateDomain(ctx, name, domain, redirect)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProjectDomainRead(ctx, d, m)
}

func resourceProjectDomainDestroy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*vercel.Client)

	name := d.Get("name").(string)
	domain := d.Get("domain").(string)

	_, err := client.Project.DeleteDomain(ctx, name, domain)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func findProjectDomain(domain string, domains []*vercel.Domain) (*vercel.Domain, error) {
	for i := 0; i < len(domains); i++ {
		current := domains[i]
		if current.Domain == domain {
			return current, nil
		}
	}

	return nil, ProjectDomainNotFoundError
}

func hydrateProjectDomain(diags diag.Diagnostics, domain *vercel.Domain, d *schema.ResourceData) diag.Diagnostics {
	if err := d.Set("domain", domain.Domain); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("redirect", domain.Redirect); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
