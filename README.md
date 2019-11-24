A service discovery for the [ECS](https://www.aliyun.com/) cloud platform compatible with [Prometheus](https://prometheus.io).

## How it works

This service gets the list of servers from the Aliyun API and generates a file which is compatible with the Prometheus `file_sd` mechanism.

## Pre-requisites

You need your Aliyun access key/secret key (token). You can create this token [in the console](https://usercenter.console.aliyun.com/).

## Installing it

Download the binary from the [Releases](https://github.com/seanly/prometheus-ecs-sd/releases) page.

## Running it

```
usage: sd adapter usage --config.file=ecs_sd_config.yml [<flags>]

Tool to generate Prometheus file_sd target files for Aliyun Ecs.

Flags:
  -h, --help                    Show context-sensitive help (also try --help-long and --help-man).
      --output.file="ecs.json"  The output filename for file_sd compatible file.
      --config.file="ecs_sd_config.yml"       The ecs sd config file.
      --web.listen-address=":9465"
                                The listen address.
      --version                 Show application version.
```

## Integration with Prometheus

Here is a Prometheus `scrape_config` snippet that configures Prometheus to scrape node_exporter assuming that it is deployed on all your Ecs servers.

```yaml
- job_name: node

  # Assuming that prometheus and prometheus-ecs-sd are started from the same directory.
  file_sd_configs:
  - files: [ "./ecs.json" ]

  # The relabeling does the following:
  # - overwrite the scrape address with the node_exporter's port.
  # - strip leading commas from the tags label.
  # - save the region label (par1/ams1).
  # - overwrite the instance label with the server's name.
  relabel_configs:
  - source_labels: [__meta_ecs_private_ip]
    replacement: "${1}:9100"
    target_label: __address__
  - source_labels: [__meta_ecs_tag_groupId]
    regex: ",(.+),"
    target_label: tags
  - source_labels: [__meta_ecs_region_id]
    target_label: region
  - source_labels: [__meta_ecs_instance_name]
    target_label: instance
```

The following meta labels are available on targets during relabeling:

* `__meta_instance_name`
* `__meta_instance_id`
* `__meta_instance_type`
* `__meta_private_ip`
* `__meta_public_ip`
* `__meta_ecs_tag_xxx`
* `__meta_ecs_region_id`: the identifier of the zone (region).


## Contributing

PRs and issues are welcome.

## License

Apache License 2.0, see [LICENSE](https://github.com/seanly/prometheus-ecs-sd/blob/master/LICENSE).
