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
	EnvNotFoundError = errors.New("project env not found")
	InvalidEnvIdError = errors.New("invalid env ID specified")
)

func importStateProjectEnvContext(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()

	// for importing of envs, specify the id in the following format: "vercel-project-name:vercel-env-id"
	if strings.Contains(id, ":") {
		values := strings.Split(id, ":")

		if len(values) != 2 {
			return nil, InvalidEnvIdError
		}

		err := d.Set("name", values[0])
		if err != nil {
			return nil, err
		}

		d.SetId(values[1])
	}

	return []*schema.ResourceData{d}, nil
}

func resourceProjectEnv() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectEnvCreate,
		ReadContext:   resourceProjectEnvRead,
		UpdateContext: resourceProjectEnvUpdate,
		DeleteContext: resourceProjectEnvDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type: schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: importStateProjectEnvContext,
		},
	}
}

func resourceProjectEnvCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*vercel.Client)

	name := d.Get("name").(string)
	envType := d.Get("type").(string)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	target := d.Get("target").([]interface{})

	projectEnv, err := client.Project.CreateProjectEnv(ctx, name, envType, key, value, extractStringSlice(target))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(projectEnv.Id)

	resourceProjectEnvRead(ctx, d, m)

	return diag.Diagnostics{}
}

func resourceProjectEnvRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*vercel.Client)

	name := d.Get("name").(string)
	id := d.Id()

	projectEnvs, err := client.Project.GetProjectEnvs(ctx, name, true)
	if err != nil {
		return diag.FromErr(err)
	}

	env, err := findEnv(id, projectEnvs)
	if err == EnvNotFoundError {
		d.SetId("")
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return hydrateProjectEnv(diag.Diagnostics{}, env, d)
}

func resourceProjectEnvUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*vercel.Client)

	id := d.Get("id").(string)
	name := d.Get("name").(string)
	envType := d.Get("type").(string)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	target := d.Get("target").([]interface{})

	env := &vercel.ProjectEnv{
		Id:     id,
		Type:   envType,
		Key:    key,
		Value:  value,
		Target: extractStringSlice(target),
	}

	_, err := client.Project.EditProjectEnv(ctx, name, env)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceProjectEnvRead(ctx, d, m)
}

func resourceProjectEnvDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*vercel.Client)

	name := d.Get("name").(string)
	id := d.Get("id").(string)

	err := client.Project.DeleteProjectEnv(ctx, name, id)
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.Diagnostics{}
}

func findEnv(envId string, envs []*vercel.ProjectEnv) (*vercel.ProjectEnv, error) {
	for i := 0; i < len(envs); i++ {
		current := envs[i]
		if current.Id == envId {
			return current, nil
		}
	}

	return nil, EnvNotFoundError
}

func extractStringSlice(getResult []interface{}) []string {
	slice := make([]string, len(getResult))

	for i := 0; i < len(getResult); i++ {
		slice[i] = getResult[i].(string)
	}

	return slice
}

func hydrateProjectEnv(diags diag.Diagnostics, env *vercel.ProjectEnv, d *schema.ResourceData) diag.Diagnostics {
	if err := d.Set("type", env.Type); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("key", env.Key); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("value", env.Value); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("target", env.Target); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
