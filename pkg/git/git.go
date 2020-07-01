package git

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/softsense/git-semver/pkg/semver"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

type Config struct {
	// Prefix to add to version strings
	Prefix string
}

type Git struct {
	highest semver.Version
	repo    *git.Repository
	cfg     Config
}

func Open(path string, cfg Config) (*Git, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open git repo %s", path)
	}

	var highest semver.Version
	highest.Prefix = cfg.Prefix

	all := make(map[string]semver.Version)

	tagrefs, err := r.Tags()
	if err != nil {
		return nil, errors.Wrap(err, "failed to list tags")
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
			return nil
		}
		if n.GT(highest) {
			highest = n
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to loop over tags")
	}

	g := &Git{
		highest: highest,
		repo:    r,
		cfg:     cfg,
	}

	return g, nil
}

func (g *Git) Increment(major, minor, patch, dev bool) (semver.Version, error) {
	newVersion, err := semver.Parse(g.highest.String())
	if err != nil {
		return semver.Version{}, errors.Wrapf(err, "failed to create version")
	}

	if patch {
		if err := newVersion.IncrementPatch(); err != nil {
			return semver.Version{}, errors.Wrap(err, "failed to increment")
		}
	}

	if minor {
		if err := newVersion.IncrementMinor(); err != nil {
			return semver.Version{}, errors.Wrap(err, "failed to increment")
		}
	}

	if major {
		if err := newVersion.IncrementMajor(); err != nil {
			return semver.Version{}, errors.Wrap(err, "failed to increment")
		}
	}

	if dev {
		head, err := g.repo.Head()
		if err != nil {
			return semver.Version{}, errors.Wrap(err, "failed to get repo head")
		}
		snapshot, err := semver.NewPRVersion(fmt.Sprintf("snapshot-%s", head.Hash().String()[:7]))
		if err != nil {
			return semver.Version{}, errors.Wrap(err, "failed to build snapshot version")
		}
		newVersion.Pre = []semver.PRVersion{snapshot}
	}

	return newVersion, nil
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
