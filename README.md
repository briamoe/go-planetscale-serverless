<div align="center">
  <h1>🪐 PlanetScale Serverless Driver for Go</h1>
  <p>Makes doing queries to PlanetScale over HTTP in Go possible, allowing for the use of PlanetScale in serverless applications. Based on the <a href="https://github.com/planetscale/database-js">PlanetScale Serverless Driver for JavaScript.</a></p>
</div>


## ⚠️ Experimental
This library is incredibly experimental, as it uses undocumented API's only exposed by the aforementioned JavaScript driver. Things are subject to change as they update the driver.

## Quick Start
Here's the essentials you need to get up and running with the library.

```go
package main

import (
	"https://github.com/briamoe/pscale-serverless-go"
)

type User struct {
	ID   int    `ps:"id"`
	Name string `ps:"name"`
}

func main() {
	c, err := planetscale.NewConnection(&planetscale.Config{
		Username: "<username>",
		Password: "<password>",
		Host:     "<host>",
	})
	if err != nil {
		panic(err)
	}

	e, err := c.Execute("SELECT * FROM users WHERE id=?", 1);
	if err != nil {
		panic(err)
	}

	var u []*User
	if err = e.Decode(&u); err != nil {
		panic(err)
	}
}

```

## Installation
Requires Go version 1.18 or higher.
```
go get -u https://github.com/briamoe/pscale-serverless-go
```

## Connection
Creating a new connection creates a shared session that can be used across queries to PlanetScale. 

To generate the credentials highlighted in `<>` below, head to your [planetscale dashboard](https://app.planetscale.com/) and select your database. Then go to `Settings` > `Passwords` > `New password`, and paste in the values it gives you.
```go
planetscale.NewConnection(&planetscale.Config{
	Username: "<username>",
	Password: "<password>",
	Host:     "<host>",
})
```

## Executing
Executing performs a query to a PlanetScale Database.
```go
connection.Execute("SELECT * FROM users WHERE id=?", 1)
```

### Selects
If a statement includes a select, you can grab the rows returned by using `Decode(out interface{})`
```go
type User struct {
	ID   int    `ps:"id"`
	Name string `ps:"name"`
}

var []*User
executed.Decode(&u)
```

The use of `ps` tags are an imporant part of the decoding process, as it determines which fields inside sql are mapped to fields in a struct. Make sure to incude these!

### Inserts or Updates
If a statement includes a insert or an update statement, you can access how many rows were affected or the insert id. These are set to 0 if a query doesn't include any of those.
```go
fmt.Println(executed.RowsAffected) // 1
fmt.Println(executed.InsertID) // 200
```