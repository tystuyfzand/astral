package arikawa_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestArikawa(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Arikawa Suite")
}
