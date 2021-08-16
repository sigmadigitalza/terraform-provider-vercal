package vercel

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vercel "github.com/sigmadigitalza/go-vercel-client/v2"
)

func resourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"framework": {
				Type:     schema.TypeString,
				Required: true,
			},
			"root_directory": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"git_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"git_repo": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"build_command": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"output_directory": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"command_for_ignoring_build_step": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*vercel.Client)

	name := d.Get("name").(string)
	framework := d.Get("framework").(string)
	rootDirectory := d.Get("root_directory").(string)
	gitType := d.Get("git_type").(string)
	gitRepo := d.Get("git_repo").(string)
	buildCommand := d.Get("build_command").(string)
	outputDirectory := d.Get("output_directory").(string)

	options := &vercel.CreateProjectOptions{
		Name:            name,
		Framework:       framework,
		RepositoryType:  gitType,
		RepositoryName:  gitRepo,
		RootDirectory:   rootDirectory,
		BuildCommand:    buildCommand,
		OutputDirectory: outputDirectory,
	}

	project, err := client.Project.CreateProject(ctx, options)
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
	rootDirectory := d.Get("root_directory").(string)
	buildCommand := d.Get("build_command").(string)
	outputDirectory := d.Get("output_directory").(string)
	commandForIgnoringBuildStep := d.Get("command_for_ignoring_build_step").(string)

	p := &vercel.Project{
		Name:                        name,
		Framework:                   framework,
		RootDirectory:               rootDirectory,
		BuildCommand:                buildCommand,
		OutputDirectory:             outputDirectory,
		CommandForIgnoringBuildStep: commandForIgnoringBuildStep,
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

	if project.RootDirectory != "" {
		if err := d.Set("root_directory", project.RootDirectory); err != nil {
			return diag.FromErr(err)
		}
	}

	if project.Link != nil && project.Link.Org != "" && project.Link.Repo != "" {
		if err := d.Set("git_type", project.Link.Type); err != nil {
			return diag.FromErr(err)
		}

		org := project.Link.Org
		repo := project.Link.Repo
		if err := d.Set("git_repo", fmt.Sprintf("%s/%s", org, repo)); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("build_command", project.BuildCommand); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("output_directory", project.OutputDirectory); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("command_for_ignoring_build_step", project.CommandForIgnoringBuildStep); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
