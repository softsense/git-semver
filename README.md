# git-semver
Create semantic versions based on git tags

## Usage:
```
A tool for bumping semantic versions based on git tags.

Usage:
  git-semver [flags]

Flags:
      --below string    only look at tags below version
  -h, --help            help for git-semver
      --major           bump major version
      --minor           bump minor version
      --patch           bump patch version (default true)
      --rc              bump rc version. will bump other version if an rc does
                        not already exist.
      --repo string     path to git repository (default "./")
      --snapshot        set snapshot version
      --prefix string   use a prefix
```
