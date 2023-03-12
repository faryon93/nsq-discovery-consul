# nsq-discovery-consul
This application is intended as a sidecar for nsqd and nsqadmin. It queries consul for a specific service and uses the [HTTP /config API](https://nsq.io/components/nsqd.html#put-confignsqlookupd_tcp_addresses) to set nsqlookupd endpoints.

The consul client library uses a [set of environment variables](https://github.com/hashicorp/consul/blob/api/v1.20.0/api/api.go#L24) to configure the consul connection.

## Usage
```shell
$: docker run faryon93/nsq-discovery-consul \
    /nsq-discovery-consul \
    --consul-lookupd-service="nsqlookupd" \
    --nsq-conf-url="http://127.0.0.1:4151/config/nsqlookupd_tcp_addresses"
```
