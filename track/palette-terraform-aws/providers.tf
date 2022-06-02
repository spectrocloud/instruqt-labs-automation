terraform {
  required_providers {
    spectrocloud = {
      version = ">= 0.1"
      source  = "spectrocloud/spectrocloud"
    }
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.0"
    }
  }
}




provider "spectrocloud" {
  host         = var.sc_host
  api_key      = var.sc_api_key
  project_name = var.sc_project_name
}

provider "aws" {
  region       = var.aws_region
  access_key   = var.aws_access_key_id
  secret_key   = var.aws_secret_access_key
}
  
