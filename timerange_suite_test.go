package timewarp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTimerange(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Timerange Suite")
}
