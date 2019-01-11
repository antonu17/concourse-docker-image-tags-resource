package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/antonu17/concourse-docker-image-tags-resource/pkg/resource"
	"github.com/sirupsen/logrus"
)

type InParams struct {
	Last int `json:"last"`
}

type InRequest struct {
	Source  resource.Source  `json:"source"`
	Version resource.Version `json:"version"`
	Params  InParams         `json:"params"`
}

type InResponse struct {
	Version  resource.Version    `json:"version"`
	Metadata []resource.Metadata `json:"metadata"`
}

func main() {
	logrus.SetOutput(os.Stderr)
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	req := InRequest{
		Source: resource.SourceDefaults,
		Params: InParams{
			Last: 1,
		},
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

	if err := req.Version.Parse(); err != nil {
		logrus.WithError(err).Fatal("invalid version given")
		return
	}

	if len(os.Args) < 2 {
		logrus.Fatal("destination path not specified")
		return
	}

	tags, err := resource.GetTags(req.Source)
	if err != nil {
		logrus.WithError(err).Fatal("could not get tags")
		return
	}

	var last []string
	if req.Params.Last > 1 {
		var versions []resource.Version
		for _, tag := range tags {
			v := resource.Version{Tag: tag}
			if err := v.Parse(); err != nil {
				logrus.WithField("tag", tag).WithError(err).Error("could not parse tag")
				continue
			}
			if v.Version.LTE(req.Version.Version) {
				versions = append(versions, v)
			}
		}
		resource.Sort(versions)
		for i := req.Params.Last; i != 0; i-- {
			last = append(last, versions[len(versions)-i].Tag)
		}
	} else {
		last = append(last, req.Version.Tag)
	}

	dest := os.Args[1]
	if err := ioutil.WriteFile(path.Join(dest, "tag"), []byte(req.Version.Tag), 0644); err != nil {
		logrus.WithError(err).Fatal("could not write to tag file")
	}
	if err := ioutil.WriteFile(path.Join(dest, "last"), []byte(strings.Join(last, "\n")), 0644); err != nil {
		logrus.WithError(err).Fatal("could not write to tag file")
	}

	response := InResponse{
		Version:  req.Version,
		Metadata: []resource.Metadata{},
	}

	if err = json.NewEncoder(os.Stdout).Encode(response); err != nil {
		logrus.WithError(err).Fatal("could not encode output")
	}
}
