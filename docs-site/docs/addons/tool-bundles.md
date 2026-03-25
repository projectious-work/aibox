---
sidebar_position: 3
title: "Tool Bundles"
---

# Tool Bundle Addons

Tool bundles install infrastructure, orchestration, and cloud CLI tools.

## Infrastructure

```toml
[addons.infrastructure.tools]
opentofu = {}      # Infrastructure-as-code (Terraform alternative)
ansible = {}       # Configuration management
packer = {}        # Machine image builder
```

OpenTofu and Packer are compiled in a multi-stage builder. Ansible is installed via pip.

## Kubernetes

```toml
[addons.kubernetes.tools]
kubectl = {}       # Kubernetes CLI
helm = {}          # Package manager
kustomize = {}     # Configuration customization
# k9s = {}         # Optional: terminal UI for Kubernetes
```

All tools are downloaded as static binaries in a multi-stage builder.

## Cloud Providers

### AWS

```toml
[addons.cloud-aws.tools]
aws-cli = {}
```

Installs the AWS CLI v2.

### Google Cloud

```toml
[addons.cloud-gcp.tools]
gcloud-cli = {}
```

Installs the Google Cloud CLI via the official APT repository.

### Azure

```toml
[addons.cloud-azure.tools]
azure-cli = {}
```

Installs the Azure CLI via pip.
