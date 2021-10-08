terraform {
  required_providers {
    vercel = {
      version = "1.0.0"
      source = "sigmadigital.io/vercel/vercel"
    }
  }
}

provider "vercel" {}

resource "vercel_project" "test_project" {
  name = "test-project"
  framework = "nextjs"
}

resource "vercel_project_env" "test_env" {
  name = vercel_project.test_project.name
  type = "encrypted"
  key = "TEST_ENV"
  value = "secret-value"
  target = [ "production" ]
}

resource "vercel_project_domain" "test_domain" {
  name = vercel_project.test_project.name
  domain = "vercel.com"
}
