# Table: rspace_document

An RSpace Document is the most fundamental element of knowledge in RSpace ELN.

## Examples

### Basic info

```sql
select
  name,
  global_id,
  tags,
  owner_username 
from
  rspace_document
```

### List documents created by a user

```sql
select
  name,
  global_id,
  tags,
  owner_username 
from
  rspace_document 
where
  owner_username = 'dcopper'
  ```

### List documents with default name 'Untitled document'

```sql
select
  name,
  global_id,
  tags,
  owner_username 
from
  rspace_document 
where
  name = 'Untitled document'
```

### List documents modified in last 28 days

```sql
select
  name,
  global_id,
  tags,
  owner_username 
from
  rspace_document 
where
  last_modified = now() - interval '28d'
```


### Get a single document by its global ID

```sql
select
  name,
  global_id,
  tags,
  owner_username 
from
  rspace_document 
where
  global_id = 'SD12345'
```