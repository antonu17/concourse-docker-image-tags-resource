package main

import (
	"encoding/json"
	"os"

	"github.com/antonu17/concourse-docker-image-tags-resource/pkg/resource"
	"github.com/sirupsen/logrus"
)

type CheckRequest struct {
	Source  resource.Source  `json:"source"`
	Version resource.Version `json:"version"`
}

func main() {
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	req := CheckRequest{
		Source:  resource.SourceDefaults,
		Version: resource.Version{},
	}

	decoder := json.NewDecoder(os.Stdin)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&req)
	if err != nil {
		logrus.Fatalf("invalid payload: %s", err)
		return
	}

	if err := resource.Setup(req.Source); err != nil {
		logrus.WithError(err).Fatal("could not setup resource")
	}

	firstCheck := req.Version.Tag == ""
	if !firstCheck {
		if err := req.Version.Parse(); err != nil {
			logrus.WithError(err).Fatal("invalid version given")
			return
		}
	}

	tags, err := resource.GetTags(req.Source)
	if err != nil {
		logrus.WithError(err).Fatal("could not get tags")
		return
	}

	var response []resource.Version

	for _, tag := range tags {
		v := resource.Version{Tag: tag}
		if err := v.Parse(); err != nil {
			logrus.WithField("tag", tag).WithError(err).Error("could not parse tag")
			continue
		}
		if firstCheck {
			if len(response) == 0 {
				response = append(response, v)
			} else if v.Version.GT(response[0].Version) {
				response[0] = v
			}
		} else {
			if v.Version.GTE(req.Version.Version) {
				response = append(response, v)
			}
		}
	}
	resource.Sort(response)

	if err = json.NewEncoder(os.Stdout).Encode(response); err != nil {
		logrus.WithError(err).Fatal("could not encode output")
	}
}
