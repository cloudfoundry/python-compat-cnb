package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/helper"
)

var (
	runtimePrefix = "python-"
)

func main() {
	context, err := detect.DefaultDetect()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create a default detection context: %s", err)
		os.Exit(100)
	}

	code, err := runDetect(context)
	if err != nil {
		context.Logger.Info(err.Error())
	}

	os.Exit(code)
}

func runDetect(context detect.Detect) (int, error) {
	runtimePath := filepath.Join(context.Application.Root, "runtime.txt")
	exists, err := helper.FileExists(runtimePath)
	if err != nil {
		return detect.FailStatusCode, err
	}

	var version string
	if exists {
		version, err = readRuntimeTxtVersion(runtimePath)
		if err != nil {
			return detect.FailStatusCode, err
		}
	}

	return context.Pass(buildplan.Plan{
		Requires: []buildplan.Required{
			{
				Name:     "python",
				Version:  version,
				Metadata: buildplan.Metadata{"launch": true},
			},
		},
	})
}

func readRuntimeTxtVersion(runtimePath string) (string, error) {
	buf, err := ioutil.ReadFile(runtimePath)
	if err != nil {
		return "", err
	}
	runtimeTxt := strings.TrimSpace(strings.TrimPrefix(string(buf), runtimePrefix))

	return runtimeTxt, nil
}
