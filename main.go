package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// DockerAuthResponse has details of response from docker registry auth request
type DockerAuthResponse struct {
	Token       string `json:"token"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	IssuedAt    string `json:"issued_at"`
}

// getAuthToken will request a token from the docker registry
func getAuthToken() (authToken string, err error) {

	// request an auth token
	resp, err := http.Get("https://auth.docker.io/token?service=registry.docker.io")
	if err != nil {
		return "", err
	}

	// process response
	var authResults DockerAuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResults)
	if err != nil {
		return "", err
	}

	return authResults.Token, err
}

// validateLogin will check a username and password with docker hub
func validateLogin(cli *client.Client, username string, password string) (bool, error) {

	// see if we have the right credentials for docker hub
	authInfo := types.AuthConfig{}
	authInfo.Username = username
	authInfo.Password = password
	// loginInfo, err := cli.RegistryLogin(ctx, authInfo)
	ctx := context.Background()
	_, err := cli.RegistryLogin(ctx, authInfo)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return true, err
}

// checkDockerHubForImage will try and find if image came from docker hub
func checkDockerHubForImage(cli *client.Client, imageName string) bool {

	// query docker hub
	ctx := context.Background()
	searchOptions := types.ImageSearchOptions{}
	searchOptions.Limit = 100
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
			fmt.Printf("found match for %s\n", imageName)
			return true
		}
	}

	// nope, image not from docker hub
	return false
}

// getTags will query docker registry to get tags available for an image
func getTags(imageName string) ([]string, error) {

	var tags []string

	// if no / add library to front
	if strings.ContainsAny(imageName, "/") == false {
		imageName = "library/" + imageName
	}
	url := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags?page_size=100", imageName)
	resp, err := http.Get(url)
	if err != nil {
		return tags, err
	}

	err = json.NewDecoder(resp.Body).Decode(&tags)
	if err != nil {
		return tags, err
	}

	return tags, err
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
				fmt.Println(err)
				continue
			}
			fmt.Println(tags)

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
