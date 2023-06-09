package git

import (
	"os"
	"strings"
	"testing"

	"github.com/mholt/archiver"
	"github.com/softsense/git-semver/pkg/semver"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	if err := archiver.Unarchive("testdata/repo.tar.gz", "testdata/"); err != nil {
		panic(err)
	}

	exitCode := m.Run()

	os.RemoveAll("testdata/repo")

	os.Exit(exitCode)
}

func TestOpen(t *testing.T) {
	g, err := Open("testdata/repo", Config{
		Prefix: "v",
	})
	if err != nil {
		t.Fatal(err)
	}

	n, err := g.Increment(false, false, true, false)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("\n" + n.String())

	require.Equal(t, "v0.0.3", n.String())
}

func TestBelow(t *testing.T) {
	tests := []struct {
		name   string
		below  semver.Version
		expect semver.Version
	}{
		{
			name:   "higher",
			below:  semver.MustParse("v9.9.9"),
			expect: semver.MustParse("v0.0.2"),
		},
		{
			name:   "same",
			below:  semver.MustParse("v0.0.2"),
			expect: semver.MustParse("v0.0.1"),
		},
		{
			name:   "below",
			below:  semver.MustParse("v0.0.1"),
			expect: semver.MustParse("v0.0.0"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g, err := Open("testdata/repo", Config{
				Prefix: "v",
				Below:  &test.below,
			})
			require.NoError(t, err)
			require.Equal(t, test.expect, g.Highest())
		})
	}
}

func TestHistory(t *testing.T) {
	g, err := Open("./testdata/repo", Config{
		Prefix: "v",
	})
	if err != nil {
		t.Fatal(err)
	}

	history, err := g.History("|||")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("\n" + history)

	for _, l := range strings.Split(history, "\n") {
		if l != "" && !strings.HasPrefix(l, "|||") {
			t.Fatalf("Expected all lines to have prefix '|||', got line: '%s'", l)
		}
	}
}
