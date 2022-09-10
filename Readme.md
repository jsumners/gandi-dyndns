# gdyndns

This is a tool for updating DNS records in [Gandi's](https://gandi.net)
[LiveDNS](https://api.gandi.net/docs/livedns/) system. It is meant to be used
in a recurring job to provide "dynamic dns" like functionalty.

For each zone in the configuration, the current DNS record will be queried.
If the values returned for the record do not contain the current public IP
address for the system, as determined by querying `http://ip-api.com/json/`,
then the record will be updated with the new IP address.

## Configuration

The tool requires a configuration file. The file, named `config.yaml`, can
be located in:

+ current working directory
+ `/etc/gdyndns/`
+ `$HOME/.gdyndns/`

Alternatively, it can be spefied by the envionrment variable
`GDYNDNS_CONFIG_FILE`, e.g `export GDYNDNS_CONFIG_FILE=/opt/gdyndns.yaml`.

The format of the configuration is as follows:

```yaml
# Must be set to your Gandi v5 LiveDNS key. See
# https://api.gandi.net/docs/authentication/
gandi_v5_api_key: super-secret-key

# Defines the set of records that should be checked and updated.
# Each record is an object with the properties:
#
# + `zone` (string): the root domain of the record.
# + `type` (string): the type of DNS record, e.g. `A`.
# + `name` (string): the subdomain to update, e.g. "foo" in "foo.example.com".
#    Set to `'@'` to update the root record.
# + `ttl` (integer, optional): Time-to-live for the record. Minimum is 300.
#    Default: 300.
records:
  - zone: example.com
    type: A
    name: '@'
    ttl: 1500

  - zone: example.com
    type: A
    name: foo
```
