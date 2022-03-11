terraform {
  required_version = ">= 1.0.0"

  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

provider "digitalocean" {}

resource "digitalocean_app" "proxy" {
  spec {
    name   = "proxy"
    region = "nyc1"

    service {
      name               = "proxy-service"
      instance_size_slug = "basic-xxs"

      git {
        branch         = "digitalocean"
        repo_clone_url = "https://github.com/br7552/proxy"
      }
    }
  }
}

output "live_url" {
  value       = digitalocean_app.proxy.live_url
  description = "The live URL of the proxy server."
}
