package semver

import (
	"testing"
)

func TestNewVersion(t *testing.T) {
	tests := map[string]struct {
		in  string
		out *Version
		err error
	}{
		"none": {
			in:  "none",
			err: ErrTooFewElements,
		},
		"simple": {
			in:  "1.2.3",
			err: nil,
			out: &Version{1, 2, 3, 0},
		},
		"simple with prefix": {
			in:  "v1.2.3",
			err: nil,
			out: &Version{1, 2, 3, 0},
		},
		"one element": {
			in:  "1",
			err: ErrTooFewElements,
		},
		"two elements": {
			in:  "1.2",
			err: ErrTooFewElements,
		},
		"with pre-release": {
			in:  "v1.2.3-4",
			out: &Version{1, 2, 3, 4},
		},
		"with extended pre-release": {
			in:  "v1.2.3-4-g42e1055",
			out: &Version{1, 2, 3, 4},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res, err := NewVersion(tc.in)
			if err != tc.err {
				t.Fatalf("expected error value of '%v' but got '%v'", tc.err, err)
			}
			if res == nil {
				if tc.out != nil {
					t.Fatalf("expected nil result but got '%v'", tc.out)
				}
				return
			}
			if res != nil && tc.out == nil {
				t.Fatalf("expected non-nil result")
			}
			if res.major != tc.out.major {
				t.Fatalf("expected major value of '%v' but got '%v'", tc.out.major, res.major)
			}
			if res.minor != tc.out.minor {
				t.Fatalf("expected minor value of '%v' but got '%v'", tc.out.minor, res.minor)
			}
			if res.patch != tc.out.patch {
				t.Fatalf("expected patch value of '%v' but got '%v'", tc.out.patch, res.patch)
			}
			if res.pre != tc.out.pre {
				t.Fatalf("expected pre-release value of '%v' but got '%v'", tc.out.pre, res.pre)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	tests := map[string]struct {
		in1 Version
		in2 Version
		out int
	}{
		"greater major": {
			in1: Version{2, 0, 0, 0},
			in2: Version{1, 2, 3, 4},
			out: 1,
		},
		"greater minor": {
			in1: Version{1, 3, 0, 0},
			in2: Version{1, 2, 3, 4},
			out: 1,
		},
		"greater patch": {
			in1: Version{1, 2, 4, 0},
			in2: Version{1, 2, 3, 4},
			out: 1,
		},
		"greater pre": {
			in1: Version{1, 2, 3, 5},
			in2: Version{1, 2, 3, 4},
			out: 1,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			res := tc.in1.Compare(tc.in2)
			if res != tc.out {
				t.Fatalf("expected '%v' but got '%v'", tc.out, res)
			}
			res = tc.in2.Compare(tc.in1)
			if res != -tc.out {
				t.Fatalf("expected '%v' but got '%v'", -tc.out, res)
			}
		})
	}
}
