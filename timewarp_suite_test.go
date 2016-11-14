package timewarp_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTimewarp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Timewarp Suite")
}
