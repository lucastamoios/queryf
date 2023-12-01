Queryf
======


This is a simple project to enable printing Go SQL queries into strings for debugging reasons.

**Attention**: This is not meant to be used in production. To avoid SQL injections you should
bind the parameters usign the `database/sql` package. Read more about it
[here](https://use-the-index-luke.com/sql/where-clause/bind-parameters).

Usage
-----

To use this package just get the resulting string from the `Format` function. After that you
can print it or do whatever you want with it.

Example:

```golang
query := "SELECT * FROM users WHERE id = $1 AND name = $2"
args := []any{1, "John"}
fmt.Println(queryf.Format(query, args...))
// Output: SELECT * FROM users WHERE id = 1 AND name = 'John'
```