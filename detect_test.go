package pythonstart_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/paketo-buildpacks/packit"
	pythonstart "github.com/paketo-buildpacks/python-start"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testDetect(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir string
		detect     packit.DetectFunc
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		err = os.WriteFile(filepath.Join(workingDir, "x.py"), []byte{}, os.ModePerm)
		Expect(err).NotTo(HaveOccurred())

		detect = pythonstart.Detect()
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("detection phase", func() {
		it("detects", func() {
			result, err := detect(packit.DetectContext{
				WorkingDir: workingDir,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Plan).To(Equal(packit.BuildPlan{
				Provides: []packit.BuildPlanProvision{},
				Requires: []packit.BuildPlanRequirement{
					{
						Name: "cpython",
						Metadata: pythonstart.BuildPlanMetadata{
							Launch: true,
						},
					},
				},
				Or: []packit.BuildPlan{
					{
						Provides: []packit.BuildPlanProvision{},
						Requires: []packit.BuildPlanRequirement{
							{
								Name: "cpython",
								Metadata: pythonstart.BuildPlanMetadata{
									Launch: true,
								},
							},
							{
								Name: "site-packages",
								Metadata: pythonstart.BuildPlanMetadata{
									Launch: true,
								},
							},
						},
					},
					{
						Provides: []packit.BuildPlanProvision{},
						Requires: []packit.BuildPlanRequirement{
							{
								Name: "conda-environment",
								Metadata: pythonstart.BuildPlanMetadata{
									Launch: true,
								},
							},
						},
					},
				},
			}))
		})

		context("When no python related files are present", func() {
			it.Before(func() {
				Expect(os.RemoveAll(filepath.Join(workingDir, "x.py"))).To(Succeed())
			})

			it("fails detection", func() {
				_, err := detect(packit.DetectContext{
					WorkingDir: workingDir,
				})
				Expect(err).To(MatchError(ContainSubstring("No *.py, environment.yml or package-list.txt found")))
			})
		})
	})
}
