package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func main() {

	// get which image we want to update
	args := os.Args
	if len(args) != 2 {
		fmt.Println("please specify which docker image you would like to check")
		fmt.Println("Usage: refresh_image IMAGE_NAME")
		os.Exit(1)
	}
	dockerImage := args[1]

	// try and connect to docker
	cli, err := client.NewEnvClient()
	if err != nil {
		// does not seem to be running
		fmt.Println(err)
		os.Exit(1)
	}

	// get list of local images
	listFilters := filters.NewArgs()
	listFilters.Add("reference", dockerImage)
	ctx := context.Background()
	// images, err := cli.ImageList(ctx, types.ImageListOptions{Filters: listFilters, All: false})
	images, err := cli.ImageList(ctx, types.ImageListOptions{Filters: listFilters})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(images) == 0 {
		fmt.Printf("no local images match: %s", dockerImage)
		os.Exit(4)
	}

	// go through each image
	var allImageTags []string

	fmt.Println("Checking for ", dockerImage, ":")
	fmt.Println("  Local:")
	for _, image := range images {

		// see when image was created
		var imageTags []string
		imageCreated := time.Unix(image.Created, 0)

		// get name from the tags, which is actually name:tag (ie node:latest)
		for _, tag := range image.RepoTags {

			// note may have several images for a specific project
			// python:3-7 & python:3.7-slim

			tagParts := strings.SplitN(tag, ":", 2)
			imageTags = append(imageTags, tagParts[1])
			allImageTags = append(allImageTags, tagParts[1])
		}

		// ignore images that are unlabeled
		// if strings.Compare(imageName, "<none>") == 0 {
		// 	continue
		// }

		if len(imageTags) > 0 {
			fmt.Printf("    %-25s %s  created: %s\n", imageTags[0], image.ID[7:19], imageCreated.Format("01-02-2006 15:04"))
		}
	}

	// now see what we can find it on docker hub
	arch := getMachineArchitecture()

	// see if we can find the source image on docker.io
	found := checkDockerHubForImage(cli, dockerImage)
	if found {

		// get list of tags for that image
		tags, err := getTagsForImage(dockerImage)
		if err != nil {
			fmt.Printf("could not find tag on docker hub: %s", err)
			// continue
		}

		fmt.Println("  Docker Hub:")
		// go through tags
		for _, tag := range tags {
			// find image hash
			image, err := getImageForArchitecture(tag.Images, arch)
			if err == nil {
				imageTimestamp, err := time.Parse(time.RFC3339, image.LastPushed)
				if err == nil {
					fmt.Printf("    %-25s %s  created: %s\n", tag.Name, image.Digest[7:19], imageTimestamp.Format("01-02-2006 15:04"))
				}
			}
		}
	} else {
		fmt.Printf("could not find %s on docker hub", dockerImage)
	}

}
