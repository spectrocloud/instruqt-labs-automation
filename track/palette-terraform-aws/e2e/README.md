# Basic Cluster demo

End-to-end example of provisioning a new AWS K8s cluster with all of its dependencies. This terraform configuration
will provision the following resources on Spectro Cloud:
- AWS subscription key pair
- Configure Palette to use an AWS Cloud Account
- Create a Palette AWS Cluster Profile
- Use Palette to deploy a Kubernetes cluster into your AWS Subscription using ec2 instances

## Instructions:

Clone this repository to a local directory, and then change directory to `aws_iaas`. Proceed with the following:
1. Follow the Spectro Cloud documentation to create an AWS cloud account with appropriate permissions:
[AWS Cloud Account](https://docs.spectrocloud.com/clusters/?clusterType=aws_cluster#awscloudaccountpermissions).
2. From the current directory, copy the template variable file `terraform.tfvars.example` to a new file `terraform.tfvars`.
3. Specify and update all the placeholder values in the `terraform.tfvars` file.
4. Initialize and run terraform: `terraform init && terraform apply`.
5. Wait for the cluster creation to finish.

Once the cluster is provisioned, the cluster _kubeconfig_ file is exported in the current working directly.

Export the kubeconfig and check cluster pods:

```shell
export KUBECONFIG=kubeconfig_<yourClusterName>
kubectl get pod -A
```

## Cleanup:

Run the destroy operation:

```shell
terraform destroy
```
