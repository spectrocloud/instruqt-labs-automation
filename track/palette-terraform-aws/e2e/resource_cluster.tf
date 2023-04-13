resource "spectrocloud_cluster_aws" "cluster" {
  name               = var.sc_cluster_name
  cluster_profile {
   id = spectrocloud_cluster_profile.profile.id
  }

  #for newly created cloud account
  cloud_account_id   = spectrocloud_cloudaccount_aws.account.id
  
  #for existing cloud account
  #cloud_account_id   = data.spectrocloud_cloudaccount_aws.account.id
  
  cloud_config {
    ssh_key_name = var.aws_ssh_key_name
    region       = var.aws_region
  }

  machine_pool {
    control_plane           = true
    control_plane_as_worker = true
    name                    = "master-pool"
    count                   = var.master-pool_node_count
    instance_type           = var.master_instance_type
    disk_size_gb            = 62
    azs                     = var.master_aws_region_az
  }

  machine_pool {
    name          = "worker-pool01"
    count         = var.worker-pool_node_count
    instance_type = var.worker_instance_type
    disk_size_gb  = 62
    azs           = var.worker_aws_region_az
  }

}
