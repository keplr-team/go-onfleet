# go-onfleet
Go client for accessing [Onfleet API](https://docs.onfleet.com/reference)

## Usage
```bash
go get github.com/keplr-team/go-onfleet
```

```go
import (
    "context"
    "fmt"
    "github.com/keplr-team/go-onfleet/onfleet"
)

client := onfleet.NewClient("API_KEY")
workers, err := client.Workers.List(context.Background(), nil)
if err != nil {
	// do stuff
}
fmt.Printf("%+v\n", workers)
```
### What's next

- [ ] unit testing
- [ ] handle errors
- [ ] handle rate limit
- [ ] add more endpoints
