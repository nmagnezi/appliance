package coreos

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2/dsl/core"
	. "github.com/onsi/gomega"
	"github.com/openshift/appliance/pkg/asset/config"
	"github.com/openshift/appliance/pkg/executer"
	"github.com/openshift/appliance/pkg/release"
)

var _ = Describe("Test CoreOS", func() {
	var (
		ctrl         *gomock.Controller
		mockExecuter *executer.MockExecuter
		mockRelease  *release.MockRelease
		testCoreOs   CoreOS
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockExecuter = executer.NewMockExecuter(ctrl)
		mockRelease = release.NewMockRelease(ctrl)
		coreOSConfig := CoreOSConfig{
			Release:    mockRelease,
			Executer:   mockExecuter,
			EnvConfig:  &config.EnvConfig{},
			CpuArch:    config.CpuArchitectureX86,
			PullSecret: "'{\"auths\":{\"\":{\"auth\":\"dXNlcjpwYXNz\"}}}'",
		}
		testCoreOs = NewCoreOS(coreOSConfig)
	})

	It("DownloadISO - success", func() {
		mockRelease.EXPECT().ExtractFile(machineOsImageName, fmt.Sprintf(coreOsFileName, config.CpuArchitectureX86)).Return("/path/to/file", nil).Times(1)
		_, err := testCoreOs.DownloadISO()
		Expect(err).ToNot(HaveOccurred())
	})

	It("DownloadISO - fail", func() {
		mockRelease.EXPECT().ExtractFile(machineOsImageName, fmt.Sprintf(coreOsFileName, config.CpuArchitectureX86)).Return("", errors.New("some error")).Times(1)
		_, err := testCoreOs.DownloadISO()
		Expect(err).To(HaveOccurred())
	})

	It("FetchCoreOSStream - fail", func() {
		mockRelease.EXPECT().ExtractFile(machineOsImageName, coreOsStream).Return("", errors.New("some error")).Times(1)
		_, err := testCoreOs.FetchCoreOSStream()
		Expect(err).To(HaveOccurred())
	})

})

func TestCoreOS(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "coreos_test")
}
