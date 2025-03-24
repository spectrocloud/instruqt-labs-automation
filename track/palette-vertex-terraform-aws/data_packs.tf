locals {
  os_version            = "1.0.0"
  k8s-1-1-30_version    = "1.30"
  k8s-1-1-31_version    = "1.31"
  k8s-1-1-32_version    = "1.32"
  cni_version           = "1.1.17"
  csi_version           = "1.35.0"
  opa_version           = "3.18.1"
  hipster_version       = "3.0.0"

}

data "spectrocloud_registry" "registry" {
  name = "Prod-palette"
}

data "spectrocloud_registry" "fips-registry" {
  name = "Palette Registry"
}

data "spectrocloud_pack" "csi" {
  name         = "csi-aws-ebs"
  registry_uid = data.spectrocloud_registry.fips-registry.id
  version      = local.csi_version
}

data "spectrocloud_pack" "cni" {
  name         = "cni-aws-vpc-eks-helm-fips"
  registry_uid = data.spectrocloud_registry.fips-registry.id
  version      = local.cni_version
}

data "spectrocloud_pack" "k8s-1-1-30" {
  name         = "kubernetes-eks"
  registry_uid = data.spectrocloud_registry.fips-registry.id
  version      = local.k8s-1-1-30_version
}

data "spectrocloud_pack" "k8s-1-1-31" {
  name         = "kubernetes-eks"
  registry_uid = data.spectrocloud_registry.fips-registry.id
  version      = local.k8s-1-1-31_version
}

data "spectrocloud_pack" "k8s-1-1-32" {
  name         = "kubernetes-eks"
  registry_uid = data.spectrocloud_registry.fips-registry.id
  version      = local.k8s-1-1-32_version
}


data "spectrocloud_pack" "os" {
  name         = "amazon-linux-eks"
  registry_uid = data.spectrocloud_registry.fips-registry.id
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