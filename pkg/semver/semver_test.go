package semver

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

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
