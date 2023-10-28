// Package version provides an interface to make handle
// the CLI logic for the versions.
package version

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Semver contains the version split into semantic version structure.
// Semver can have either a patch or a release candidate value, not both at the same time.
type Semver struct {
	// Contains the major version of the semantic structure.
	Major *uint64

	// Contains the minor version of the semantic structure.
	Minor *uint64

	// Contains the patch version of the semantic structure.
	Patch *uint64

	// Contains the release candidate number of the semantic structure.
	ReleaseCandidate *uint64
}

// GetVersion returns back the stringified representation of the semantic version structure.
func (s Semver) GetVersion() string {
	semver := ""

	if s.Major == nil {
		return ""
	}

	semver += fmt.Sprint(*s.Major)

	if s.Minor != nil {
		semver += fmt.Sprintf(".%d", *s.Minor)
	}

	if s.Patch != nil {
		semver += fmt.Sprintf(".%d", *s.Patch)
	} else if s.ReleaseCandidate != nil {
		semver += fmt.Sprintf("rc%d", *s.ReleaseCandidate)
	}

	return semver
}

// parseNumber converts a string into *uint64.
// In case of an error while converting the value, parseNumber return nil.
func parseNumber(str string) *uint64 {
	num, err := strconv.ParseUint(str, 10, 64)

	if err != nil {
		return nil
	}

	return &num
}

// ParseSemver parses the given stringified version into a semantic version structure.
//
// If the parse is successful it will store it
// in the value pointed to by semver.
//
// If the parse fails, ParseSemver will return an error.
func ParseSemver(version string, semver *Semver) error {
	var major *uint64
	var minor *uint64
	var patch *uint64
	var rc *uint64

	r := regexp.MustCompile(`(\d{1,2})(\.?\d{1,2})?(\.?\d{1,2}|rc\d{1,2})?`)

	if !r.MatchString(version) {
		return errors.New("invalid Go version")
	}

	groups := r.FindStringSubmatch(version)[1:]
	major = parseNumber(groups[0])
	minor = parseNumber(strings.TrimPrefix(groups[1], "."))

	if strings.HasPrefix(groups[2], "rc") {
		rc = parseNumber(strings.TrimPrefix(groups[2], "rc"))
	} else {
		patch = parseNumber(strings.TrimPrefix(groups[2], "."))
	}

	semver.Major = major
	semver.Minor = minor
	semver.Patch = patch
	semver.ReleaseCandidate = rc

	return nil
}
