package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr string
	}{
		{
			name: "merges the three collections",
			args: []string{"-c1", "9,7,5,3", "-c2", "0,2,4,8", "-c3", "1,6,10"},
			want: "[0 1 2 3 4 5 6 7 8 9 10]\n",
		},
		{
			name: "spaces and missing flags are fine",
			args: []string{"-c1", " 3, 1 ", "-c2", "0"},
			want: "[0 1 3]\n",
		},
		{
			name: "no input prints empty",
			args: nil,
			want: "[]\n",
		},
		{
			name:    "not an integer",
			args:    []string{"-c2", "0,x"},
			wantErr: `-c2: "x" is not an integer`,
		},
		{
			name: "duplicates across collections",
			args: []string{"-c1", "9,5,3,1", "-c2", "0,2,8,10", "-c3", "1,2,6,11"},
			want: "[0 1 1 2 2 3 5 6 8 9 10 11]\n",
		},
		{
			name:    "unknown flag",
			args:    []string{"-c4", "1"},
			wantErr: "flag provided but not defined",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var out bytes.Buffer
			err := run(tt.args, &out)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if out.String() != tt.want {
				t.Errorf("output = %q, want %q", out.String(), tt.want)
			}
		})
	}
}
