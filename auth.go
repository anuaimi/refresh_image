package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
