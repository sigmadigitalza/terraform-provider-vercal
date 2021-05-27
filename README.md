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

The provider requires the following environmental variables:

| Variable | Required | Description |
| --- | :---: | --- |
| VERCEL_TOKEN | ✅ | A valid Vercel API token |
| VERCEL_TEAM_ID | - | A Vercel Team ID for working with a team rather than the token's user |

See the [example](./examples/main.tf) directory for an example usage.
