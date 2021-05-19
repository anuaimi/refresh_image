package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// getArchitecture will return the hardware running on
func getArchitecture() string {
	// amd64, arm64 are main values
	return runtime.GOARCH
}

func main() {

	// log.PrintColors = true
	// log.PrintTimestamp = false

	// see if docker running locally
	// try and connect to it
	cli, err := client.NewEnvClient()
	if err != nil {
		// does not seem to be running
		fmt.Println(err)
		os.Exit(1)
	}

	// don't need credentials - so below is commented out

	// get an auth token from docker registry
	// authToken, err := getAuthToken()
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// fmt.Println(authToken)

	// username := ""
	// password := ""
	// loginValid, err := validateLogin(cli, username, password)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(3)
	// }

	// get list of local images
	ctx := context.Background()
	images, err := cli.ImageList(ctx, types.ImageListOptions{All: false})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// go through each image
	for _, image := range images {
		var imageName string
		var imageTags []string

		// get name from the tags, which is actually name:tag (ie node:latest)
		for _, tag := range image.RepoTags {

			// note may have several images for a specific project
			// python:3-7 & python:3.7-slim

			tagParts := strings.SplitN(tag, ":", 2)
			imageName = tagParts[0]
			imageTags = append(imageTags, tagParts[1])
		}

		// ignore images that are unlabeled
		if strings.Compare(imageName, "<none>") == 0 {
			continue
		}

		fmt.Printf("%-25s | %-25s\n",
			imageName,
			imageTags,
		)

		// fmt.Printf("%-25s | %-15s | %s | %-30s\n",
		// 	imageName,
		// 	imageTags,
		// 	image.ID[7:19],
		// 	time.Unix(image.Created, 0),
		// )

		// see if we can find the source image on docker.io
		fmt.Printf("looking for %s\n", imageName)

		found := checkDockerHubForImage(cli, imageName)
		if found {

			// get list of tags for that image
			tags, err := getTags(imageName)
			if err != nil {
				fmt.Printf("could not find tag on docker hub: %s", err)
				continue
			}

			// go through tags
			for _, localTag := range imageTags {

				// look for matching labels
				foundTag := false
				var tagInfo DockerImageTag
				for _, tag := range tags {
					fmt.Println(tag.Name)
					if strings.Compare(localTag, tag.Name) == 0 {
						fmt.Printf("found matching tag: %s", localTag)
						foundTag = true
						tagInfo = tag
					}
				}
				if foundTag {
					fmt.Println("check if newer image and if so pull")
					fmt.Println(tagInfo.ID)
					fmt.Println(getArchitecture())
					// go through images and find right arch
				}
			}

		} else {
			fmt.Printf("could not find %s on docker hub", imageName)
		}

		// login
		var authInfo = ""
		encodedJSON, err := json.Marshal(authInfo)
		if err != nil {
			panic(err)
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)

		// pull an image
		// imageName := "docker.io/library/alpine:latest"
		options := types.ImagePullOptions{}
		options.RegistryAuth = authStr
		out, err := cli.ImagePull(ctx, imageName, options)
		if err != nil {
			panic(err)
		}
		defer out.Close()

		io.Copy(os.Stdout, out)

	}

	// if _, err := ioutil.ReadAll(out); err != nil {
	// 	panic(err)
	// }

}
