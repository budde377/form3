package main_test

import (
	. "./"
	"context"
	"io"
	"net/http/httptest"
)

var testConfig = OpenConfigFile("./test/config.json")
var testCtx = context.WithValue(context.Background(), ContextConfig, testConfig)

var id, _ = StringToID("5cdd382e9549af35c3b94301")
var id2, _ = StringToID("5cdd382e9549af35c3b94302")

var paymentSample = Payment{
	ID:             *id,
	Version:        0,
	OrganisationID: "743d5b63-8e6f-432e-a8fa-c5d8d2ee5fcb",
	Attributes: PaymentAttributes{
		Amount: "100.21",
		BeneficiaryParty: PaymentParty{
			AccountName:       "W Owens",
			AccountNumber:     "31926819",
			AccountNumberCode: "BBAN",
			AccountType:       0,
			Address:           "1 The Beneficiary Localtown SE2",
			BankID:            "403000",
			BankIDCode:        "GBDSC",
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
			BankID:            "203301",
			BankIDCode:        "GBDSC",
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
		PaymentID:            "123456789012345678",
		PaymentPurpose:       "Paying for goods/services",
		PaymentScheme:        "FPS",
		PaymentType:          "Credit",
		ProcessingDate:       "2017-01-18",
		Reference:            "Payment for Em's piano lessons",
		SchemePaymentSubType: "InternetBanking",
		SchemePaymentType:    "ImmediatePayment",
		SponsorParty: PaymentSponsorParty{
			AccountNumber: "56781234",
			BankID:        "123123",
			BankIDCode:    "GBDSC",
		},
	},
}

var paymentSampleAttributesJSON = `
{
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
}`

type mockDb struct {
	Payments []Payment
	error    error
}

func (d mockDb) DeletePayment(ctx context.Context, id ID) error {
	return d.error
}

func (d mockDb) UpdatePayment(ctx context.Context, id ID, organizationId string, attributes PaymentAttributes) error {
	return d.error
}

func (d mockDb) CreatePayment(ctx context.Context, organizationId string, attributes PaymentAttributes) (*ID, error) {
	if d.error != nil {
		return nil, d.error
	}
	return StringToID("5cdd382e9549af35c3b94301")
}

func (d mockDb) Connect(ctx context.Context) error {
	return d.error
}

func (d mockDb) Close(ctx context.Context) error {
	return d.error
}

func (d mockDb) GetPayments(ctx context.Context, size int, after *ID) (*[]PaymentSummary, error) {
	summaries := make([]PaymentSummary, 0)
	for _, v := range d.Payments {
		summaries = append(summaries, PaymentSummary{
			ID: v.ID,
		})
	}
	return &summaries, nil
}

func (d mockDb) GetPaymentByID(ctx context.Context, id ID) (*Payment, error) {
	for _, v := range d.Payments {
		if v.ID == id {
			return &v, nil
		}
	}
	return nil, nil
}

func performRequest(ctx context.Context, method, path string) *httptest.ResponseRecorder {
	return performRequestBody(ctx, method, path, nil)
}

func performRequestBody(ctx context.Context, method, path string, body io.Reader) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, body)
	req.Header.Set("content-type", "application/json")
	w := httptest.NewRecorder()
	RootRoute().ServeHTTP(w, req.WithContext(ctx))
	return w
}