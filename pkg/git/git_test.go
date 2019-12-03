package git

import (
	"fmt"
	"testing"
)

func TestOpen(t *testing.T) {
	g, err := Open("../../")
	if err != nil {
		t.Fatal(err)
	}

	n, err := g.Increment(false, false, true, false)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(n.String())
}
