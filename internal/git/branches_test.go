package git

import (
	"reflect"
	"testing"
)

func TestParseBranches(t *testing.T) {
	tests := []struct {
		name   string
		output string
		want   []string
	}{
		{
			name:   "empty output",
			output: "",
			want:   []string{},
		},
		{
			name:   "single branch",
			output: "  origin/main",
			want:   []string{"main"},
		},
		{
			name: "multiple branches",
			output: `  origin/main
  origin/develop
  origin/feature/test`,
			want: []string{"main", "develop", "feature/test"},
		},
		{
			name: "with HEAD reference",
			output: `  origin/HEAD -> origin/main
  origin/main
  origin/develop`,
			want: []string{"main", "develop"},
		},
		{
			name: "already stripped prefix",
			output: `  main
  develop`,
			want: []string{"main", "develop"},
		},
		{
			name:   "whitespace variations",
			output: "\n  origin/main  \n\n  origin/develop\n",
			want:   []string{"main", "develop"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := _parseBranches(tt.output)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("_parseBranches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeduplicateBranches(t *testing.T) {
	tests := []struct {
		name     string
		branches []string
		want     []string
	}{
		{
			name:     "no duplicates",
			branches: []string{"main", "develop", "feature"},
			want:     []string{"main", "develop", "feature"},
		},
		{
			name:     "with duplicates",
			branches: []string{"main", "develop", "main", "feature", "develop"},
			want:     []string{"main", "develop", "feature"},
		},
		{
			name:     "all duplicates",
			branches: []string{"main", "main", "main"},
			want:     []string{"main"},
		},
		{
			name:     "empty list",
			branches: []string{},
			want:     []string{},
		},
		{
			name:     "single branch",
			branches: []string{"main"},
			want:     []string{"main"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := _deduplicateBranches(tt.branches)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("_deduplicateBranches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaultBranches(t *testing.T) {
	got := _defaultBranches()

	expectedBranches := []string{"main", "master", "develop"}
	if !reflect.DeepEqual(got, expectedBranches) {
		t.Errorf("_defaultBranches() = %v, want %v", got, expectedBranches)
	}

	if len(got) != 3 {
		t.Errorf("_defaultBranches() length = %d, want 3", len(got))
	}
}
