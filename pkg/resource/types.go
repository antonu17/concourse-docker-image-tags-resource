package resource

import (
	"regexp"
	"sort"
	"time"

	"github.com/blang/semver"
	"github.com/pkg/errors"
)

var Regexp *regexp.Regexp

var SourceDefaults = Source{
	Repository:  "",
	Regexp:      "^v?([.0-9]*)$",
	Insecure:    false,
	ForceNonSSL: false,
	SkipPing:    false,
	Timeout:     time.Minute,
	AuthURL:     "",
	Username:    "",
	Password:    "",
	Debug:       false,
}

type Source struct {
	Repository  string        `json:"repository"`
	Regexp      string        `json:"regexp"`
	Insecure    bool          `json:"insecure"`
	ForceNonSSL bool          `json:"force_non_ssl"`
	SkipPing    bool          `json:"skip_ping"`
	Timeout     time.Duration `json:"timeout"`
	AuthURL     string        `json:"auth_url"`
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	Debug       bool          `json:"debug"`
}

type Version struct {
	Tag     string         `json:"tag,omitempty"`
	Version semver.Version `json:"-"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (v *Version) Parse() (err error) {
	match := Regexp.FindStringSubmatch(v.Tag)
	if match == nil {
		return errors.New("regexp doesn't match")
	}
	if len(match) > 1 {
		v.Version, err = semver.ParseTolerant(match[1])
	} else {
		v.Version, err = semver.ParseTolerant(match[0])
	}
	return
}

// Versions represents multiple versions.
type Versions []Version

// Len returns length of version collection
func (s Versions) Len() int {
	return len(s)
}

// Swap swaps two versions inside the collection by its indices
func (s Versions) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less checks if version at index i is less than version at index j
func (s Versions) Less(i, j int) bool {
	return s[i].Version.LT(s[j].Version)
}

// Sort sorts a slice of versions
func Sort(versions []Version) {
	sort.Sort(Versions(versions))
}
