# refresh_images

one of the challenges with docker images is keeping track of updates. While the Docker CLI can search for images, there is no easy way to check for new tags for that image. So you end up going to the Docker Hub website, for every single image you are using.

The tool will list what versions (tags) of an image you currently have and then will also list the most recent tags available on Docker Hub.  

## Background

Originally I wanted to create a simple 'update' and 'upgrade' tool.  For popular images, there are check for new versions of local docker images.  But this is challenging has many of the common images have a large number of tags and 'streams'.  Node is a great example of this.  It might be possible to add some repo specific knowledge so that the tool knows all of the streams and can track the current update path
