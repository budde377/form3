// +build integration

package main_test

import (
	. "./"
	. "github.com/onsi/ginkgo"
)

var c = OpenConfigFile("./test/config.json")

var _ = Describe("RoutesIntegration", func() {

	var db = InitDb(c)

})
