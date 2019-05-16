package main_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	. "./"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


var _ = Describe("Routes", func() {
	var db Db = mockDb{}
	ctx := context.WithValue(testCtx, ContextDb, db)
	gin.SetMode(gin.TestMode)

	Describe("V1", func() {
		Describe("GET /", func() {
			It("should succeed with OK", func() {
				w := performRequest(ctx, "GET", "/")
				Expect(w.Code).To(Equal(http.StatusOK))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(string(r)).To(Equal("OK"))
			})
			It("should fail with invalid method", func() {
				w := performRequest(ctx, "POST", "/")
				Expect(w.Code).To(Equal(http.StatusMethodNotAllowed))
			})
		})
		Describe("GET /v1/payments/", func() {
			It("should list payments", func() {
				w := performRequest(ctx, "GET", "/v1/payments/")
				Expect(w.Code).To(Equal(http.StatusOK))
				Expect(w.Header().Get("content-type")).To(ContainSubstring("application/json"))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(r).To(MatchJSON(`{"data": [], "links": {"self": "http://example.com/v1/payments/?count=10", "next": null}}`))
			})
			It("should default zero count to 10", func() {
				w := performRequest(ctx, "GET", "/v1/payments/?count=0")
				Expect(w.Code).To(Equal(http.StatusOK))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(r).To(MatchJSON(`{"data": [], "links": {"self": "http://example.com/v1/payments/?count=10", "next": null}}`))
			})
			It("should default negative count size to 10", func() {
				w := performRequest(ctx, "GET", "/v1/payments/?count=-1")
				Expect(w.Code).To(Equal(http.StatusOK))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(r).To(MatchJSON(`{"data": [], "links": {"self": "http://example.com/v1/payments/?count=10", "next": null}}`))
			})
			It("should list payments with updated count", func() {
				w := performRequest(ctx, "GET", "/v1/payments/?count=20")
				Expect(w.Code).To(Equal(http.StatusOK))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(r).To(MatchJSON(`{"data": [], "links": {"self": "http://example.com/v1/payments/?count=20", "next": null}}`))
			})
			It("should list payments with updated after", func() {
				w := performRequest(ctx, "GET", "/v1/payments/?count=20&after=5cdd382e9549af35c3b94301")
				Expect(w.Code).To(Equal(http.StatusOK))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(r).To(MatchJSON(`{"data": [], "links": {"self": "http://example.com/v1/payments/?count=20&after=5cdd382e9549af35c3b94301", "next": null}}`))
			})
			It("should list payments from database", func() {
				c := context.WithValue(ctx, ContextDb, mockDb{
					Payments: []Payment{
						{ID: *id},
						}})
				w := performRequest(c, "GET", "/v1/payments/?count=10")
				Expect(w.Code).To(Equal(http.StatusOK))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(r).To(MatchJSON(`{
					"data": [
						{"id": "5cdd382e9549af35c3b94301", "links": {"self": "http://example.com/v1/payments/5cdd382e9549af35c3b94301/"}}], 
					"links": {"self": "http://example.com/v1/payments/?count=10", "next": null}}`))
			})
			It("should have next link if there are too many items", func() {
				c := context.WithValue(ctx, ContextDb, mockDb{
					Payments: []Payment{
						{ID: *id},
						{ID: *id2},
					}})
				w := performRequest(c, "GET", "/v1/payments/?count=1")
				Expect(w.Code).To(Equal(http.StatusOK))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(r).To(MatchJSON(`{
					"data": [
						{"id": "5cdd382e9549af35c3b94301", "links": {"self": "http://example.com/v1/payments/5cdd382e9549af35c3b94301/"}}], 
					"links": {"self": "http://example.com/v1/payments/?count=1", "next": "http://example.com/v1/payments/?count=1&after=5cdd382e9549af35c3b94301"}}`))
			})
		})

		Describe("GET /v1/payments/{id}", func() {
			It("should return 404 on not found payment", func() {
				w := performRequest(ctx, "GET", "/v1/payments/non-existing-id")
				Expect(w.Code).To(Equal(http.StatusNotFound))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(string(r)).To(ContainSubstring(http.StatusText(http.StatusNotFound)))
			})
			It("should payment resource when in database", func() {
				c := context.WithValue(ctx, ContextDb, mockDb{
					Payments: []Payment{
						paymentSample,
					}})
				w := performRequest(c, "GET", "/v1/payments/5cdd382e9549af35c3b94301")
				Expect(w.Code).To(Equal(http.StatusOK))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(string(r)).To(MatchJSON(`
				{
					"type": "Payment",
				  	"id": "5cdd382e9549af35c3b94301",
				  	"version": 0,
				  	"organisation_id": "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
					"attributes": {
						"amount": "100.21",
						"beneficiary_party": {
						  "account_name": "W Owens",
						  "account_number": "31926819",
						  "account_number_code": "BBAN",
						  "account_type": 0,
						  "address": "1 The Beneficiary Localtown SE2",
						  "bank_id": "403000",
						  "bank_id_code": "GBDSC",
						  "name": "Wilfred Jeremiah Owens"
						},
						"charges_information": {
						  "bearer_code": "SHAR",
						  "sender_charges": [
							{
							  "amount": "5.00",
							  "currency": "GBP"
							},
							{
							  "amount": "10.00",
							  "currency": "USD"
							}
						  ],
						  "receiver_charges_amount": "1.00",
						  "receiver_charges_currency": "USD"
						},
						"currency": "GBP",
						"debtor_party": {
						  "account_name": "EJ Brown Black",
						  "account_number": "GB29XABC10161234567801",
						  "account_number_code": "IBAN",
						  "account_type": 0,
						  "address": "10 Debtor Crescent Sourcetown NE1",
						  "bank_id": "203301",
						  "bank_id_code": "GBDSC",
						  "name": "Emelia Jane Brown"
						},
						"end_to_end_reference": "Wil piano Jan",
						"fx": {
						  "contract_reference": "FX123",
						  "exchange_rate": "2.00000",
						  "original_amount": "200.42",
						  "original_currency": "USD"
						},
						"numeric_reference": "1002001",
						"payment_id": "123456789012345678",
						"payment_purpose": "Paying for goods/services",
						"payment_scheme": "FPS",
						"payment_type": "Credit",
						"processing_date": "2017-01-18",
						"reference": "Payment for Em's piano lessons",
						"scheme_payment_sub_type": "InternetBanking",
						"scheme_payment_type": "ImmediatePayment",
						"sponsor_party": {
						  "account_number": "56781234",
						  "bank_id": "123123",
						  "bank_id_code": "GBDSC"
						}
					},
					"links": {"self": "http://example.com/v1/payments/5cdd382e9549af35c3b94301/"}
				}`))
			})
		})

		Describe("PUT /v1/payments/{id}", func() {
			It("should return 404 on not found payment", func() {
				w := performRequest(ctx, "PUT", "/v1/payments/non-existing-id")
				Expect(w.Code).To(Equal(http.StatusNotFound))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(string(r)).To(ContainSubstring(http.StatusText(http.StatusNotFound)))
			})
			It("should return 400 on invalid body", func() {
				c := context.WithValue(ctx, ContextDb, mockDb{
					Payments: []Payment{
						paymentSample,
					}})
				w := performRequestBody(c, "PUT", "/v1/payments/5cdd382e9549af35c3b94301", strings.NewReader("{"))
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
			It("should return 204 on success", func() {
				c := context.WithValue(ctx, ContextDb, mockDb{
					Payments: []Payment{
						paymentSample,
					}})
				w := performRequestBody(c, "PUT", "/v1/payments/5cdd382e9549af35c3b94301", strings.NewReader("{}"))
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
			It("should return 500 on internal server error", func() {
				c := context.WithValue(ctx, ContextDb, mockDb{
					error: errors.New("noooo"),
					Payments: []Payment{
						paymentSample,
					}})
				w := performRequestBody(c, "PUT", "/v1/payments/5cdd382e9549af35c3b94301", strings.NewReader("{}"))
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(string(r)).ToNot(ContainSubstring("noooo"))

			})
		})

		Describe("DELETE /v1/payments/{id}", func() {
			It("should return 404 on not found payment", func() {
				w := performRequest(ctx, "DELETE", "/v1/payments/non-existing-id")
				Expect(w.Code).To(Equal(http.StatusNotFound))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(string(r)).To(ContainSubstring(http.StatusText(http.StatusNotFound)))
			})
			It("should return 204 on success", func() {
				c := context.WithValue(ctx, ContextDb, mockDb{
					Payments: []Payment{
						paymentSample,
					}})
				w := performRequestBody(c, "DELETE", "/v1/payments/5cdd382e9549af35c3b94301", strings.NewReader("{}"))
				Expect(w.Code).To(Equal(http.StatusNoContent))
			})
			It("should return 500 on internal server error", func() {
				c := context.WithValue(ctx, ContextDb, mockDb{
					error: errors.New("noooo"),
					Payments: []Payment{
						paymentSample,
					}})
				w := performRequestBody(c, "DELETE", "/v1/payments/5cdd382e9549af35c3b94301", strings.NewReader("{}"))
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(string(r)).ToNot(ContainSubstring("noooo"))

			})
		})
		Describe("POST /v1/payments/", func() {
			It("should return 400 on invalid body", func() {
				c := context.WithValue(ctx, ContextDb, mockDb{
					Payments: []Payment{
						paymentSample,
					}})
				w := performRequestBody(c, "POST", "/v1/payments", strings.NewReader("{"))
				Expect(w.Code).To(Equal(http.StatusBadRequest))
			})
			It("should return 200 on success", func() {
				c := context.WithValue(ctx, ContextDb, mockDb{
					Payments: []Payment{
						paymentSample,
					}})
				w := performRequestBody(c, "POST", "/v1/payments", strings.NewReader("{}"))
				Expect(w.Code).To(Equal(http.StatusCreated))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(string(r)).To(MatchJSON(`
				{ "id": "5cdd382e9549af35c3b94301", "links": {"self": "http://example.com/v1/payments/5cdd382e9549af35c3b94301/"} }
				`))
			})
			It("should return 500 on internal server error", func() {
				c := context.WithValue(ctx, ContextDb, mockDb{
					error: errors.New("noooo"),
					Payments: []Payment{
						paymentSample,
					}})
				w := performRequestBody(c, "POST", "/v1/payments", strings.NewReader("{}"))
				Expect(w.Code).To(Equal(http.StatusInternalServerError))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(string(r)).ToNot(ContainSubstring("noooo"))

			})
		})

	})
})
