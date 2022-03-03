# Basic Cluster Profile creation

Example of creating a cluster profile in Palette. This terraform configuration
will provision the following resources on Spectro Cloud:
- Create a Palette Edge Cluster Profile
- Use customized manifes files for K8s and Spectro Proxy 

## Instructions:

Clone this repository to a local directory, and then change directory to `cluster_profile`. Proceed with the following:
1. Follow the Spectro Cloud documentation to create Palette trial
[Spectro Cloud Palette Trial ](https://www.spectrocloud.com/free-trial/).
2. From the current directory, copy the template variable file `terraform.tfvars.example` to a new file `terraform.tfvars`.
3. Specify and update all the placeholder values in the `terraform.tfvars` file.
4. Initialize and run terraform: `terraform init && terraform apply`.
5. Wait for the cluster profile creation to finish.

## Cleanup:

Run the destroy operation:

```shell
terraform destroy
```
