package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type LocalDocker struct {
	cli    *client.Client
	images []types.ImageSummary
}

// Init will connect to local docker daemon
func (ld *LocalDocker) Init() error {

	// try and connect to docker
	cli, err := client.NewEnvClient()
	if err != nil {
		// does not seem to be running
		return err
	}
	ld.cli = cli

	return nil
}

// Find will search local docker images for requested image
func (ld *LocalDocker) Find(image string) error {

	// get local docker images that match requested name
	listFilters := filters.NewArgs()
	listFilters.Add("reference", image)

	ctx := context.Background()
	// images, err := cli.ImageList(ctx, types.ImageListOptions{Filters: listFilters, All: false})
	images, err := ld.cli.ImageList(ctx, types.ImageListOptions{Filters: listFilters})
	if err != nil {
		return err
	}

	if len(images) == 0 {
		msg := fmt.Sprintf("no local images match: %s", image)
		return errors.New(msg)
	}

	ld.images = images

	return nil
}

// ListTags will print key info on each tag for current image
func (ld *LocalDocker) ListTags() {

	fmt.Println("  Local Images:")

	// for each image tag
	for _, image := range ld.images {

		var tags []string

		imageTimestamp := time.Unix(image.Created, 0)

		// get name from the tags, which is actually name:tag (ie node:latest)
		for _, fullTag := range image.RepoTags {

			// note may have several images for a specific project
			// python:3-7 & python:3.7-slim

			tagParts := strings.SplitN(fullTag, ":", 2)
			tag := tagParts[1]

			// ignore images that are unlabeled
			// if strings.Compare(imageName, "<none>") == 0 {
			// 	continue
			// }

			tags = append(tags, tag)
		}

		if len(tags) > 0 {
			fmt.Printf("    %-25s %s  created: %s\n", tags[0], image.ID[7:19], imageTimestamp.Format("01-02-2006 15:04"))
		}

	}

}
