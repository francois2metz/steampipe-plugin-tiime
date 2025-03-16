# Table: tiime_client

The `tiime_client` table can be used to query information about the client within Tiime.

## Examples

### List clients

```sql
select
  id,
  name
from
  tiime_client;
```

### List clients with unpaid invoices

```sql
select
  id,
  name,
  balance_including_taxes
from
  tiime_client
where
  payment_status = 'late'
```
