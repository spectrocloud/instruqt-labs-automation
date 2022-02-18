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


data "spectrocloud_pack" "proxy" {
  name = "spectro-proxy"
  version  = "1.0.0"
}

data "spectrocloud_pack" "csi" {
  name = "csi-rook-ceph"
  version  = "1.8.0"
}

data "spectrocloud_pack" "cni" {
  name    = "cni-calico"
  version = "3.19.0"
}

data "spectrocloud_pack" "k8s" {
  name    = "kubernetes"
  version = "1.21.5"
}

data "spectrocloud_pack" "ubuntu" {
  name = "ubuntu-edge"
  version  = "20.04"
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
  }
pack {
    name   = data.spectrocloud_pack.k8s.name
    tag    = data.spectrocloud_pack.k8s.version
    uid    = data.spectrocloud_pack.k8s.id
    values = <<-EOT
pack:
  k8sHardening: True
  #CIDR Range for Pods in cluster
  # Note : This must not overlap with any of the host or service network
  podCIDR: "192.168.0.0/16"
  #CIDR notation IP range from which to assign service cluster IPs
  # Note : This must not overlap with any IP ranges assigned to nodes for pods.
  serviceClusterIpRange: "10.96.0.0/12"

# KubeAdm customization for kubernetes hardening. Below config will be ignored if k8sHardening property above is disabled
kubeadmconfig:
  apiServer:
    certSANs:
    - "cluster-{{ .spectro.system.cluster.uid }}.{{ .spectro.system.reverseproxy.server }}"
    extraArgs:
      # Note : secure-port flag is used during kubeadm init. Do not change this flag on a running cluster
      secure-port: "6443"
      anonymous-auth: "true"
      insecure-port: "0"
      profiling: "false"
      disable-admission-plugins: "AlwaysAdmit"
      default-not-ready-toleration-seconds: "60"
      default-unreachable-toleration-seconds: "60"
      enable-admission-plugins: "AlwaysPullImages,NamespaceLifecycle,ServiceAccount,NodeRestriction,PodSecurityPolicy"
      audit-log-path: /var/log/apiserver/audit.log
      audit-policy-file: /etc/kubernetes/audit-policy.yaml
      audit-log-maxage: "30"
      audit-log-maxbackup: "10"
      audit-log-maxsize: "100"
      authorization-mode: RBAC,Node
      tls-cipher-suites: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256"
    extraVolumes:
      - name: audit-log
        hostPath: /var/log/apiserver
        mountPath: /var/log/apiserver
        pathType: DirectoryOrCreate
      - name: audit-policy
        hostPath: /etc/kubernetes/audit-policy.yaml
        mountPath: /etc/kubernetes/audit-policy.yaml
        readOnly: true
        pathType: File
  controllerManager:
    extraArgs:
      profiling: "false"
      terminated-pod-gc-threshold: "25"
      pod-eviction-timeout: "1m0s"
      use-service-account-credentials: "true"
      feature-gates: "RotateKubeletServerCertificate=true"
  scheduler:
    extraArgs:
      profiling: "false"
  kubeletExtraArgs:
    read-only-port : "0"
    event-qps: "0"
    feature-gates: "RotateKubeletServerCertificate=true"
    protect-kernel-defaults: "true"
    tls-cipher-suites: "TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_256_GCM_SHA384,TLS_RSA_WITH_AES_128_GCM_SHA256"
  files:
    - path: hardening/audit-policy.yaml
      targetPath: /etc/kubernetes/audit-policy.yaml
      targetOwner: "root:root"
      targetPermissions: "0600"
    - path: hardening/privileged-psp.yaml
      targetPath: /etc/kubernetes/hardening/privileged-psp.yaml
      targetOwner: "root:root"
      targetPermissions: "0600"
    - path: hardening/90-kubelet.conf
      targetPath: /etc/sysctl.d/90-kubelet.conf
      targetOwner: "root:root"
      targetPermissions: "0600"
  preKubeadmCommands:
    # For enabling 'protect-kernel-defaults' flag to kubelet, kernel parameters changes are required
    - 'echo "====> Applying kernel parameters for Kubelet"'
    - 'sysctl -p /etc/sysctl.d/90-kubelet.conf'
  postKubeadmCommands:
    # Apply the privileged PodSecurityPolicy on the first master node ; Otherwise, CNI (and other) pods won't come up
    # Sometimes api server takes a little longer to respond. Retry if applying the pod-security-policy manifest fails
    - 'export KUBECONFIG=/etc/kubernetes/admin.conf && [ -f "$KUBECONFIG" ] && { echo " ====> Applying PodSecurityPolicy" ; until $(kubectl apply -f /etc/kubernetes/hardening/privileged-psp.yaml > /dev/null ); do echo "Failed to apply PodSecurityPolicies, will retry in 5s" ; sleep 5 ; done ; } || echo "Skipping PodSecurityPolicy for worker nodes"'

# Client configuration to add OIDC based authentication flags in kubeconfig
#clientConfig:
  #oidc-issuer-url: "{{ .spectro.pack.kubernetes.kubeadmconfig.apiServer.extraArgs.oidc-issuer-url }}"
  #oidc-client-id: "{{ .spectro.pack.kubernetes.kubeadmconfig.apiServer.extraArgs.oidc-client-id }}"
  #oidc-client-secret: 1gsranjjmdgahm10j8r6m47ejokm9kafvcbhi3d48jlc3rfpprhv
  #oidc-extra-scope: profile,email
    EOT 
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
    values = <<-EOT
manifests:
  storageclass:
    contents: |
      apiVersion: ceph.rook.io/v1
      kind: CephFilesystem
      metadata:
        name: myfs
        namespace: rook-ceph # namespace:cluster
      spec:
        metadataPool:
          replicated:
            size: 1
            requireSafeReplicaSize: false
        dataPools:
          - name: replicated
            failureDomain: osd
            replicated:
              size: 1
              requireSafeReplicaSize: false
        preserveFilesystemOnDelete: false
        metadataServer:
          activeCount: 1
          activeStandby: true
      ---
      apiVersion: storage.k8s.io/v1
      kind: StorageClass
      metadata:
        name: spectro-storage-class
        annotations:
          storageclass.kubernetes.io/is-default-class: "true"
      # Change "rook-ceph" provisioner prefix to match the operator namespace if needed
      provisioner: rook-ceph.cephfs.csi.ceph.com # driver:namespace:operator
      parameters:
        # clusterID is the namespace where the rook cluster is running
        # If you change this namespace, also change the namespace below where the secret namespaces are defined
        clusterID: rook-ceph # namespace:cluster

        # CephFS filesystem name into which the volume shall be created
        fsName: myfs

        # Ceph pool into which the volume shall be created
        # Required for provisionVolume: "true"
        pool: myfs-data0

        # The secrets contain Ceph admin credentials. These are generated automatically by the operator
        # in the same namespace as the cluster.
        csi.storage.k8s.io/provisioner-secret-name: rook-csi-cephfs-provisioner
        csi.storage.k8s.io/provisioner-secret-namespace: rook-ceph # namespace:cluster
        csi.storage.k8s.io/controller-expand-secret-name: rook-csi-cephfs-provisioner
        csi.storage.k8s.io/controller-expand-secret-namespace: rook-ceph # namespace:cluster
        csi.storage.k8s.io/node-stage-secret-name: rook-csi-cephfs-node
        csi.storage.k8s.io/node-stage-secret-namespace: rook-ceph # namespace:cluster

        # (optional) The driver can use either ceph-fuse (fuse) or ceph kernel client (kernel)
        # If omitted, default volume mounter will be used - this is determined by probing for ceph-fuse
        # or by setting the default mounter explicitly via --volumemounter command-line argument.
        # mounter: kernel
      reclaimPolicy: Delete
      allowVolumeExpansion: true
      #Supported binding modes are Immediate, WaitForFirstConsumer
      volumeBindingMode: "WaitForFirstConsumer"
      mountOptions:
        # uncomment the following line for debugging
        #- debug
  cluster:
    contents: |
      kind: ConfigMap
      apiVersion: v1
      metadata:
        name: rook-config-override
        namespace: rook-ceph # namespace:cluster
      data:
        config: |
          [global]
          osd_pool_default_size = 1
          mon_warn_on_pool_no_redundancy = false
          bdev_flock_retry = 20
          bluefs_buffered_io = false
      ---
      apiVersion: ceph.rook.io/v1
      kind: CephCluster
      metadata:
        name: my-cluster
        namespace: rook-ceph # namespace:cluster
      spec:
        dataDirHostPath: /var/lib/rook
        cephVersion:
          image: quay.io/ceph/ceph:v16.2.7
          allowUnsupported: true
        mon:
          count: 1
          allowMultiplePerNode: true
        mgr:
          count: 1
          allowMultiplePerNode: true
        dashboard:
          enabled: true
        crashCollector:
          disable: true
        storage:
          useAllNodes: true
          useAllDevices: true
          deviceFilter: sdb
        healthCheck:
          daemonHealth:
            mon:
              interval: 45s
              timeout: 600s
        disruptionManagement:
          managePodBudgets: true   
    EOT
  }
  pack {
    name   = data.spectrocloud_pack.proxy.name
    tag    = data.spectrocloud_pack.proxy.version
    uid    = data.spectrocloud_pack.proxy.id
    values = data.spectrocloud_pack.proxy.values
  }
}
