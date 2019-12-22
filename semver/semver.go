package semver

import (
	"errors"
	"strconv"
	"strings"
)

var ErrTooFewElements = errors.New("version contains less than 3 elements")

type Version struct {
	major, minor, patch, pre uint64
}

func NewVersion(v string) (*Version, error) {
	elements := strings.Split(v, ".")
	if len(elements) != 3 {
		return nil, ErrTooFewElements
	}
	if elements[0][0] == 'v' {
		elements[0] = elements[0][1:]
	}
	major, err := strconv.ParseUint(elements[0], 10, 64)
	if err != nil {
		return nil, err
	}
	minor, err := strconv.ParseUint(elements[1], 10, 64)
	if err != nil {
		return nil, err
	}
	patchElements := strings.Split(elements[2], "-")
	patch, err := strconv.ParseUint(patchElements[0], 10, 64)
	if err != nil {
		return nil, err
	}
	pre := uint64(0)
	if len(patchElements) > 1 {
		pre, err = strconv.ParseUint(patchElements[1], 10, 64)
		if err != nil {
			return nil, err
		}
	}

	return &Version{
		major: major,
		minor: minor,
		patch: patch,
		pre:   pre,
	}, nil
}

func (v Version) Compare(o Version) int {
	if v.major < o.major {
		return -1
	}
	if v.major > o.major {
		return 1
	}
	if v.minor < o.minor {
		return -1
	}
	if v.minor > o.minor {
		return 1
	}
	if v.patch < o.patch {
		return -1
	}
	if v.patch > o.patch {
		return 1
	}
	if v.pre < o.pre {
		return -1
	}
	if v.pre > o.pre {
		return 1
	}
	return 0
}
