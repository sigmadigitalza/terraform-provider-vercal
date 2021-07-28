package vercel

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vercel "github.com/sigmadigitalza/go-vercel-client"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext: resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type: schema.TypeString,
				Computed: true,
			},
			"name": {
				Type: schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"framework": {
				Type: schema.TypeString,
				Required: true,
			},
			"gitType": {
				Type: schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"gitRepo": {
				Type: schema.TypeString,
				Optional: false,
				ForceNew: true,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*vercel.Client)

	name := d.Get("name").(string)
	framework := d.Get("framework").(string)
	gitType := d.Get("gitType").(string)
	gitRepo := d.Get("gitRepo").(string)

	project, err := client.Project.CreateProject(ctx, name, framework, gitType, gitRepo)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(project.Id)

	resourceProjectRead(ctx, d, m)

	var diags diag.Diagnostics
	return diags
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*vercel.Client)
	name := d.Get("name").(string)

	project, err := client.Project.GetProject(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	var diags diag.Diagnostics
	return hydrateProject(diags, project, d)
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*vercel.Client)

	name := d.Get("name").(string)
	framework := d.Get("framework").(string)

	p := &vercel.Project{
		Name: name,
		Framework: framework,
	}

	_, err := client.Project.UpdateProject(ctx, name, p)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProjectRead(ctx, d, m)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*vercel.Client)
	name := d.Get("name").(string)

	err := client.Project.DeleteProject(ctx, name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	var diags diag.Diagnostics
	return diags
}

func hydrateProject(diags diag.Diagnostics, project *vercel.Project, d *schema.ResourceData) diag.Diagnostics {
	if err := d.Set("id", project.Id); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", project.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("framework", project.Framework); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
