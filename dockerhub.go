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
	"time"

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

type DockerImageTags []DockerImageTag

type Repository struct {
	name string
	tags []DockerImageTag
	cli  *client.Client
}

// Init will open connection to local docker daemon
func (r *Repository) Init() error {

	// connect to local docker
	cli, err := client.NewEnvClient()
	if err != nil {
		// does not seem to be running
		return err
	}
	r.cli = cli

	return nil
}

// Find the given repo on docker hub
func (r *Repository) Find(name string) error {

	r.name = name

	// search
	ctx := context.Background()
	searchOptions := types.ImageSearchOptions{}
	searchOptions.Limit = 100
	matches, err := r.cli.ImageSearch(ctx, r.name, searchOptions)
	if err != nil {
		return err
	}

	// get lots of partial matches, check each to see if exact match
	for _, match := range matches {
		if strings.Compare(match.Name, r.name) == 0 {
			return nil
		}
	}

	// could not find desired image
	return errors.New("repo does not exist on Docker Hub")
}

// GetTags for a given repository
func (r *Repository) GetTags() error {

	var queryResults DockerTagQueryResults

	// if no '/' add library to front
	var imageName string
	if strings.ContainsAny(r.name, "/") == false {
		imageName = "library/" + r.name
	} else {
		imageName = r.name
	}

	url := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags?page_size=100", imageName)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	// decode json into golang struct
	err = json.NewDecoder(resp.Body).Decode(&queryResults)
	if err != nil {
		return err
	}

	if queryResults.Count > 100 {
		// need to get rest of tags

		// make channel to get results

		// loop from 100 to count (mod 100)
		//   can use queryResults.next as is (it has page_size param)
		//   call go routine to get give page
		//   go routine passes back using channel

		// we wait until go routines finish or timeout
		// merge all the results back together
	}

	r.tags = queryResults.Results

	return nil

}

// ListTags will print all the tags in the repository
func (r *Repository) ListTags() {

	arch := getMachineArchitecture()

	fmt.Println("  Docker Hub:")

	// for each tag
	for _, tag := range r.tags {
		// get correct image for our type of machine
		image, err := getImageForArchitecture(tag.Images, arch)
		if err == nil {
			imageTimestamp, err := time.Parse(time.RFC3339, image.LastPushed)
			if err == nil {
				fmt.Printf("    %-25s %s  created: %s\n", tag.Name, image.Digest[7:19], imageTimestamp.Format("01-02-2006 15:04"))
			}
		}
	}
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
