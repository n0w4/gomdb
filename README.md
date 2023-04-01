# gomdb

is a very simple in memory database

for my proposal I shouldn't delete documents. only:

- insert 
- update 
- find 

## Install

```shell
go get github.com/n0w4/gompb
```


## Example usage

```go
db := gomdb.NewMemoryDB("exampleDB")

doc1 := map[string]interface{}{
    "name": "John",
    "age":  30,
}

doc2 := map[string]interface{}{
    "name": "Jane",
    "age":  28,
}

doc3 := map[string]interface{}{
    "name": "Jaoline",
    "age":  35,
}

db.InsertOnCollection("users", doc1)
db.InsertOnCollection("users", doc2)
db.InsertOnCollection("users", doc3)

filter := map[string]interface{}{
    "name": `ne$`,
    "age": 28,
}

// FIND SOMEONE
filteredUsers := db.FindOnCollection("users", filter)


// UPDATE SOMEONE
update := map[string]interface{}{
    "age": 36,
}

totalUpdated := db.UpdateOnCollection("users", filter, update)

```