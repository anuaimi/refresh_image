# refresh_image

one of the challenges with docker images is keeping track of updates. The Docker CLI has not command to list the tags for an image on a remote repository like Docker Hub. So you end up having to go to the Docker Hub website...for every image you are using.

`refresh_image` will list both the versions (tags) you have locally and the most recent tags available on Docker Hub.  

```bash
$ ./refresh_images postgres
Checking for postgres :
  Local:
    13.2-alpine               e43172f9204f  created: 04-14-2021 23:52
    13.2                      1f0815c1cb6e  created: 02-12-2021 14:21
    13.1-alpine               8c6053d81a45  created: 01-28-2021 22:01
  Docker Hub:
    latest                    4bca01db9119  created: 05-14-2021 20:32
    11.12                     f06f1a9b6a82  created: 05-14-2021 20:30
    11                        f06f1a9b6a82  created: 05-14-2021 20:30
    10.17                     cb8328596b02  created: 05-14-2021 20:29
    10                        cb8328596b02  created: 05-14-2021 20:29
    alpine                    b2db3c347304  created: 05-14-2021 20:33
    9.6.22-alpine             0fceaf685127  created: 05-14-2021 20:34
    9.6.22                    b67f157ec6a1  created: 05-14-2021 20:33
    .
    .
    .
    9.6-alpine                0fceaf685127  created: 05-14-2021 20:34
```

## Background

Originally I wanted to create a simple 'update' and 'upgrade' tool for the images I had locally. But this is challenging as many of the common docker images have a large number of tags and 'streams' (ie slim, stretch, alpine, etc).  Node is a great example of this. A possible future enhancement might be to add some repo specific knowledge so that the tool knows all of the various streams

