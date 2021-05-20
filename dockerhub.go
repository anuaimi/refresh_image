package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerImageInfo struct {
	Architecture string `json:"architecture"`
	Features     string `json:"features"`
	Variant      string `json:"variant"`
	Digest       string `json:"digest"`
	OS           string `json:"os"`
	OSFeatures   string `json:"os_features"`
	OSVersion    string `json:"os_version"`
	Size         int64  `json:"size"`
	Status       string `json:"status"`
	LastPulled   string `json:"last_pulled"`
	LastPushed   string `json:"last_pushed"`
}

type DockerImageTag struct {
	Creator             int               `json:"creator"`
	ID                  int               `json:"id"`
	Images              []DockerImageInfo `json:"images"`
	LastUpdated         string            `json:"last_updated"`
	LastUpdater         int               `json:"last_updater"`
	LastUpdaterUsername string            `json:"last_updater_username"`
	Name                string            `json:"name"`
	Repository          int               `json:"repository"`
	FullSize            int               `json:"full_size"`
	V2                  bool              `json:"v2"`
	TagStatus           string            `json:"tag_status"`
	TagLastPulled       string            `json:"tag_last_pulled"`
	TagLastPushed       string            `json:"tag_last_pushed"`
}

type DockerTagQueryResults struct {
	Count   int              `json:"count"`
	Next    string           `json:"next"`
	Results []DockerImageTag `json:"results"`
}

// checkDockerHubForImage will try and find if image came from docker hub
func checkDockerHubForImage(cli *client.Client, imageName string) bool {

	// query docker hub
	ctx := context.Background()
	searchOptions := types.ImageSearchOptions{}
	searchOptions.Limit = 50
	// searchOptions.Limit = 100
	results, err := cli.ImageSearch(ctx, imageName, searchOptions)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// get lots of partial matches, check each to see if exact match
	for _, result := range results {

		// look for exact match
		if strings.Compare(result.Name, imageName) == 0 {
			// need to get tags!!! and decide if to a pull
			// fmt.Printf("found match for %s\n", imageName)
			return true
		}
	}

	// nope, image not from docker hub
	return false
}

// getTags will query docker registry to get 1st 100 tags available for an image
func getTagsForImage(imageName string) ([]DockerImageTag, error) {

	var queryResults DockerTagQueryResults

	// if no / add library to front
	if strings.ContainsAny(imageName, "/") == false {
		imageName = "library/" + imageName
	}
	url := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags?page_size=50", imageName)
	resp, err := http.Get(url)
	if err != nil {
		return queryResults.Results, err
	}

	err = json.NewDecoder(resp.Body).Decode(&queryResults)
	if err != nil {
		return queryResults.Results, err
	}

	// if there are more than 100 tags, need to loop to get rest
	// if queryResults.Count > 100, then use queryResults.next as URL to get next 100
	// can use queryResults.next as is (it has page_size param)

	return queryResults.Results, err
}

// getImageForArchitecture will find correct image for given architecture
func getImageForArchitecture(images []DockerImageInfo, arch string) (DockerImageInfo, error) {
	for _, image := range images {
		if image.Architecture == arch {
			return image, nil
		}
	}
	return DockerImageInfo{}, errors.New("could not find image for given architecture")
}

// pullImage will download the desired docker image
func pullImage(cli *client.Client, image string, desiredTag string) error {

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
	ctx := context.Background()
	dockerImage := image + ":" + desiredTag
	out, err := cli.ImagePull(ctx, dockerImage, options)
	if err != nil {
		return err
	}
	defer out.Close()

	io.Copy(os.Stdout, out)

	// if _, err := ioutil.ReadAll(out); err != nil {
	// 	panic(err)
	// }

	return nil
}
