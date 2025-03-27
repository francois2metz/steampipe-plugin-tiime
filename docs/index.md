---
organization: francois2metz
category: ["SaaS"]
brand_color: "#7680ff"
display_name: "Tiime"
short_name: "tiime"
description: "Steampipe plugin for querying invoices from Tiime."
og_description: "Query Tiime with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/francois2metz/tiime-social-graphic.png"
icon_url: "/images/plugins/francois2metz/tiime.svg"
---

# Tiime + Steampipe

[Tiime](https://www.tiime.fr/) is an invoicing software.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

```sql
select
  id,
  number,
  emission_date
from
  tiime_invoice
```

```
+----------+--------+----------------------+
| id       | number | emission_date        |
+----------+--------+----------------------+
| 13863849 | 129    | 2025-03-04T00:00:00Z |
| 13603916 | 124    | 2025-02-22T00:00:00Z |
| 14074064 | 131    | 2025-03-14T00:00:00Z |
| 14074126 | 132    | 2025-03-14T00:00:00Z |
| 13684506 | 127    | 2025-02-28T00:00:00Z |
| 13788235 | 126    | 2025-02-28T00:00:00Z |
| 13982585 | 130    | 2025-03-10T00:00:00Z |
| 12650352 | 125    | 2025-02-22T00:00:00Z |
| 13793693 | 128    | 2025-03-01T00:00:00Z |
+----------+--------+----------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/francois2metz/tiime/tables)**

## Get started

### Install

Download and install the latest Tiime plugin:

```bash
steampipe plugin install ghcr.io/francois2metz/tiime
```

### Credentials

| Item        | Description                                                                                                                                                                   |
|-------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| Credentials | Tiime requires a login and a password                                                                                                                                         |
| Permissions | Tokens have the same permissions as the user who creates them.                                                                                                                |
| Radius      | Each connection represents a single Scalingo account.                                                                                                                         |
| Resolution  | 1. Credentials explicitly set in a steampipe config file (`~/.steampipe/config/tiime.spc`)<br />2. Credentials specified in environment variables, e.g., `TIIME_EMAIL`. |

### Configuration

Installing the latest tiime plugin will create a config file (`~/.steampipe/config/tiime.spc`) with a single connection named `tiime`:

```hcl
connection "tiime" {
    plugin = "ghcr.io/francois2metz/tiime"

    # The Tiime email
    # This can also be set via the `TIIME_EMAIL` environment variable.
    # email = "test@example.net"

    # The Tiime password
    # This can also be set via the `TIIME_PASSWORD` variable.
    # password = "EefDeabJeshejror"

    # The company id to use
    # This can also be set via the `TIIME_COMPANY_ID` variable.
    # company_id = 1234
}
```

### Credentials from Environment Variables

The Tiime plugin will use the following environment variables to obtain credentials **only if other argument (`email`, `password`) is not specified** in the connection:

```sh
export TIIME_EMAIL=text@example.net
export TIIME_PASSWORD=EefDeabJeshejror
export TIIME_COMPANY_ID=1234
```

## Get Involved

* Open source: https://github.com/francois2metz/steampipe-plugin-tiime
