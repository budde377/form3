// +build integration

package main_test

import (
	"encoding/json"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"strings"
)

var _ = Describe("RoutesIntegration", func() {
	BeforeEach(func() {
		_ = db.(TestDb).Drop(testCtx)
	})
	It("should create and read a payment", func() {
		// Get list of payments
		w := performRequest(testDbCtx, "GET", "/v1/payments")
		Expect(w.Code).To(Equal(http.StatusOK))
		r, _ := ioutil.ReadAll(w.Body)
		Expect(string(r)).To(MatchJSON(`
		{
			"data": [],
			"links": {
				"self": "http://example.com/v1/payments/?count=10",
			  	"next": null
			}
		}

		`))

		// Create a new payment
		w = performRequestBody(testDbCtx, "POST", "/v1/payments", strings.NewReader(fmt.Sprintf(`
		 { "attributes": %s }
		`, paymentSampleAttributesJSON)))

		Expect(w.Code).To(Equal(http.StatusCreated))
		dec := json.NewDecoder(w.Body)
		var res struct{ ID string `json:"id"` }
		_ = dec.Decode(&res)
		Expect(res.ID).To(MatchRegexp("[0-9a-z]{24}"))

		// List payments again
		w = performRequest(testDbCtx, "GET", "/v1/payments")
		Expect(w.Code).To(Equal(http.StatusOK))
		r, _ = ioutil.ReadAll(w.Body)
		Expect(r).To(MatchJSON(fmt.Sprintf(`
		{
			"data": [{"id": "%s", "links": {"self": "http://example.com/v1/payments/%s/"}}],
			"links": {
				"self": "http://example.com/v1/payments/?count=10",
			  	"next": null
			}
		}`, res.ID, res.ID)))


		// Fetch payment
		w = performRequest(testDbCtx, "GET", fmt.Sprintf("/v1/payments/%s", res.ID))
		Expect(w.Code).To(Equal(http.StatusOK))
		r, _ = ioutil.ReadAll(w.Body)
		Expect(r).To(MatchJSON(fmt.Sprintf(`
		{
			"id": "%s",
			"attributes": %s,
			"version": 0,
			"organisation_id": "",
			"links": {
				"self": "http://example.com/v1/payments/%s/"
			},
			"type": "Payment"
		}`, res.ID, paymentSampleAttributesJSON, res.ID)))


		// Update payment
		w = performRequestBody(testDbCtx, "PUT", fmt.Sprintf("/v1/payments/%s",res.ID), strings.NewReader(fmt.Sprintf(`
		 	{ "attributes": %s, "organisation_id": "org1"}
		`, paymentSampleAttributesJSON)))
		Expect(w.Code).To(Equal(http.StatusNoContent))

		// Read again
		w = performRequest(testDbCtx, "GET", fmt.Sprintf("/v1/payments/%s", res.ID))
		Expect(w.Code).To(Equal(http.StatusOK))
		r, _ = ioutil.ReadAll(w.Body)
		Expect(r).To(MatchJSON(fmt.Sprintf(`
		{
			"id": "%s",
			"attributes": %s,
			"version": 1,
			"organisation_id": "org1",
			"links": {
				"self": "http://example.com/v1/payments/%s/"
			},
			"type": "Payment"
		}`, res.ID, paymentSampleAttributesJSON, res.ID)))


		// Delete payment
		w = performRequestBody(testDbCtx, "DELETE", fmt.Sprintf("/v1/payments/%s",res.ID), strings.NewReader(fmt.Sprintf(`
		 	{ "attributes": %s, "organisation_id": "org1"}
		`, paymentSampleAttributesJSON)))
		Expect(w.Code).To(Equal(http.StatusNoContent))

		// Read again
		w = performRequest(testDbCtx, "GET", fmt.Sprintf("/v1/payments/%s", res.ID))
		Expect(w.Code).To(Equal(http.StatusNotFound))

		// List
		w = performRequest(testDbCtx, "GET", "/v1/payments")
		Expect(w.Code).To(Equal(http.StatusOK))
		r, _ = ioutil.ReadAll(w.Body)
		Expect(string(r)).To(MatchJSON(`
		{
			"data": [],
			"links": {
				"self": "http://example.com/v1/payments/?count=10",
			  	"next": null
			}
		}

		`))

	})
})
