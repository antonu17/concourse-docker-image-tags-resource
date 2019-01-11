package resource

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/genuinetools/reg/registry"
	"github.com/genuinetools/reg/repoutils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Setup(source Source) error {
	if source.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if source.Repository == "" {
		return errors.New("repository is required")
	}
	Regexp = regexp.MustCompile(source.Regexp)
	return nil
}

func GetTags(source Source) ([]string, error) {
	image, err := registry.ParseImage(source.Repository)
	if err != nil {
		return nil, err
	}

	// Use the auth-url domain if provided.
	authDomain := source.AuthURL
	if authDomain == "" {
		authDomain = image.Domain
	}

	auth, err := repoutils.GetAuthConfig(source.Username, source.Password, authDomain)
	if err != nil {
		return nil, err
	}

	// Prevent non-ssl unless explicitly forced
	if !source.ForceNonSSL && strings.HasPrefix(auth.ServerAddress, "http:") {
		return nil, fmt.Errorf("Attempted to use insecure protocol! Use force-non-ssl option to force")
	}

	// Create the registry client.
	r, err := registry.New(auth, registry.Opt{
		Insecure: source.Insecure,
		Debug:    source.Debug,
		SkipPing: source.SkipPing,
		NonSSL:   source.ForceNonSSL,
		Timeout:  source.Timeout,
	})
	if err != nil {
		return nil, err
	}

	return r.Tags(image.Path)
}
