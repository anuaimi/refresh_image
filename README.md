# refresh_image

one of the challenges with local docker images is keeping track of updates. The Docker CLI does not have a command to list the newer tags for an image on a remote repository (like Docker Hub). So to see if there is a newer image, you end up having to go to the Docker Hub website...for every image you are using.

When you call `refresh_image`, you give it the image you want tag infor for.  It will list both all the local tags for the image and all the most recent tags on Docker Hub.  It will exclude from the Docker Hub results: a) any tag that is older than the oldest local tag and b) if any of the local tags are versions, all older versions on Docker Hub.

Below you can see the results for the `alpine` image.  You can see that there are two local tags and that Docker Hub has a few newer ones.  You can now easily issue a `docker pull` with whatever new tag you want to download

```bash
$ ./refresh_image alpine
  Local Images:
    latest                    6dbb9cc54074  created: 04-14-2021 15:19
    3.13                      28f6e2705743  created: 02-17-2021 16:19
  Docker Hub:
    latest                    def822f9851c  created: 04-14-2021 19:37
    3.13.5                    def822f9851c  created: 04-14-2021 19:37
    3.13                      def822f9851c  created: 04-14-2021 19:37
    3.13.4                    e103c1b4bf01  created: 03-31-2021 20:13
    3.13.3                    4266485e304a  created: 03-25-2021 22:24
    edge                      45fbb9ea28b1  created: 03-25-2021 22:24
    20210212                  45fbb9ea28b1  created: 03-25-2021 22:24
    3.13.2                    4661fb57f789  created: 02-17-2021 21:33
```

## Background

Originally I wanted to create a simple 'update' and 'upgrade' tool for the images I had locally. But this is challenging as many of the common docker images have a large number of tags and 'streams' (ie slim, stretch, alpine, etc).  Node is a great example of this. A possible future enhancement might be to add some repo specific knowledge so that the tool knows all of the various streams
