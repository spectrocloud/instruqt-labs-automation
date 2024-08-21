# If looking up a cluster profile instead of creating a new one
# data "spectrocloud_cluster_profile" "profile" {
#   # id = <uid>
#   name = "eks-basic"
# }
#

resource "spectrocloud_cluster_profile" "infra_k8s-1-1-29" {
  name        = "tf-infra-eks-profile"
  description = "basic eks cp"
  cloud       = "eks"
  type        = "infra"
  tags        = ["cloud:eks"]
  version     = "1.1.29"

  pack {
    name   = data.spectrocloud_pack.os.name
    tag    = local.os_version
    uid    = data.spectrocloud_pack.os.id
    values = data.spectrocloud_pack.os.values
  }


pack {
    name   = data.spectrocloud_pack.k8s-1-1-29.name
    tag    = local.k8s-1-1-29_version
    uid    = data.spectrocloud_pack.k8s-1-1-29.id
    values = data.spectrocloud_pack.k8s-1-1-29.values
  }

  pack {
    name   = data.spectrocloud_pack.cni.name
    tag    = local.cni_version
    uid    = data.spectrocloud_pack.cni.id
    values = data.spectrocloud_pack.cni.values
  }

  pack {
    name   = data.spectrocloud_pack.csi.name
    tag    = local.csi_version
    uid    = data.spectrocloud_pack.csi.id
    values = data.spectrocloud_pack.csi.values
  }

}

