package git

import (
	"fmt"
	"testing"
)

func TestOpen(t *testing.T) {
	g, err := Open("../../", Config{
		Prefix: "v",
	})
	if err != nil {
		t.Fatal(err)
	}

	n, err := g.Increment(false, false, true, false)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(n.String())
}

func TestHistory(t *testing.T) {
	g, err := Open("../../", Config{
		Prefix: "v",
	})
	if err != nil {
		t.Fatal(err)
	}

	history, err := g.History()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Print(history)
}
