package integration_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/dagger"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	. "github.com/onsi/gomega"
)

var (
	bpDir, pythonCompatURI, pythonURI, pipURI string
)

func TestIntegration(t *testing.T) {
	var err error
	Expect := NewWithT(t).Expect
	bpDir, err = dagger.FindBPRoot()
	Expect(err).NotTo(HaveOccurred())
	pythonCompatURI, err = dagger.PackageBuildpack(bpDir)
	Expect(err).ToNot(HaveOccurred())
	defer dagger.DeleteBuildpack(pythonCompatURI)

	pythonURI, err = dagger.GetLatestBuildpack("python-runtime-cnb")
	Expect(err).ToNot(HaveOccurred())
	defer dagger.DeleteBuildpack(pythonURI)

	pipURI, err = dagger.GetLatestBuildpack("pip-cnb")
	Expect(err).ToNot(HaveOccurred())
	defer dagger.DeleteBuildpack(pipURI)

	spec.Run(t, "Integration", testIntegration, spec.Report(report.Terminal{}))
}

func testIntegration(t *testing.T, when spec.G, it spec.S) {
	var Expect func(interface{}, ...interface{}) GomegaAssertion
	it.Before(func() {
		Expect = NewWithT(t).Expect
	})

	when("building a simple app", func() {
		it("using the python version in the runtime.txt", func() {
			app, err := dagger.PackBuild(filepath.Join("testdata", "simple_app"), pythonURI, pythonCompatURI, pipURI)
			Expect(err).ToNot(HaveOccurred())

			app.SetHealthCheck("", "3s", "1s")

			err = app.Start()
			if err != nil {
				_, err = fmt.Fprintf(os.Stderr, "App failed to start: %v\n", err)
				containerID, imageName, volumeIDs, err := app.Info()
				Expect(err).NotTo(HaveOccurred())
				fmt.Printf("ContainerID: %s\nImage Name: %s\nAll leftover cached volumes: %v\n", containerID, imageName, volumeIDs)

				containerLogs, err := app.Logs()
				Expect(err).NotTo(HaveOccurred())
				fmt.Printf("Container Logs:\n %s\n", containerLogs)
				t.FailNow()
			}

			body, _, err := app.HTTPGet("/version")
			Expect(err).ToNot(HaveOccurred())
			Expect(body).To(ContainSubstring("2.7."))
		})
	})

}
