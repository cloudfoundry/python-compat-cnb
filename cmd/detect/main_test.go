package main

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/python-runtime-cnb/python"

	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestUnitDetect(t *testing.T) {
	spec.Run(t, "Detect", testDetect, spec.Report(report.Terminal{}))
}

func testDetect(t *testing.T, when spec.G, it spec.S) {
	var factory *test.DetectFactory

	it.Before(func() {
		RegisterTestingT(t)
		factory = test.NewDetectFactory(t)
	})

	when("there is a runtime.txt", func() {
		version := "1.2.3"
		runtimeTxt := fmt.Sprintf("python-%s\n", version)

		it.Before(func() {
			Expect(helper.WriteFile(filepath.Join(factory.Detect.Application.Root, "runtime.txt"), 0666, runtimeTxt)).To(Succeed())
		})

		it("should pass with the requested version of python", func() {
			code, err := runDetect(factory.Detect)
			Expect(err).NotTo(HaveOccurred())
			Expect(code).To(Equal(detect.PassStatusCode))

			Expect(factory.Plans.Plan).To(Equal(buildplan.Plan{
				Requires: []buildplan.Required{
					{
						Name:     python.Dependency,
						Version:  version,
						Metadata: buildplan.Metadata{"launch": true},
					},
				},
			}))
		})
	})
}
