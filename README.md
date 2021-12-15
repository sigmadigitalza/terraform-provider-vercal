# Terraform Provider for Vercel

This is a Terraform provider which is used to configure Vercel.

## Requirements

* [Terraform](https://www.terraform.io/downloads.html) 0.14+
* [Go](https://golang.org/doc/install) 1.16.0 or higher

## Installing the provider

Enter the provider directory and run the following command:

```shell
make install
```

## Using the provider

The provider requires the following environment variables:

| Variable | Required | Description |
| --- | :---: | --- |
| VERCEL_TOKEN | ✅ | A valid Vercel API token |
| VERCEL_TEAM_ID | - | A Vercel Team ID for working with a team rather than the token's user |

See the [example](./examples/main.tf) directory for an example usage.

## Importing existing resources

Any IDs of existing resources required for importing with the following commands can be found using the
[Vercel API](https://vercel.com/docs/rest-api#endpoints/projects/find-a-project-by-id-or-name)

### Vercel projects

Use the following format to import a Vercel project:

```shell
terraform import vercel_project.test_project <vercel-project-name>
```

### Domains

Use the following format to import a domain:

```shell
terraform import vercel_project_domain.test_domain <vercel-project-name>:<domain-name>
```

### Environmental Variables

Use the following format to import an env:

```shell
terraform import vercel_project_env.test_env <vercel-project-name>:<env-id>
```
