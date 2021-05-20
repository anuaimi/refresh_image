package main

import "runtime"

// getArchitecture will return the hardware running on
func getMachineArchitecture() string {
	// amd64, arm64 are main values
	return runtime.GOARCH
}
