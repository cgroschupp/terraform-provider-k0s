# Terraform provider for k0s (using embedded k0sctl)

> ⚠️ The provider name on Terraform registry has been renamed from `alessiodionisi/k0s` to `cgroschupp/k0s`, please update your `required_providers` block.

Terraform provider to create and manage [k0s](https://k0sproject.io) Kubernetes clusters, using embedded [k0sctl](https://github.com/k0sproject/k0sctl).

## Getting started

You can install stable releases of the provider from the [Terrafom registry](https://registry.terraform.io/providers/cgroschupp/k0s/latest).

## Uploading inline files

Each entry in `k0s_cluster.hosts[].files` accepts either `src` (a local file path)
or `content` (inline text), but not both. For example, inside a host configuration:

```hcl
files = [
  {
    name    = "app-config"
    dst     = "/etc/app.conf"
    perm    = "0600"
    content = <<-EOT
      enabled=true
      endpoint=https://example.com
    EOT
  }
]
```

Inline content must not be empty. It is uploaded through k0sctl's file upload
phase before k0s configuration and initialization. Although `content` is marked
sensitive, its value is stored in Terraform state; protect state access accordingly.

## Contributing

If you wish to work on the provider, you'll first need these requirements:

- [Terraform](https://www.terraform.io) >= 1.0
- [Go](https://golang.org) >= 1.21
- [Task](https://taskfile.dev) >= 1.3

To compile the provider, run `task install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `task generate`.

In order to run the full suite of acceptance tests, run `task acceptance-tests`.

_Note:_ acceptance tests create real resources, and often cost money to run.
