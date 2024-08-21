/* If using an existing cluster profile 

data "spectrocloud_cluster_profile" "addon_profile" {
  name = "<name-of-existing-cluster-profile>"
}
*/
resource "spectrocloud_cluster_eks" "cluster" {
  name = "tf-cluster01-eks"
  tags = ["cloud:eks","env:prod"]

  cluster_profile {
    id = spectrocloud_cluster_profile.infra_k8s-1-1-28.id
  }

  cluster_profile {
    id = spectrocloud_cluster_profile.add-on.id
  }

  cloud_account_id = data.spectrocloud_cloudaccount_aws.account.id

  cloud_config {
    ssh_key_name = var.aws_ssh_key_name
    region       = var.aws_region
    #vpc_id       = var.aws_vpc_id
    #azs          = var.azs != [] ? var.azs : null
    #az_subnets   = var.master_azs_subnets_map != {} ? var.master_azs_subnets_map : null
  }

  machine_pool {
    name          = "worker-pool01"
    count         = 3
    instance_type = "t3.large"
    #azs           = var.azs != [] ? var.azs : null
    #az_subnets    = var.master_azs_subnets_map != {} ? var.master_azs_subnets_map : null
    disk_size_gb  = 60
  }
}
