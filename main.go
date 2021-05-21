package main

import (
	"fmt"
	"os"
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

	// see what tags we have locally for given image
	docker := LocalDocker{}
	err := docker.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = docker.Find(dockerImage)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	// show user list of tags that are local
	docker.ListTags()

	// get filters for our search on Docker hub
	oldestTimestamp, err := docker.GetOldestImage()
	if err != nil {
		// really shouldn't ever happen as long as we have one image
		fmt.Println(err)
		os.Exit(3)
	}
	oldestVersion, err := docker.GetMinVersion()
	if err != nil {
		// will happen when image tags aren't in semver format
	}
	// localTags := docker.GetAllTags()

	// now see what is on Docker Hub

	repository := Repository{}
	err = repository.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// check Docker Hub for repo we desire
	err = repository.Find(dockerImage)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = repository.GetTags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// display the tags
	repository.ListFilteredTags(oldestTimestamp, *oldestVersion)

	// go through local tags
	// for _, tag := range allImageTags {

	// 	// if 'latest' see if digest & date the same
	// 	// look for exact match - see if digest the same

	// 	// just the tag no digest
	// 	tag, err := repository.FindTag(tag)
	// 	if err != nil {
	// 		// is it newer?
	// 	}

	// }

}
