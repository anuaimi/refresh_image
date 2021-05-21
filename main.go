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

	docker := LocalDocker{}
	err := docker.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = docker.Find(dockerImage)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	docker.ListTags()

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
	repository.ListTags()

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
