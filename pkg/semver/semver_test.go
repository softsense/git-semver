package semver

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	version_1_2_3 = Version{
		Major: 1,
		Minor: 2,
		Patch: 3,
	}
	version_1_2_2 = Version{
		Major: 1,
		Minor: 2,
		Patch: 2,
	}
	version_1_2_3_WithPR = Version{
		Major: 1,
		Minor: 2,
		Patch: 3,
		Pre: []PRVersion{
			{
				VersionNum: 4,
				IsNum:      true,
			},
			{
				VersionStr: "alpha5",
			},
		},
	}
	version_1_2_3_WithBuild = Version{
		Major: 1,
		Minor: 2,
		Patch: 3,
		Build: []string{
			"6",
			"alpha7",
		},
	}
	version_1_2_3_WithPRAndBuild = Version{
		Major: 1,
		Minor: 2,
		Patch: 3,
		Pre: []PRVersion{
			{
				VersionNum: 4,
				IsNum:      true,
			},
			{
				VersionStr: "alpha5",
			},
		},
		Build: []string{
			"6",
			"alpha7",
		},
	}
	version_1_2_2_WithPRAndBuild = Version{
		Major: 1,
		Minor: 2,
		Patch: 2,
		Pre: []PRVersion{
			{
				VersionNum: 4,
				IsNum:      true,
			},
			{
				VersionStr: "alpha5",
			},
		},
		Build: []string{
			"6",
			"alpha7",
		},
	}
)

