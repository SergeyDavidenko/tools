## Golang tools


### Install logger
```
go get -u github.com/SergeyDavidenko/tools/pkg/logger 
```

### Install config
```
go get -u github.com/SergeyDavidenko/tools/pkg/config 
```

### How to use config
``` sh
mkdir config
touch example.yml
```

``` yml
http:
  api:
    hostString: ":8080"
    readTimeout: 15s
    writeTimeout: 15s
    maxHeaderMegabytes: 4
  healtz:
    hostString: ":1499"
    readTimeout: 15s
    writeTimeout: 15s
    maxHeaderMegabytes: 4
redis:
  hostString: "localhost:6379"
  login: ""
  password: ""
  connectionsCount: 5
  dbNum: 1
postgres:
  hostString: "localhost"
  port: 5432
  login: "go"
  password: "go"
  dbName: "example"
  connectionsCount: 20
mongodb:
  hostString: "localhost"
  port: 27017
  login: "go"
  password: "go"
  dbName: "example"
custom:
  version: "0.0.1"
```

``` go
package main

import (
	"log"

	"github.com/SergeyDavidenko/tools/pkg/config"
)

func main() {
	cfg, err := config.New("config/", "example")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg.GetHTTP("api").HostString)
	log.Println(cfg.BuildDSNPostgres())
	log.Println(cfg.GetRedis().HostString)
	log.Println(cfg.GetCustom("version"))

}
```