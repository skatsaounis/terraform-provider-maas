## 2.0.0 (Dec 13, 2023)

NEW:

* Add support for resource and data source maas_device (#119)
* Deploy a machine based on its system_id (#118)
* Add data sources for machine and physical nic (#7)
* Add support for TLS configuration options (#101)

IMPROVEMENTS:

* Consume versioned `maas/gomaasclient` starting with 0.1.0 (#124)
* Fixes on expecting the VLAN id and parsing big numbers (#121)
* chore: replace deprecated package (#114)
* Fix maas_instance updates (#117)
* ci: special labels to trigger integration tests (#109)
* Fix parsing of power parameters (#96)
  * It includes a breaking change to how power parameters have to be declared. Please consult the example: `terraform-provider-maas/examples/1-machines.tf`.
* Update dependencies

## 1.3.1 (Oct 12, 2023)

IMPROVEMENTS:

* Update `gomaasclient` to include:
  * bugfix related to retry logic

## 1.3.0 (Sep 28, 2023)

NEW:

* Add `comment`, `definition`, `kernel_opts` fields to `tag` resource

IMPROVEMENTS:

* Update `gomaasclient` to include:
  * changes related to proper parsing of machine fields
  * changes related to retry improvements
* Update dependencies
* Update Go version to 1.20

## 1.2.0 (May 12, 2023)

NEW:

* Add `enable_hw_sync` to `deploy_params` of `maas_instance` resource

## 1.1.0 (Mar 10, 2023)

NEW:

* Add release process guide
* Add GitHub Actions workflow for releasing
* Add Dependabot configuration

IMPROVEMENTS:

* Update documentation and documentation structure to use terraform-plugin-docs tool
* Refactor API client references to use the maas repo
* Modify Resource `maas_machine` with timeouts support
* Modify Resource `maas_instance` with timeouts support
* Update dependencies
* Update Go version to 1.18

## 1.0.0 (Jul 21, 2021)

NEW:

* Resource `maas_block_device`
* Resource `maas_fabric`
* Resource `maas_vlan`
* Resource `maas_subnet`
* Resource `maas_space`
* Resource `maas_subnet_ip_range`
* Resource `maas_dns_domain`
* Resource `maas_dns_record`
* Resource `maas_user`
* Resource `Implement importers for the existing managed resources:`
* Resource `maas_machine`
* Resource `maas_instance`
* Resource `maas_tag`
* Resource `maas_network_interface_physical`
* Resource `maas_vm_host`
* Resource `maas_vm_host_machine`

IMPROVEMENTS:

* Allow env vars `MAAS_API_KEY` and `MAAS_API_URL` to be used for the provider configuration.
* Use VM host naming instead of Pod.
* Add validation for maas_machine resource power_type argument.
* Update VM host machine network and storage params.
* Properly implement the network and storage parameters for the `maas_vm_host_machine` resource.
* Remove managed argument from `maas_subnet` resource and data source. This is considered `true` by default on MAAS 2.0 and newer.
* Update docs and examples.
