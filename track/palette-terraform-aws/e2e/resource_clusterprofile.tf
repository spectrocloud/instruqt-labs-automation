# If looking up a cluster profile instead of creating a new one
# data "spectrocloud_cluster_profile" "profile" {
#   # id = <uid>
#   name = var.cluster_profile_name
# }

# # Example of a Basic add-on profile
# resource "spectrocloud_cluster_profile" "cp-addon-aws" {
#   name        = var.cluster_profile_name
#   description = var.cluster_profile_description
#   cloud       = var.cloud
#   type        = var.type
#   pack {
#     name = "spectro-byo-manifest"
#     tag  = "1.0.x"
#     uid  = "5faad584f244cfe0b98cf489"
#     # layer  = ""
#     values = <<-EOT
#       manifests:
#         byo-manifest:
#           contents: |
#             # Add manifests here
#             apiVersion: v1
#             kind: Namespace
#             metadata:
#               labels:
#                 app: wordpress
#                 app3: wordpress3
#               name: wordpress
#     EOT
#   }
# }

data "spectrocloud_registry" "registry" {
  name = "Public Repo"
}

#data "spectrocloud_pack" "argo" {
  #name = "argo-cd"
  #registry_uid = data.spectrocloud_registry.registry.id
  #version  = "3.26.7"
#}

#data "spectrocloud_pack" "k8sdash" {
  #name = "k8s-dashboard"
  #registry_uid = data.spectrocloud_registry.registry.id
  #version  = "2.4.0"
#}

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


#data "spectrocloud_pack" "csi" {
  #name = "csi-aws"
  #registry_uid = data.spectrocloud_registry.registry.id
  #version  = "1.0.0"
#}

data "spectrocloud_pack" "csi" {
  name = "csi-aws"
  registry_uid = data.spectrocloud_registry.registry.id
  version  = "1.0.0"
}

#data "spectrocloud_pack" "cni" {
  #name    = "cni-calico"
  #registry_uid = data.spectrocloud_registry.registry.id
  #version = "3.22.0"
#}

data "spectrocloud_pack" "cni" {
  name    = "cni-calico"
  registry_uid = data.spectrocloud_registry.registry.id
  version = "3.24.1"
}

#data "spectrocloud_pack" "k8s" {
  #name    = "kubernetes"
  #registry_uid = data.spectrocloud_registry.registry.id
  #version = "1.21.10"
#}

data "spectrocloud_pack" "k8s" {
  name    = "kubernetes"
  registry_uid = data.spectrocloud_registry.registry.id
  version = "1.25.4"
}

#data "spectrocloud_pack" "ubuntu" {
  #name = "ubuntu-aws"
  #registry_uid = data.spectrocloud_registry.registry.id
  #version  = "20.04"
#}

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
    name   = "ubuntu-aws"
    tag    = data.spectrocloud_pack.ubuntu.version
    uid    = data.spectrocloud_pack.ubuntu.id
    values = data.spectrocloud_pack.ubuntu.values
  }

pack {
    name   = "kubernetes"
    tag    = data.spectrocloud_pack.k8s.version
    uid    = data.spectrocloud_pack.k8s.id
    values = data.spectrocloud_pack.k8s.values
  }

pack {
    name   = "cni-calico"
    tag    = data.spectrocloud_pack.cni.version
    uid    = data.spectrocloud_pack.cni.id
    values = data.spectrocloud_pack.cni.values
  }

pack {
    name   = "csi-aws"
    tag    = data.spectrocloud_pack.csi.version
    uid    = data.spectrocloud_pack.csi.id
    values = data.spectrocloud_pack.csi.values
  }

#pack {
    #name   = "open-policy-agent"
    #tag    = data.spectrocloud_pack.opa.version
    #uid    = data.spectrocloud_pack.opa.id
    #values = data.spectrocloud_pack.opa.values
  #}
  
pack {
    name   = "k8sdash"
    tag    = data.spectrocloud_pack.k8sdash.version
    uid    = data.spectrocloud_pack.k8sdash.id
    values = data.spectrocloud_pack.k8sdash.values
  }  

#pack {
    #name   = "argo-cd"
    #tag    = data.spectrocloud_pack.argo.version
    #uid    = data.spectrocloud_pack.argo.id
    #values = data.spectrocloud_pack.argo.values
  #}  

}
