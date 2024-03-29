data "spectrocloud_pack" "proxy" {
  name = "spectro_proxy"
  version  = "1.3.0"
}

data "spectrocloud_pack" "csi" {
  name = "csi-gcp-driver"
  version  = "1.8.2"
}

data "spectrocloud_pack" "cni" {
  name    = "cni-calico"
  version = "3.25.1"
}

data "spectrocloud_pack" "k8s" {
  name    = "kubernetes"
  version = "1.25.9"
}

data "spectrocloud_pack" "ubuntu" {
  name = "ubuntu-gcp"
  version  = "22.04"
}

resource "spectrocloud_cluster_profile" "profile" {
  name        = var.cluster_profile_name
  description = var.cluster_profile_description
  tags        = var.cluster_profile_tags
  cloud       = var.cloud
  type        = var.type


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
    values = file("${path.cwd}/values_yaml/spectro_proxy.yaml")
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
    name   = data.spectrocloud_pack.spectro_proxy.name
    tag    = data.spectrocloud_pack.spectro_proxy.version
    uid    = data.spectrocloud_pack.spectro_proxy.id
    values = data.spectrocloud_pack.spectro_proxy.values
  }

}