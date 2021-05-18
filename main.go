package main

import (
	"context"
	"fmt"

	// "strings"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func main() {

	// log.PrintColors = true
	// log.PrintTimestamp = false

	// check if docker is running

	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// log.CheckError(err)
	images, err := cli.ImageList(context.Background(),
		types.ImageListOptions{
			All: false,
		})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, image := range images {
		for _, tag := range image.RepoTags {

			fmt.Println(tag)
			// tagParts := strings.SplitN(tag, ":", 2)

			// log.Infof(
			// 		"%-25s | %-15s | %s | %-30s | %10s",
			// 		tagParts[0],
			// 		tagParts[1],
			// 		image.ID[7:19],
			// 		time.Unix(image.Created, 0),
			// 		formatter.FileSize(image.Size),
			// )
		}
	}

}
