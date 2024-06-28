package git

import (
	"fmt"
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
	// repo.tar.gz has:
	// v0.0.1
	// v0.0.2
	// v0.0.3-rc1
	//
	// repo-rc.tar.gz has:
	// v0.0.1
	// v0.0.2
	// v0.3.0

	for _, name := range []string{"repo.tar.gz", "repo-rc.tar.gz"} {
		if err := archiver.Unarchive(fmt.Sprintf("testdata/%s", name), "testdata/"); err != nil {
			panic(err)
		}
	}

	exitCode := m.Run()

	os.RemoveAll("testdata/repo")
	os.RemoveAll("testdata/repo-rc")

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

func TestOpenNonExistingPath(t *testing.T) {
	_, err := Open("testdata/does-not-exist", Config{
		Prefix: "v",
	})

	require.Error(t, err)
	require.EqualError(t, err, "open git repo testdata/does-not-exist: repository does not exist")
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

func TestRC(t *testing.T) {
	tests := []struct {
		name      string
		expect    semver.Version
		includeRC bool
	}{
		{
			name:   "no rc",
			expect: semver.MustParse("v0.0.2"),
		},
		{
			name:      "with rc",
			expect:    semver.MustParse("v0.0.3-rc1"),
			includeRC: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g, err := Open("testdata/repo", Config{
				Prefix:    "v",
				IncludeRC: test.includeRC,
			})
			require.NoError(t, err)
			require.Equal(t, test.expect, g.Highest())
		})
	}
}

func TestRC_Increment(t *testing.T) {
	type increment struct {
		major bool
		minor bool
		patch bool
		dev   bool
		rc    bool
	}

	tests := []struct {
		name      string
		repo      string
		expect    semver.Version
		below     *semver.Version
		prefix    string
		includeRC bool
		increment increment
	}{
		{
			name:      "no rc, increment patch",
			repo:      "repo",
			expect:    semver.MustParse("v0.0.3"),
			increment: increment{patch: true},
		},
		{
			name:      "no rc, increment minor",
			repo:      "repo",
			expect:    semver.MustParse("v0.1.0"),
			increment: increment{minor: true},
		},
		{
			name:      "no rc, increment major",
			repo:      "repo",
			expect:    semver.MustParse("v1.0.0"),
			increment: increment{major: true},
		},
		{
			name:      "with rc, increment patch",
			repo:      "repo-rc",
			expect:    semver.MustParse("v0.0.3-rc1"),
			below:     ptr(semver.MustParse("v0.1.0")),
			includeRC: true,
			increment: increment{patch: true, rc: true},
		},
		{
			name:      "with rc, increment minor",
			repo:      "repo-rc",
			expect:    semver.MustParse("v0.1.0-rc1"),
			below:     ptr(semver.MustParse("v0.0.3")),
			includeRC: true,
			increment: increment{minor: true, rc: true},
		},
		{
			name:      "with rc, increment minor where exists",
			repo:      "repo-rc",
			expect:    semver.MustParse("v0.4.0-rc1"),
			below:     ptr(semver.MustParse("v0.9.0")),
			includeRC: true,
			increment: increment{minor: true, rc: true},
		},
		{
			name:      "with rc, increment major",
			repo:      "repo-rc",
			expect:    semver.MustParse("v1.0.0-rc1"),
			below:     ptr(semver.MustParse("v0.1.0")),
			includeRC: true,
			increment: increment{major: true, rc: true},
		},
		{
			name:      "with rc, increment rc",
			repo:      "repo-rc",
			expect:    semver.MustParse("v0.0.2-rc1"),
			below:     ptr(semver.MustParse("v0.1.0")),
			includeRC: true,
			increment: increment{rc: true},
		},
		{
			name:      "with rc, increment existing rc",
			repo:      "repo",
			expect:    semver.MustParse("v0.0.3-rc2"),
			below:     ptr(semver.MustParse("v0.1.0")),
			includeRC: true,
			increment: increment{rc: true},
		},
		{
			name:      "no rc, increment dev (snapshot)",
			repo:      "repo-rc",
			expect:    semver.MustParse("v0.0.2-snapshot-cf85392"),
			below:     ptr(semver.MustParse("v0.1.0")),
			increment: increment{dev: true},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			g, err := Open(fmt.Sprintf("testdata/%s", test.repo), Config{
				Prefix:    "v",
				Below:     test.below,
				IncludeRC: test.includeRC,
			})
			require.NoError(t, err)
			got, err := g.Increment(test.increment.major, test.increment.minor, test.increment.patch, test.increment.dev, test.increment.rc)
			require.NoError(t, err)
			require.Equal(t, test.expect, got)
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

func ptr[T any](v T) *T {
	return &v
}
