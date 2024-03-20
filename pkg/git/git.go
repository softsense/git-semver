package git

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/softsense/git-semver/pkg/semver"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var (
	prNumFromCommit = regexp.MustCompile(`\(#([0-9]+)\)($|\n)`)
	gitHostAndPath  = regexp.MustCompile(`git@(.+):(.*)\.git`)
)

type Config struct {
	// Prefix to add to version strings
	Prefix string

	// Only look at tags below version
	Below *semver.Version

	// Include ReleaseCandidate Version
	IncludeRC bool
}

type Git struct {
	highest semver.Version
	repo    *git.Repository
	cfg     Config
}

func Open(path string, cfg Config) (*Git, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("open git repo %s: %w", path, err)
	}

	var highest semver.Version
	highest.Prefix = cfg.Prefix

	all := make(map[string]semver.Version)

	tagrefs, err := r.Tags()
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		n, err := parseTagRef(string(t.Name()))
		if err != nil {
			return nil
		}

		// only care about tags with the same prefix
		if n.Prefix != cfg.Prefix {
			return nil
		}

		v := format(n)
		allN, ok := all[v]
		if ok {
			if n.GT(allN) {
				all[v] = n
			}
		} else {
			all[v] = n
		}

		if len(n.Pre) > 0 {
			if !cfg.IncludeRC || !strings.HasPrefix(n.Pre[0].String(), "rc") {
				return nil
			}
		}
		if cfg.Below != nil && n.GTE(*cfg.Below) {
			return nil
		}
		if n.GT(highest) {
			highest = n
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("loop over tags: %w", err)
	}

	g := &Git{
		highest: highest,
		repo:    r,
		cfg:     cfg,
	}

	return g, nil
}

func (g *Git) Increment(major, minor, patch, dev bool, rc bool) (semver.Version, error) {
	newVersion, err := semver.Parse(g.highest.String())
	if err != nil {
		return semver.Version{}, fmt.Errorf("create version: %w", err)
	}

	if rc {
		found := false
		for i, pre := range newVersion.Pre {
			if strings.HasPrefix(pre.VersionStr, "rc") {
				rcNumStr := strings.Replace(pre.VersionStr, "rc", "", 1)
				rcNum, err := strconv.Atoi(rcNumStr)
				if err != nil {
					return semver.Version{}, fmt.Errorf("parse rc number: %w", err)
				}
				rcNum++
				newVersion.Pre[i], err = semver.NewPRVersion(fmt.Sprintf("rc%d", rcNum))
				if err != nil {
					return semver.Version{}, fmt.Errorf("create rc version: %w", err)
				}
				found = true
				patch = false
				minor = false
				major = false
				break
			}
		}
		if !found {
			rcVer, err := semver.NewPRVersion("rc1")
			if err != nil {
				return semver.Version{}, fmt.Errorf("build rc version: %w", err)
			}
			newVersion.Pre = []semver.PRVersion{rcVer}
		}
	}

	if patch {
		if err := newVersion.IncrementPatch(); err != nil {
			return semver.Version{}, fmt.Errorf("increment: %w", err)
		}
	}

	if minor {
		if err := newVersion.IncrementMinor(); err != nil {
			return semver.Version{}, fmt.Errorf("increment: %w", err)
		}
	}

	if major {
		if err := newVersion.IncrementMajor(); err != nil {
			return semver.Version{}, fmt.Errorf("increment: %w", err)
		}
	}

	if dev {
		head, err := g.repo.Head()
		if err != nil {
			return semver.Version{}, fmt.Errorf("get repo head: %w", err)
		}
		snapshot, err := semver.NewPRVersion(fmt.Sprintf("snapshot-%s", head.Hash().String()[:7]))
		if err != nil {
			return semver.Version{}, fmt.Errorf("build snapshot version: %w", err)
		}
		newVersion.Pre = []semver.PRVersion{snapshot}
	}

	return newVersion, nil
}

func (g *Git) History(prefix string) (string, error) {
	head, err := g.repo.Head()
	if err != nil {
		return "", fmt.Errorf("get head: %w", err)
	}

	var prevHash *plumbing.Hash
	prevRef, err := g.repo.Tag(g.highest.String())
	if err != nil {
		fmt.Printf("Tag %s not found, including the entire history\n", g.highest.String())
	} else {
		cIter, err := g.repo.Log(&git.LogOptions{From: prevRef.Hash()})
		if err != nil {
			return "", fmt.Errorf("get log from %s: %w", prevRef.Hash(), err)
		}
		c, err := cIter.Next()
		if err == nil {
			prevHash = &c.Hash
		}
	}

	cIter, err := g.repo.Log(&git.LogOptions{From: head.Hash()})
	if err != nil {
		return "", fmt.Errorf("get log from %s: %w", head.Hash(), err)
	}
	out := make([]string, 0)
	_ = cIter.ForEach(func(c *object.Commit) error {
		if prevHash != nil && c.Hash.String() == prevHash.String() {
			return errors.New("EOF")
		}
		msg := fmt.Sprintf("%s* %s %s\n", prefix, c.Hash.String()[:7], strings.ReplaceAll(strings.TrimSuffix(c.Message, "\n"), "\n", "\n  "))
		if prefix != "" {
			msg = strings.ReplaceAll(msg, "\n", fmt.Sprintf("\n%s", prefix))
		}
		msg += "\n"
		msg = insertPullRequestURL(msg, g)

		out = append(out, msg)

		return nil
	})

	return strings.Join(out, ""), nil
}

func (g *Git) Highest() semver.Version {
	return g.highest
}

// insertPullRequestURL replaces GitHub PR references in commit messages
// with the full URL to the PR.
func insertPullRequestURL(msg string, git *Git) string {
	remotes, err := git.repo.Remotes()
	if err != nil || len(remotes) < 1 {
		return msg
	}
	if len(remotes[0].Config().URLs) < 1 {
		return msg
	}
	url := remotes[0].Config().URLs[0]
	if !(strings.HasPrefix(url, "git@github.com") || strings.HasPrefix(url, "https://github.com")) {
		return msg
	}

	// Convert non-HTTPS URL to HTTPS
	url = gitHostAndPath.ReplaceAllString(url, "https://$1/$2")

	link := fmt.Sprintf("[(#$1)](%s/pull/$1)$2", url)
	msg = prNumFromCommit.ReplaceAllString(msg, link)

	return msg
}

func parseTagRef(t string) (semver.Version, error) {
	s := strings.Replace(t, "refs/tags/", "", 1)
	v, err := semver.Parse(s)
	if err != nil {
		return semver.Version{}, err
	}
	return v, nil
}

func format(v semver.Version) string {
	return fmt.Sprintf("%s%d.%d.%d", v.Prefix, v.Major, v.Minor, v.Patch)
}
