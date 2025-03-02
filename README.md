# Tiime plugin for Steampipe

Use SQL to query invoices from [Tiime][].

- **[Get started →](docs/index.md)**
- Documentation: [Table definitions & examples](docs/tables)

## Quick start

Install the plugin with [Steampipe][]:

    steampipe plugin install francois2metz/tiime

## Development

To build the plugin and install it in your `.steampipe` directory

    make

Copy the default config file:

    cp config/tiime.spc ~/.steampipe/config/tiime.spc

## License

Apache 2

[steampipe]: https://steampipe.io
[tiime]: https://www.tiime.fr/
