_# Terraform Provider for Grafana Adaptive Metrics

- Grafana website: https://grafana.com
- Grafana Cloud website: https://grafana.com/products/cloud/
- Grafana Adaptive Metrics website: https://grafana.com/docs/grafana-cloud/cost-management-and-billing/reduce-costs/metrics-costs/control-metrics-usage-via-grafana-adaptive-metrics/

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.20

## Development

This repository is built on the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework).


### Building the provider

Build the provider using the Go `install` command:

```shell
go install
```

This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

### Using the provider

Add the following to your `.terraformrc` to test with a local version of the provider:

```
provider_installation {
  dev_overrides {
      "registry.terraform.io/grafana/grafana-adaptive-metrics" = "/$GOPATH/bin"
  }

  # For all other providers, install them directly from their origin provider
  # registries as normal. If you omit this, Terraform will _only_ use
  # the dev_overrides block, and so no other providers will be available.
  direct {}
}
```

### Debugging the provider

1. Build the provider:
    ```
    go build -gcflags "all=-N -l" -o terraform-provider-grafana-adaptive-metrics .
   ```
2. Run w/ delve:
    ```
    dlv exec --accept-multiclient --listen=:2345 --continue --headless ./terraform-provider-grafana-adaptive-metrics -- -debug`
    ```
3. Connect your IDE debugger to the delve instance (listening on port 2345).
4. The `dlv` command will output something that starts with `TF_REATTACH_PROVIDERS`; prepend that to the terraform command you're testing. For example:
    ```
   TF_REATTACH_PROVIDERS='{"registry.terraform.io/my-org/my-provider":{"Protocol":"grpc","Pid":3382870,"Test":true,"Addr":{"Network":"unix","String":"/tmp/plugin713096927"}}}' terraform plan
    ```

### Running acceptance tests

In order to run the full suite of Acceptance tests, run `make testacc`.

```shell
make testacc
```

Acceptance tests expect the `GRAFANA_CLOUD_API_URL` and `GRAFANA_CLOUD_API_KEY` environment variables to be set.

### Updating documentation

To generate or update documentation, run `go generate`.
