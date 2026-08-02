## Cenari of tests

First test
```go
query := `
SELECT id, type, bu, service, environment, action, realms, accounts, message, created_by, created_at, update_by, update_by
FROM rules
WHERE bu = $1
ORDER BY created_at DESC
`
```
