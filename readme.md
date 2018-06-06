# Gist store
Use Github Gists as a store persisting option.

## Why?
I needed a simple way to persistent data between scheduled jobs. 

## Example
```go
import (
	"os"

	gist "github.com/X11/go-gist-store"
)

func main() {
	type TestState struct {
		Hello string `json:"hello"`
	}

	gs := gist.New(os.Getenv("GITHUB_GIST_ID"), os.Getenv("GITHUB_GIST_FILENAME"), os.Getenv("GITHUB_AUTHENTICATION"))

	s := TestState{}
	gs.Get(&s)
	s.Hello = "world"
	gs.Update(&s)
}

```

