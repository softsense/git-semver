package git

import (
	"os"
	"strings"
	"testing"

	"github.com/mholt/archiver"
	"github.com/softsense/git-semver/pkg/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/storage/memory"
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

	n, err := g.Increment(false, false, true, false, false)
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

func TestInsertPullRequestURL(t *testing.T) {
	tests := []struct {
		name      string
		remoteUrl string
		msg       string
		want      string
	}{
		{
			name:      "Regular Github pull request ref - SSH remote",
			remoteUrl: "git@github.com:foo/bar.git",
			msg:       "15f53d3 Some commit message (#1)",
			want:      "15f53d3 Some commit message [(#1)](https://github.com/foo/bar/pull/1)",
		},
		{
			name:      "Multi-line commit with PR ref - SSH remote",
			remoteUrl: "git@github.com:foo/bar.git",
			msg:       "15f53d3 Some commit message (#1)\n\nSome description",
			want:      "15f53d3 Some commit message [(#1)](https://github.com/foo/bar/pull/1)\n\nSome description",
		},
		{
			name:      "Regular Github pull request ref - HTTPS remote",
			remoteUrl: "https://github.com/foo/bar",
			msg:       "15f53d3 Some commit message (#1)",
			want:      "15f53d3 Some commit message [(#1)](https://github.com/foo/bar/pull/1)",
		},
		{
			name:      "Multi-line commit with PR ref - HTTPS remote",
			remoteUrl: "https://github.com/foo/bar",
			msg:       "15f53d3 Some commit message (#1)\n\nSome description",
			want:      "15f53d3 Some commit message [(#1)](https://github.com/foo/bar/pull/1)\n\nSome description",
		},
		{
			name:      "Non-GitHub remote - SSH",
			remoteUrl: "git@example.com:foo/bar.git",
			msg:       "15f53d3 Some commit message (#1)",
			want:      "15f53d3 Some commit message (#1)",
		},
		{
			name:      "Non-GitHub remote - HTTPS",
			remoteUrl: "https://example.com/foo/bar",
			msg:       "15f53d3 Some commit message (#1)",
			want:      "15f53d3 Some commit message (#1)",
		},
	}

	for _, test := range tests {
		tc := test // don't close over loop variable
		t.Run(tc.name, func(t *testing.T) {
			g := &Git{
				repo: &git.Repository{
					Storer: memory.NewStorage(),
				},
			}
			_, err := g.repo.CreateRemote(&config.RemoteConfig{
				Name:  "origin",
				URLs:  []string{tc.remoteUrl},
				Fetch: nil,
			})
			require.NoError(t, err)

			msg := insertPullRequestURL(tc.msg, g)
			assert.Equal(t, tc.want, msg)
		})
	}
}
