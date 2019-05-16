package main_test

import (
	"github.com/google/logger"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestForm3(t *testing.T) {
	defer logger.Init("Form3 API", false, false, ioutil.Discard).Close()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Form3 Suite")
}
