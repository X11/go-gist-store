# Gist States
Use Github Gists as state manangement between jobs.

## Why?
I needed a simple way to persistent data between scheduled jobs. 

## Example
```go
import (
	"os"

	gstate "github.com/X11/go-gstate"
)

func main() {
	type TestState struct {
		Hello string `json:"hello"`
	}

	gs := gstate.New(os.Getenv("GITHUB_GIST_ID"), os.Getenv("GITHUB_GIST_FILENAME"), os.Getenv("GITHUB_AUTHENTICATION"))

	s := TestState{}
	gs.Get(&s)
	s.Hello = "world"
	gs.Update(&s)
}

```

