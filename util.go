package main

import "runtime"

// getArchitecture will return the hardware running on
func getMachineArchitecture() string {
	// amd64, arm64 are main values
	return runtime.GOARCH
}

func findTagMatch(exact bool) string {
	// either match string exactly or
	// extract version number
	//   version number same or larger & has same suffix

	return ""
}

// findVersionOfLatest will look on Docker Hub to see if can find version of latest
func findTagOfLatest(digest string) string {
	// look for digest match
	// return tag of match (that is not latest)
	// note, there may be several ( 1, 1.4, 1.4.2)
	return ""
}

// breakdowntag will try and look for a version number
func breakdownTag(tag string) (string, string) {

	// extract:
	//   1.15.3-nanoserver
	//   1.15.3-nanoserver-1809
	//   1
	//   1-buster
	//   1.15-buster
	//   v4.2.0
	//   1.0.0

	// don't split
	//   latest
	//   slim-buster
	//   nanoserver-1809
	//   alpine3.13
	//   trusty-20191215
	//   latest-linux-amd64

	// cases:
	// 1. tag is only numbers and dots
	// 2. starts with numbers and dots but has alpha text after that
	// 3. everything else (don't break apart)

	return "", ""
}
