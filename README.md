Queryf
======

This is a simple project to enable printing Go SQL queries into strings for debugging reasons.

**Attention**: You should not use it to make requests to the database, to do that you should bind the parameters.

Usage
-----
To use you need just to assign the Print the variable to a string then print it to console:

```golang
printedQuery := queryf.Print(query, args...)
fmt.Println(printedQuery)
```