
data "spectrocloud_registry" "registry" {
  name = "Public Repo"
}

data "spectrocloud_pack" "argo" {
  name = "argo-cd"
  registry_uid = data.spectrocloud_registry.registry.id
  version  = "3.26.7"
}

data "spectrocloud_pack" "k8sdash" {
  name = "spectro-k8s-dashboard"
  registry_uid = data.spectrocloud_registry.registry.id
  version  = "2.7.0"
}

#data "spectrocloud_pack" "opa" {
  #name = "open-policy-agent"
  #registry_uid = data.spectrocloud_registry.registry.id
  #version  = "3.7.0"
#}

data "spectrocloud_pack" "csi" {
  name = "csi-aws-ebs"
  registry_uid = data.spectrocloud_registry.registry.id
  version  = "1.16.0"
}

data "spectrocloud_pack" "cni" {
  name    = "cni-calico"
  registry_uid = data.spectrocloud_registry.registry.id
  version = "3.24.1"
}

data "spectrocloud_pack" "k8s" {
  name    = "kubernetes"
  registry_uid = data.spectrocloud_registry.registry.id
  version = "1.25.4"
}

data "spectrocloud_pack" "ubuntu" {
  name = "ubuntu-aws"
  registry_uid = data.spectrocloud_registry.registry.id
  version  = "22.04"
}

resource "spectrocloud_cluster_profile" "profile" {
  name        = var.sc_cp_profile_name
  description = var.sc_cp_profile_description
  tags        = var.sc_cp_profile_tags
  cloud       = var.sc_cp_cloud
  type        = var.sc_cp_type

pack {
    name   = data.spectrocloud_pack.ubuntu.name
    tag    = data.spectrocloud_pack.ubuntu.version
    uid    = data.spectrocloud_pack.ubuntu.id
    values = data.spectrocloud_pack.ubuntu.values
  }

pack {
    name   = data.spectrocloud_pack.k8s.name
    tag    = data.spectrocloud_pack.k8s.version
    uid    = data.spectrocloud_pack.k8s.id
    values = data.spectrocloud_pack.k8s.values
  }

pack {
    name   = data.spectrocloud_pack.cni.name
    tag    = data.spectrocloud_pack.cni.version
    uid    = data.spectrocloud_pack.cni.id
    values = data.spectrocloud_pack.cni.values
  }

pack {
    name   = data.spectrocloud_pack.csi.name
    tag    = data.spectrocloud_pack.csi.version
    uid    = data.spectrocloud_pack.csi.id
    values = data.spectrocloud_pack.csi.values
  }

pack {
    name   = data.spectrocloud_pack.k8sdash.name
    tag    = data.spectrocloud_pack.k8sdash.version
    uid    = data.spectrocloud_pack.k8sdash.id
    values = data.spectrocloud_pack.k8sdash.values
  }  

pack {
    name   = "argo-cd"
    tag    = data.spectrocloud_pack.argo.version
    uid    = data.spectrocloud_pack.argo.id
    values = data.spectrocloud_pack.argo.values
  }  

}
