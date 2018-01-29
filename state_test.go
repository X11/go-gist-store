package gstate_test

import (
	"os"
	"testing"

	gstate "github.com/X11/go-gstate"
)

func TestShouldFetchTheState(t *testing.T) {
	type TestState struct {
		Hello string `json:"hello"`
	}

	gs := gstate.New(os.Getenv("GITHUB_GIST_ID"), os.Getenv("GITHUB_GIST_FILENAME"), os.Getenv("GITHUB_AUTHENTICATION"))

	s := TestState{}
	gs.Get(&s)

	if s.Hello != "world" && s.Hello != "world2" {
		t.Fail()
	}
}

func TestShouldFetchAndUpdateTheState(t *testing.T) {
	type TestState struct {
		Hello string `json:"hello"`
	}

	// Reset
	gs := gstate.New(os.Getenv("GITHUB_GIST_ID"), os.Getenv("GITHUB_GIST_FILENAME"), os.Getenv("GITHUB_AUTHENTICATION"))
	s := TestState{}
	gs.Get(&s)
	s.Hello = "world"
	gs.Update(&s)
	gs.SetFetched(false)

	// First check
	s2 := TestState{}
	gs.Get(&s2)

	if s2.Hello != "world" {
		t.Fail()
	}

	// Update again
	s2.Hello = "world2"

	gs.Update(&s2)
	gs.SetFetched(false)

	// Second check
	s3 := TestState{}
	gs.Get(&s3)

	if s3.Hello != "world2" {
		t.Fail()
	}
}
