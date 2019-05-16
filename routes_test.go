package main_test

import (
	. "./"
	"context"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

var config = OpenConfigFile("./test/config.json")

type mockDb struct {
	Payments []Payment
}

func (d mockDb) GetPayments(ctx context.Context, size int, after *string) (*[]PaymentSummary, error) {
	summaries := make([]PaymentSummary, 0)
	for _, v := range d.Payments {
		summaries = append(summaries, PaymentSummary{
			Id: v.Id,
		})
	}
	return &summaries, nil
}

func (d mockDb) GetPaymentById(ctx context.Context, id string) (*Payment, error) {
	for _, v := range d.Payments {
		if v.Id == id {
			return &v, nil
		}
	}
	return nil, nil
}

func performRequest(ctx context.Context, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	RootRoute().ServeHTTP(w, req.WithContext(ctx))
	return w
}

var _ = Describe("Routes", func() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, ContextConfig, config)
	ctx = context.WithValue(ctx, ContextDb, mockDb{})

	BeforeSuite(func() {
		gin.SetMode(gin.TestMode)
	})

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
				w := performRequest(ctx, "GET", "/v1/payments/?count=20&after=1")
				Expect(w.Code).To(Equal(http.StatusOK))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(r).To(MatchJSON(`{"data": [], "links": {"self": "http://example.com/v1/payments/?count=20&after=1", "next": null}}`))
			})
			It("should list payments from database", func() {
				c := context.WithValue(ctx, ContextDb, mockDb{
					Payments: []Payment{
						{Id: "1"},
					}})
				w := performRequest(c, "GET", "/v1/payments/?count=10")
				Expect(w.Code).To(Equal(http.StatusOK))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(r).To(MatchJSON(`{
					"data": [
						{"id": "1", "links": {"self": "http://example.com/v1/payments/1/"}}], 
					"links": {"self": "http://example.com/v1/payments/?count=10", "next": null}}`))
			})
			It("should have next link if there are too many items", func() {
				c := context.WithValue(ctx, ContextDb, mockDb{
					Payments: []Payment{
						{Id: "a"},
						{Id: "b"},
					}})
				w := performRequest(c, "GET", "/v1/payments/?count=1")
				Expect(w.Code).To(Equal(http.StatusOK))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(r).To(MatchJSON(`{
					"data": [
						{"id": "a", "links": {"self": "http://example.com/v1/payments/a/"}}], 
					"links": {"self": "http://example.com/v1/payments/?count=1", "next": "http://example.com/v1/payments/?count=1&after=a"}}`))
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
						{
							Id:             "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
							Version:        0,
							OrganisationId: "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
							Attributes: PaymentAttributes{
								Amount: "100.21",
								BeneficiaryParty: PaymentParty{
									AccountName:       "W Owens",
									AccountNumber:     "31926819",
									AccountNumberCode: "BBAN",
									AccountType:       0,
									Address:           "1 The Beneficiary Localtown SE2",
									BankId:            "403000",
									BankIdCode:        "GBDSC",
									Name:              "Wilfred Jeremiah Owens",
								},
								ChargesInformation: PaymentChargesInformation{
									BearerCode: "SHAR",
									SenderCharges: []PaymentSenderCharge{
										{Amount: "5.00", Currency: "GBP"},
										{Amount: "10.00", Currency: "USD"},
									},
									ReceiverChargesAmount:   "1.00",
									ReceiverChargesCurrency: "USD",
								},
								Currency: "GBP",
								DebtorParty: PaymentParty{
									AccountName:       "EJ Brown Black",
									AccountNumber:     "GB29XABC10161234567801",
									AccountNumberCode: "IBAN",
									AccountType:       0,
									Address:           "10 Debtor Crescent Sourcetown NE1",
									BankId:            "203301",
									BankIdCode:        "GBDSC",
									Name:              "Emelia Jane Brown",
								},
								EndToEndReference: "Wil piano Jan",
								Fx: PaymentFx{
									ContractReference: "FX123",
									ExchangeRate:      "2.00000",
									OriginalAmount:    "200.42",
									OriginalCurrency:  "USD",
								},
								NumericReference:     "1002001",
								PaymentId:            "123456789012345678",
								PaymentPurpose:       "Paying for goods/services",
								PaymentScheme:        "FPS",
								PaymentType:          "Credit",
								ProcessingDate:       "2017-01-18",
								Reference:            "Payment for Em's piano lessons",
								SchemePaymentSubType: "InternetBanking",
								SchemePaymentType:    "ImmediatePayment",
								SponsorParty: PaymentSponsorParty{
									AccountNumber: "56781234",
									BankId:        "123123",
									BankIdCode:    "GBDSC",
								},
							},
						},
					}})
				w := performRequest(c, "GET", "/v1/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43")
				Expect(w.Code).To(Equal(http.StatusOK))
				r, _ := ioutil.ReadAll(w.Body)
				Expect(string(r)).To(MatchJSON(`
				{
					"type": "Payment",
				  	"id": "4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43",
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
					"links": {"self": "http://example.com/v1/payments/4ee3a8d8-ca7b-4290-a52c-dd5b6165ec43/"}
				}`))
			})
		})
	})
})
