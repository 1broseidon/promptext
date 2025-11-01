Language: Go
Version: 1.22.0
Project Structure:
└── config.yaml
└── go.mod
└── main.go


## Source Files

### config.yaml (4 lines)
```yaml
server:
  port: 8080
  host: localhost

```

### go.mod (4 lines)
```mod
module example.com/go-service

go 1.22.0

```

### main.go (14 lines)
```go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Go!")
	})
	http.ListenAndServe(":8080", nil)
}

```
