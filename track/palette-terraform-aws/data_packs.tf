locals {
  os_version            = "1.0.0"
  k8s-1-1-27_version    = "1.27"
  k8s-1-1-28_version    = "1.28"
  k8s-1-1-29_version    = "1.29"
  cni_version           = "1.0"
  csi_version           = "1.30.0"
  opa_version           = "3.15.1"
  hipster_version       = "3.0.0"

}

data "spectrocloud_registry" "registry" {
  name = "Public Repo"
}

data "spectrocloud_pack" "csi" {
  name         = "csi-aws-ebs"
  registry_uid = data.spectrocloud_registry.registry.id
  version      = local.csi_version
}

data "spectrocloud_pack" "cni" {
  name         = "cni-aws-vpc-eks"
  registry_uid = data.spectrocloud_registry.registry.id
  version      = local.cni_version
}

data "spectrocloud_pack" "k8s-1-1-27" {
  name         = "kubernetes-eks"
  registry_uid = data.spectrocloud_registry.registry.id
  version      = local.k8s-1-1-27_version
}

data "spectrocloud_pack" "k8s-1-1-28" {
  name         = "kubernetes-eks"
  registry_uid = data.spectrocloud_registry.registry.id
  version      = local.k8s-1-1-28_version
}

data "spectrocloud_pack" "k8s-1-1-29" {
  name         = "kubernetes-eks"
  registry_uid = data.spectrocloud_registry.registry.id
  version      = local.k8s-1-1-29_version
}


data "spectrocloud_pack" "os" {
  name         = "amazon-linux-eks"
  registry_uid = data.spectrocloud_registry.registry.id
  version      = local.os_version
}

data "spectrocloud_pack" "opa" {
  name         = "open-policy-agent"
  registry_uid = data.spectrocloud_registry.registry.id
  version      = local.opa_version
}

data "spectrocloud_pack" "hipster" {
  name         = "sapp-hipster"
  registry_uid = data.spectrocloud_registry.registry.id
  version      = local.hipster_version
}