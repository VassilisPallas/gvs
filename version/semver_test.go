package version_test

import (
	"errors"
	"testing"

	"github.com/VassilisPallas/gvs/version"
	"github.com/google/go-cmp/cmp"
)

func intToUnsigned(num int64) *uint64 {
	res := uint64(num)
	return &res
}

func TestParseSemver(t *testing.T) {
	testCases := []struct {
		testTitle      string
		version        string
		expectedSemver *version.Semver
		expectedError  error
	}{
		{
			testTitle:      "should return an error when regex does not compile the passed version",
			version:        "some-version",
			expectedSemver: &version.Semver{},
			expectedError:  errors.New("invalid Go version"),
		},
		{
			testTitle:      "should parse only the major version",
			version:        "1",
			expectedSemver: &version.Semver{Major: intToUnsigned(1)},
			expectedError:  nil,
		},
		{
			testTitle:      "should parse minor version",
			version:        "1.23",
			expectedSemver: &version.Semver{Major: intToUnsigned(1), Minor: intToUnsigned(23)},
			expectedError:  nil,
		},
		{
			testTitle:      "should parse patch version",
			version:        "1.23.3",
			expectedSemver: &version.Semver{Major: intToUnsigned(1), Minor: intToUnsigned(23), Patch: intToUnsigned(3)},
			expectedError:  nil,
		},
		{
			testTitle:      "should parse rc version",
			version:        "1.23rc2",
			expectedSemver: &version.Semver{Major: intToUnsigned(1), Minor: intToUnsigned(23), ReleaseCandidate: intToUnsigned(2)},
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			semver := &version.Semver{}
			err := version.ParseSemver(tc.version, semver)

			if tc.expectedError == nil && err != nil {
				t.Errorf("error should be nil, instead got %q", err.Error())
				return
			}

			if tc.expectedError != nil && err.Error() != tc.expectedError.Error() {
				t.Errorf("error should be %q, instead got %q", tc.expectedError.Error(), err.Error())
				return
			}

			if !cmp.Equal(tc.expectedSemver, semver) {
				t.Errorf("Wrong object received, got=%s", cmp.Diff(tc.expectedSemver, semver))
			}
		})
	}
}

func TestGetVersion(t *testing.T) {
	testCases := []struct {
		testTitle       string
		semver          *version.Semver
		expectedVersion string
	}{
		{
			testTitle:       "should return an empty string when major is not specified",
			semver:          &version.Semver{},
			expectedVersion: "",
		},
		{
			testTitle:       "should return the major version",
			semver:          &version.Semver{Major: intToUnsigned(1)},
			expectedVersion: "1",
		},
		{
			testTitle:       "should return the major and the minor version",
			semver:          &version.Semver{Major: intToUnsigned(1), Minor: intToUnsigned(23)},
			expectedVersion: "1.23",
		},
		{
			testTitle:       "should return the major, minor and the patch version",
			semver:          &version.Semver{Major: intToUnsigned(1), Minor: intToUnsigned(23), Patch: intToUnsigned(3)},
			expectedVersion: "1.23.3",
		},
		{
			testTitle:       "should return the major, minor and the rc version",
			semver:          &version.Semver{Major: intToUnsigned(1), Minor: intToUnsigned(23), ReleaseCandidate: intToUnsigned(2)},
			expectedVersion: "1.23rc2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			res := tc.semver.GetVersion()

			if res != tc.expectedVersion {
				t.Errorf("version should be %q, instead got %q", tc.expectedVersion, res)
			}
		})
	}
}
