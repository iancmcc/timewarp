package timewarp_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTimewarp(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Timewarp Suite")
}

var _ = BeforeSuite(func() {
	rand.Seed(GinkgoRandomSeed())
})
