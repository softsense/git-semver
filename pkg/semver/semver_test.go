package semver

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    Version
		wantErr bool
	}{
		{
			name:    "1.0.0",
			version: "1.0.0",
			want: Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			wantErr: false,
		},
		{
			name:    "0.0.1",
			version: "0.0.1",
			want: Version{
				Major: 0,
				Minor: 0,
				Patch: 1,
			},
			wantErr: false,
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
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