func TestOperators(t *testing.T) {
	tests := []struct {
		name    string
		op      func(v Version, o Version) bool
		version Version
		other   Version
		want    bool
	}{
		// Equal
		{
			name:    "Equals: same version should be equal",
			op:      Version.Equals,
			version: version_1_2_3,
			other:   version_1_2_3,
			want:    true,
		},
		{
			name:    "Equals: same version should be equal - with prerelease and build",
			op:      Version.Equals,
			version: version_1_2_3_WithPRAndBuild,
			other:   version_1_2_3_WithPRAndBuild,
			want:    true,
		},
		{
			name:    "Equals: different version should be unequal",
			op:      Version.Equals,
			version: version_1_2_3,
			other:   version_1_2_2,
			want:    false,
		},
		{
			name:    "Equals: different prerelease should not be equal",
			op:      Version.Equals,
			version: version_1_2_3,
			other:   version_1_2_3_WithPR,
			want:    false,
		},

		// NE
		{
			name:    "NE: different version should be ne",
			op:      Version.NE,
			version: version_1_2_3,
			other:   version_1_2_2,
			want:    true,
		},
		{
			name:    "NE: same version should not be ne",
			op:      Version.NE,
			version: version_1_2_3,
			other:   version_1_2_3,
			want:    false,
		},

		// GT
		{
			name:    "GT: greater version should be gt",
			op:      Version.GT,
			version: version_1_2_3,
			other:   version_1_2_2,
			want:    true,
		},
		{
			name:    "GT: equal version should not be gt",
			op:      Version.GT,
			version: version_1_2_3,
			other:   version_1_2_3,
			want:    false,
		},

		// GTE
		{
			name:    "GTE: greater version should be gte",
			op:      Version.GTE,
			version: version_1_2_3,
			other:   version_1_2_2,
			want:    true,
		},
		{
			name:    "GTE: equal version should be gte",
			op:      Version.GTE,
			version: version_1_2_3,
			other:   version_1_2_3,
			want:    true,
		},
		{
			name:    "GTE: lesser version should not be gte",
			op:      Version.GTE,
			version: version_1_2_2,
			other:   version_1_2_3,
			want:    false,
		},

		// LT
		{
			name:    "LT: lesser version should be lt",
			op:      Version.LT,
			version: version_1_2_2,
			other:   version_1_2_3,
			want:    true,
		},
		{
			name:    "LT: equal version should not be lt",
			op:      Version.LT,
			version: version_1_2_3,
			other:   version_1_2_3,
			want:    false,
		},

		// LTE
		{
			name:    "LTE: lesser version should be lte",
			op:      Version.LTE,
			version: version_1_2_2,
			other:   version_1_2_3,
			want:    true,
		},
		{
			name:    "LTE: equal version should be lte",
			op:      Version.LTE,
			version: version_1_2_3,
			other:   version_1_2_3,
			want:    true,
		},
		{
			name:    "LTE: greater version should not be lte",
			op:      Version.LTE,
			version: version_1_2_3,
			other:   version_1_2_2,
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.op(tt.version, tt.other)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    Version
		wantErr error
	}{
		// valid
		{
			name:    "1.0.0",
			version: "1.0.0",
			want: Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
		},
		{
			name:    "0.0.1",
			version: "0.0.1",
			want: Version{
				Major: 0,
				Minor: 0,
				Patch: 1,
			},
		},
		{
			name:    "v0.1.0",
			version: "v0.1.0",
			want: Version{
				Major:  0,
				Minor:  1,
				Patch:  0,
				Prefix: "v",
			},
		},
		{
			name:    "0.0.1-alphanumeric1",
			version: "0.0.1-alphanumeric1",
			want: Version{
				Major: 0,
				Minor: 0,
				Patch: 1,
				Pre: []PRVersion{
					{
						VersionStr: "alphanumeric1",
					},
				},
			},
		},
		{
			name:    "0.0.1-1",
			version: "0.0.1-1",
			want: Version{
				Major: 0,
				Minor: 0,
				Patch: 1,
				Pre: []PRVersion{
					{
						VersionNum: 1,
						IsNum:      true,
					},
				},
			},
		},
		{
			name:    "0.0.1-1.alphanumeric2",
			version: "0.0.1-1.alphanumeric2",
			want: Version{
				Major: 0,
				Minor: 0,
				Patch: 1,
				Pre: []PRVersion{
					{
						VersionNum: 1,
						IsNum:      true,
					},
					{
						VersionStr: "alphanumeric2",
					},
				},
			},
		},
		{
			name:    "0.0.1+alphanumeric1",
			version: "0.0.1+alphanumeric1",
			want: Version{
				Major: 0,
				Minor: 0,
				Patch: 1,
				Build: []string{
					"alphanumeric1",
				},
			},
		},
		{
			name:    "0.0.1+1",
			version: "0.0.1+1",
			want: Version{
				Major: 0,
				Minor: 0,
				Patch: 1,
				Build: []string{
					"1",
				},
			},
		},
		{
			name:    "0.0.1+1.alphanumeric2",
			version: "0.0.1+1.alphanumeric2",
			want: Version{
				Major: 0,
				Minor: 0,
				Patch: 1,
				Build: []string{
					"1",
					"alphanumeric2",
				},
			},
		},
		{
			name:    "1.2.3-4.alphanumeric5+6.alphanumeric7",
			version: "1.2.3-4.alphanumeric5+6.alphanumeric7",
			want: Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Pre: []PRVersion{
					{
						VersionNum: 4,
						IsNum:      true,
					},
					{
						VersionStr: "alphanumeric5",
					},
				},
				Build: []string{
					"6",
					"alphanumeric7",
				},
			},
		},

		// invalid
		{
			name:    "empty",
			version: "",
			wantErr: errors.New("version string empty"),
		},
		{
			name:    "0",
			version: "0",
			wantErr: errors.New("no Major.Minor.Patch elements found"),
		},
		{
			name:    "0.0",
			version: "0.0",
			wantErr: errors.New("no Major.Minor.Patch elements found"),
		},
		// {
		// 	// TODO: This case is not defensively handled, it confuses it with a prefix.
		// 	name:    "major.0.1",
		// 	version: "major.0.1",
		// 	wantErr: errors.New(`only numbers`),
		// },
		{
			name:    "0NaN.0.1",
			version: "0NaN.0.1",
			wantErr: errors.New(`invalid character(s) found in major number "0NaN"`),
		},
		{
			name:    "01.0.1",
			version: "01.0.1",
			wantErr: errors.New(`major number must not contain leading zeroes "01"`),
		},
		{
			name:    "0.0NaN.1",
			version: "0.0NaN.1",
			wantErr: errors.New(`invalid character(s) found in minor number "0NaN"`),
		},
		{
			name:    "0.01.1",
			version: "0.01.1",
			wantErr: errors.New(`minor number must not contain leading zeroes "01"`),
		},
		{
			name:    "0.0.1NaN",
			version: "0.0.1NaN",
			wantErr: errors.New(`invalid character(s) found in patch number "1NaN"`),
		},
		{
			name:    "0.1.01",
			version: "0.1.01",
			wantErr: errors.New(`patch number must not contain leading zeroes "01"`),
		},
		{
			name:    "0.0.1-",
			version: "0.0.1-",
			wantErr: errors.New("prerelease is empty"),
		},
		{
			name:    "0.0.1-01",
			version: "0.0.1-01",
			wantErr: errors.New(`numeric PreRelease version must not contain leading zeroes "01"`),
		},
		{
			name:    "0.0.1-@",
			version: "0.0.1-@",
			wantErr: errors.New(`invalid character(s) found in prerelease "@"`),
		},
		{
			name:    "0.0.1+",
			version: "0.0.1+",
			wantErr: errors.New("build meta data is empty"),
		},
		{
			name:    "0.0.1+@",
			version: "0.0.1+@",
			wantErr: errors.New(`invalid character(s) found in build meta data "@"`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.version)

			if tt.wantErr != nil {
				require.EqualError(t, err, tt.wantErr.Error())
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestParseTolerant(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    Version
		wantErr error
	}{
		{
			name:    "Parses version that don't need modification",
			version: "1.2.3",
			want:    version_1_2_3,
		},
		{
			name:    "Trims leading spaces",
			version: "  1.2.3",
			want:    version_1_2_3,
		},
		{
			name:    "Trims trailing spaces",
			version: "1.2.3  ",
			want:    version_1_2_3,
		},
		{
			name:    "Removes leading zero in version components",
			version: "01.02.03",
			want:    version_1_2_3,
		},
		{
			name:    "Adds missing patch version component",
			version: "1.2",
			want: Version{
				Major: 1,
				Minor: 2,
				Patch: 0,
			},
		},
		{
			name:    "Adds missing minor and patch version components",
			version: "1",
			want: Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
		},
		{
			name:    "Does not accept short versions with prerelease or build metadata using '+' prefix",
			version: "1.0+alpha1",
			want: Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			wantErr: fmt.Errorf("short version cannot contain PreRelease/Build meta data"),
		},
		{
			name:    "Does not accept short versions with prerelease or build metadata using '+' prefix",
			version: "1.0-alpha1",
			want: Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			wantErr: fmt.Errorf("short version cannot contain PreRelease/Build meta data"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTolerant(tt.version)

			if tt.wantErr != nil {
				require.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.Equal(t, tt.want, got)
			}
		})
	}
}

func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		name                string
		version             string
		other               string
		additionalPRVersion string
		want                int
	}{
		{
			name:    "Equal versions should be equal - major",
			version: "1.0.0",
			other:   "1.0.0",
			want:    0,
		},
		{
			name:    "Equal versions should be equal - mixed",
			version: "1.2.3",
			other:   "1.2.3",
			want:    0,
		},
		{
			name:    "Equal PR versions should be equal",
			version: "1.2.3-rc1",
			other:   "1.2.3-rc1",
			want:    0,
		},

		{
			name:    "Higher major version is considered greater",
			version: "1.2.3",
			other:   "2.0.0",
			want:    -1,
		},
		{
			name:    "Higher minor version is considered greater",
			version: "1.2.3",
			other:   "1.3.0",
			want:    -1,
		},
		{
			name:    "Higher patch version is considered greater",
			version: "1.2.3",
			other:   "1.2.4",
			want:    -1,
		},
		{
			name:    "Higher PR version is considered greater",
			version: "1.2.3-rc1",
			other:   "1.2.3-rc2",
			want:    -1,
		},

		{
			name:    "Lower major version is considered lesser",
			version: "2.0.0",
			other:   "1.2.3",
			want:    1,
		},
		{
			name:    "Lower minor version is considered lesser",
			version: "1.3.0",
			other:   "1.2.3",
			want:    1,
		},
		{
			name:    "Lower patch version is considered lesser",
			version: "1.2.4",
			other:   "1.2.3",
			want:    1,
		},
		{
			name:    "Lower PR version is considered lesser",
			version: "1.2.3-rc2",
			other:   "1.2.3-rc1",
			want:    1,
		},

		{
			name:    "Regular version is considered greater than PR version",
			version: "1.2.3-rc1",
			other:   "1.2.3",
			want:    -1,
		},
		{
			name:    "PR version is considered lesser than regular version",
			version: "1.2.3",
			other:   "1.2.3-rc1",
			want:    1,
		},

		{
			name:                "Additional PR version is considered greater",
			version:             "1.2.3-rc1",
			other:               "1.2.3-rc1",
			additionalPRVersion: "rc2",
			want:                -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := New(tt.version)
			require.NoError(t, err)
			o, err := New(tt.other)
			require.NoError(t, err)
			if len(tt.additionalPRVersion) > 0 {
				pr, err := NewPRVersion(tt.additionalPRVersion)
				require.NoError(t, err)
				o.Pre = append(o.Pre, pr)
			}

			got := v.Compare(*o)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestVersion_Validate(t *testing.T) {
	tests := []struct {
		name    string
		version Version
		wantErr error
	}{
		{
			name: "Valid version without PR versions or build is valid",
			version: Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Pre: []PRVersion{
					{
						VersionNum: 4,
						IsNum:      true,
					},
					{
						VersionStr: "alpha5",
					},
				},
				Build: []string{
					"6",
					"alpha7",
				},
				Prefix: "v",
			},
		},
		{
			name: "Valid version with PR versions and build is valid",
			version: Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Pre: []PRVersion{
					{
						VersionNum: 4,
						IsNum:      true,
					},
					{
						VersionStr: "alpha5",
					},
				},
				Build: []string{
					"6",
					"alpha7",
				},
				Prefix: "v",
			},
		},
		{
			name: "Version with empty PR version is invalid",
			version: Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Pre: []PRVersion{
					{
						VersionStr: "",
					},
				},
			},
			wantErr: fmt.Errorf(`prerelease can not be empty ""`),
		},
		{
			name: "Version with PR version that consists of invalid characters is invalid",
			version: Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Pre: []PRVersion{
					{
						VersionStr: "@#!",
					},
				},
			},
			wantErr: fmt.Errorf(`invalid character(s) found in prerelease "@#!"`),
		},
		{
			name: "Version with PR version that ends with invalid characters is invalid",
			version: Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Pre: []PRVersion{
					{
						VersionStr: "alpha5@#!",
					},
				},
			},
			wantErr: fmt.Errorf(`invalid character(s) found in prerelease "alpha5@#!"`),
		},
		{
			name: "Version with empty build version is invalid",
			version: Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: []string{""},
			},
			wantErr: fmt.Errorf(`build meta data can not be empty ""`),
		},
		{
			name: "Version with build that consists of invalid characters is invalid",
			version: Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: []string{"@#!"},
			},
			wantErr: fmt.Errorf(`invalid character(s) found in build meta data "@#!"`),
		},
		{
			name: "Version with build that ends with invalid characters is invalid",
			version: Version{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: []string{"alpha7@#!"},
			},
			wantErr: fmt.Errorf(`invalid character(s) found in build meta data "alpha7@#!"`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.version.Validate()
			if tt.wantErr != nil {
				require.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestVersion_IncrementMajor(t *testing.T) {
	version, err := Parse("1.2.3")
	require.NoError(t, err)

	require.Equal(t, uint64(1), version.Major)
	require.Equal(t, uint64(2), version.Minor)
	require.Equal(t, uint64(3), version.Patch)

	version.IncrementMajor()

	require.Equal(t, uint64(2), version.Major, "major not incremented")
	require.Equal(t, uint64(0), version.Minor, "minor not reset")
	require.Equal(t, uint64(0), version.Patch, "patch not reset")
}

func TestVersion_IncrementMinor(t *testing.T) {
	version, err := Parse("1.2.3")
	require.NoError(t, err)

	require.Equal(t, uint64(1), version.Major)
	require.Equal(t, uint64(2), version.Minor)
	require.Equal(t, uint64(3), version.Patch)

	version.IncrementMinor()

	require.Equal(t, uint64(1), version.Major)
	require.Equal(t, uint64(3), version.Minor, "minor not incremented")
	require.Equal(t, uint64(0), version.Patch, "patch not reset")
}

func TestVersion_IncrementPatch(t *testing.T) {
	version, err := Parse("1.2.3")
	require.NoError(t, err)

	require.Equal(t, uint64(1), version.Major)
	require.Equal(t, uint64(2), version.Minor)
	require.Equal(t, uint64(3), version.Patch)

	version.IncrementPatch()

	t.Run("increments patch", func(t *testing.T) {
		require.Equal(t, uint64(1), version.Major)
		require.Equal(t, uint64(2), version.Minor)
		require.Equal(t, uint64(4), version.Patch, "patch not incremented")
	})
}

func TestVersion_String(t *testing.T) {
	type fields struct {
		Major  uint64
		Minor  uint64
		Patch  uint64
		Pre    []PRVersion
		Build  []string
		Prefix string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "1.0.0",
			fields: fields{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			want: "1.0.0",
		},
		{
			name: "0.0.1",
			fields: fields{
				Major: 0,
				Minor: 0,
				Patch: 1,
			},
			want: "0.0.1",
		},
		{
			name: "v0.1.0",
			fields: fields{
				Major:  0,
				Minor:  1,
				Patch:  0,
				Prefix: "v",
			},
			want: "v0.1.0",
		},
		{
			name: "v1.2.3-4.alpha5+6.alpha7",
			fields: fields{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Pre: []PRVersion{
					{
						VersionNum: 4,
						IsNum:      true,
					},
					{
						VersionStr: "alpha5",
					},
				},
				Build: []string{
					"6",
					"alpha7",
				},
				Prefix: "v",
			},
			want: "v1.2.3-4.alpha5+6.alpha7",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Version{
				Major:  tt.fields.Major,
				Minor:  tt.fields.Minor,
				Patch:  tt.fields.Patch,
				Pre:    tt.fields.Pre,
				Build:  tt.fields.Build,
				Prefix: tt.fields.Prefix,
			}
			if got := v.String(); got != tt.want {
				t.Errorf("Version.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPRVersion_Compare(t *testing.T) {
	tests := []struct {
		name    string
		version PRVersion
		other   PRVersion
		want    int
	}{

		{
			name: "Equal num pr version is equal",
			version: PRVersion{
				IsNum:      true,
				VersionNum: 1,
			},
			other: PRVersion{
				IsNum:      true,
				VersionNum: 1,
			},
			want: 0,
		},
		{
			name: "Equal non-num pr version is equal",
			version: PRVersion{
				IsNum:      false,
				VersionStr: "alpha5",
			},
			other: PRVersion{
				IsNum:      false,
				VersionStr: "alpha5",
			},
			want: 0,
		},
		{
			name: "Num pr version is less than non-num",
			version: PRVersion{
				IsNum:      true,
				VersionNum: 1,
			},
			other: PRVersion{
				IsNum:      false,
				VersionStr: "alpha5",
			},
			want: -1,
		},
		{
			name: "Non-num pr version is greater than num",
			version: PRVersion{
				IsNum:      false,
				VersionStr: "alpha5",
			},
			other: PRVersion{
				IsNum:      true,
				VersionNum: 1,
			},

			want: 1,
		},
		{
			name: "Greater num pr version is greater",
			version: PRVersion{
				IsNum:      true,
				VersionNum: 2,
			},
			other: PRVersion{
				IsNum:      true,
				VersionNum: 1,
			},
			want: 1,
		},
		{
			name: "Lesser num pr version is lesser",
			version: PRVersion{
				IsNum:      true,
				VersionNum: 1,
			},
			other: PRVersion{
				IsNum:      true,
				VersionNum: 2,
			},
			want: -1,
		},
		{
			name: "Longer non-num pr version is greater",
			version: PRVersion{
				IsNum:      false,
				VersionStr: "alphabeta",
			},
			other: PRVersion{
				IsNum:      false,
				VersionStr: "alpha5",
			},
			want: 1,
		},
		{
			name: "Shorter non-num pr version is lesser",
			version: PRVersion{
				IsNum:      false,
				VersionStr: "alpha5",
			},
			other: PRVersion{
				IsNum:      false,
				VersionStr: "alphabeta",
			},
			want: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.version.Compare(tt.other)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestNewBuildVersion(t *testing.T) {
	tests := []struct {
		name         string
		buildVersion string
		wantBuild    string
		wantErr      error
	}{
		{
			name:         "Correct build on valid buildversion",
			buildVersion: "alpha5",
			wantBuild:    "alpha5",
		},
		{
			name:         "Error on empty buildversion",
			buildVersion: "",
			wantErr:      fmt.Errorf("buildversion is empty"),
		},
		{
			name:         "Error on invalid characters",
			buildVersion: "alpha5!@#",
			wantErr:      fmt.Errorf(`invalid character(s) found in build meta data "alpha5!@#"`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBuildVersion(tt.buildVersion)
			require.Equal(t, tt.wantBuild, b)
			if tt.wantErr != nil {
				require.EqualError(t, err, tt.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
