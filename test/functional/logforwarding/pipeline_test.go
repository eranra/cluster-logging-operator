package fluent_test

import (
	"fmt"
	"github.com/ViaQ/logerr/log"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	logging "github.com/openshift/cluster-logging-operator/pkg/apis/logging/v1"
	"github.com/openshift/cluster-logging-operator/pkg/constants"
	"github.com/openshift/cluster-logging-operator/pkg/utils"
	"github.com/openshift/cluster-logging-operator/test/functional"
	"github.com/openshift/cluster-logging-operator/test/helpers/fluentd"
	"github.com/openshift/cluster-logging-operator/test/runtime"
	k8score "k8s.io/api/core/v1"
	"os"
	"time"
)

const message = "This is a test message"

// Extend client.Test for tests using logGenerator and receiver.
type FluentdFunctionalFramework struct {
	*functional.FluentdFunctionalFramework
	receiver     *fluentd.Receiver
	logGenerator *k8score.Pod
}

func SetEnvironmentVariables() error {
	// set fluentd image to default if not exist in env variable
	fluentdImage := utils.GetComponentImage(constants.FluentdName)
	if fluentdImage == "" {
		_ = os.Setenv("FLUENTD_IMAGE", "quay.io/openshift/origin-logging-fluentd:latest")
	}

	// set scripts_dir  to default if not exist in env variable
	scriptsDir := os.Getenv("SCRIPTS_DIR")
	if scriptsDir == "" {
		_ = os.Setenv("SCRIPTS_DIR", os.Getenv("GOPATH")+"/src/github.com/openshift/cluster-logging-operator/scripts")
	}

	// set logging_share_dir to default if not exist in env variable
	loggingShareDir := os.Getenv("LOGGING_SHARE_DIR")
	if loggingShareDir == "" {
		_ = os.Setenv("LOGGING_SHARE_DIR", os.Getenv("GOPATH")+"/src/github.com/openshift/cluster-logging-operator/files")
	}
	return nil
}
func NewFramework(message string) *FluentdFunctionalFramework {

	// set environment variables if not exist
	_ = SetEnvironmentVariables()

	fr := functional.NewFluentdFunctionalFramework()
	lg := runtime.NewLogGenerator(fr.Namespace, "log-generator", 10000, 0, message)
	lr := fluentd.NewReceiver(fr.Namespace, "receiver")
	functional.NewClusterLogForwarderBuilder(fr.Forwarder).
		FromInput(logging.InputNameApplication).
		ToFluentForwardOutput()
	framework := &FluentdFunctionalFramework{
		FluentdFunctionalFramework: fr,
		logGenerator:               lg,
		receiver:                   lr,
	}
	err := framework.Deploy()
	Expect(err).To(BeNil())

	return framework
}

func (framework *FluentdFunctionalFramework) Close() {
	framework.FluentdFunctionalFramework.Cleanup()
}

var _ = Describe("[ClusterLogForwarder] Function testing of fluentd forward pipeline formatting", func() {

	var (
		framework *FluentdFunctionalFramework
	)

	BeforeEach(func() {
		framework = NewFramework(message)
	})

	AfterEach(func() {
		framework.Close()
	})

	Context("with app/infra/audit receiver", func() {
		BeforeEach(func() {
			// Create logs receiver
			framework.receiver.AddSource(&fluentd.Source{Name: "application", Type: "forward", Port: 24224})
			err := framework.receiver.Create(framework.Test.Client)
			Expect(err).To(BeNil())

			// Create logs generator
			err = framework.Test.Client.Create(framework.logGenerator)
			Expect(err).To(BeNil())
		})

		It("forwards application logs only", func() {
			log.V(2).Info("Executing check metrics test")
			//cmd := fmt.Sprintf("curl -ksv https://%s.%s:24231/metrics", framework.Name, framework.Namespace)
			//metrics, _ := framework.RunCommand(cmd)
			//Expect(metrics).To(ContainSubstring("fluentd_output_status_buffer_total_bytes"))

			// Create containers directory
			cmd := fmt.Sprintf("mkdir /var/log/containers")
			_,err := framework.RunCommand(cmd)
			Expect(err).To(BeNil())

			// Create application log file with single line
			//cmd = fmt.Sprintf(`sh -c "echo \"test\" >> /var/log/containers/application.log"`)
			//_,err = framework.RunCommand(cmd)
			//Expect(err).To(BeNil())

			var logFilename= fmt.Sprintf("%s/%s_%s_%s-%s.log",
			"/var/log/containers",framework.Name,framework.Namespace,framework.Name,"0000")
			// Create application log file with single line
			var logLine = `2020-11-04T08:11:00.945011359+00:00 stderr F level=info ts=2020-11-04T08:11:00.944Z caller=main.go:330 msg=\\\"Starting Prometheus\\\" version=\\\"(version=2.15.2, branch=rhaos-4.5-rhel-7, revision=e1955f6306a7aa5b5be4509683913e431f95ad45)\\\"`
			cmd = fmt.Sprintf(`sh -c "echo \"%s\" >> %s"`, logLine,logFilename)
			_,err = framework.RunCommand(cmd)
			Expect(err).To(BeNil())

			// log source container log file content
			cmd = fmt.Sprintf(`cat %s`,logFilename)
			output1, err1 := framework.RunCommand(cmd)
			Expect(err1).To(BeNil())
			log.V(2).Info("The file "+logFilename+" looks like:\n", "output", output1)

			// log destination (after pipeline) file
			time.Sleep(100*5000)
			var outputFilename = "/opt/app-root/data/application"
			output2,err2 := runtime.Exec (framework.receiver.Pod, "cat", outputFilename).CombinedOutput()
			Expect(err2).To(BeNil())
			log.V(2).Info("The file "+outputFilename+" looks like:\n", "output", string(output2))
			log.V(2).Info("Done !")
		})
	})
})
