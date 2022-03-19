terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = "~> 2.4"
    }
  }
}

provider "digitalocean" {
  token = var.do_token
}

provider "helm" {
  kubernetes {
    host  = digitalocean_kubernetes_cluster.do1.endpoint
    token = digitalocean_kubernetes_cluster.do1.kube_config[0].token
    cluster_ca_certificate = base64decode(
      digitalocean_kubernetes_cluster.do1.kube_config[0].cluster_ca_certificate
    )
  }
}

resource "digitalocean_kubernetes_cluster" "do1" {
  name    = "do-cluster-1"
  region  = var.region
  version = var.k8s-version

  node_pool {
    name       = "autoscale-worker-pool"
    size       = "s-1vcpu-2gb"
    auto_scale = true
    min_nodes  = 1
    max_nodes  = 5
  }
}

resource "helm_release" "prometheus" {
  depends_on = [helm_release.cert-manager]

  name             = "prometheus"
  repository       = "https://prometheus-community.github.io/helm-charts"
  chart            = "kube-prometheus-stack"
  namespace        = "prometheus"
  create_namespace = true
  set {
    name  = "prometheusOperator.admissionWebhooks.certManager.enabled"
    value = "true"
  }
}

resource "helm_release" "cert-manager" {
  depends_on = [digitalocean_kubernetes_cluster.do1]

  name             = "cert-manager"
  repository       = "https://charts.jetstack.io"
  chart            = "cert-manager"
  namespace        = "cert-manager"
  create_namespace = true
  set {
    name  = "installCRDs"
    value = "true"
  }
}
