# Table: tiime_invoice

The `tiime_invoice` table can be used to query information about the invoices within Tiime.

## Examples

### List invoices

```sql
select
  id,
  emission_date,
  total_excluding_taxes
from
  tiime_invoice;
```

### List unpaid invoice

```sql
select
  id,
  emission_date,
  total_excluding_taxes,
  client_name
from
  tiime_invoice
where
  status != 'paid';
```

### List draft invoice

```sql
select
  id,
  emission_date,
  total_excluding_taxes
from
  tiime_invoice
where
  status = 'draft';
```
