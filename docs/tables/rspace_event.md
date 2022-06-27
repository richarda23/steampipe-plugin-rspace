# Table: rspace_document

An RSpace Event records an action performed in RSpace.

The table can be queried by domain and / or action.

Only events performed in the last 6 months are queryable via the RSpace API, and a maximum of 1000 events
are returned from RSpace API.

Valid domains are:

- GROUP
- INV_SAMPLE
- INV_SUBSAMPLE
- INV_CONTAINER
- NOTEBOOK
- RECORD
- USER

Valid actions are:

- EXPORT
- CREATE
- READ
- SHARE
- SIGN
- WRITE

## Examples

### Basic info

```sql
select
  name,
  username,
  full_name,
  timestamp,
  payload
from
  rspace_event
```

### List 'create document' actions

```sql
select
  username,
  timestamp,
  payload
from
  rspace_event 
where
  domain = 'RECORD'
and
  action = 'CREATE'
limit 10
  ```

