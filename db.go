package main

import (
	"context"
)

type PaymentSummary struct {
	Id string
}

type PaymentParty struct {
	AccountName       string
	AccountNumber     string
	AccountNumberCode string
	AccountType       int
	Address           string
	BankId            string
	BankIdCode        string
	Name              string
}

type PaymentSenderCharge struct {
	Amount   string
	Currency string
}

type PaymentChargesInformation struct {
	BearerCode              string
	ReceiverChargesAmount   string
	ReceiverChargesCurrency string
	SenderCharges           []PaymentSenderCharge
}

type PaymentFx struct {
	ContractReference string
	ExchangeRate      string
	OriginalAmount    string
	OriginalCurrency  string
}

type PaymentSponsorParty struct {
	AccountNumber string
	BankId        string
	BankIdCode    string
}

type PaymentAttributes struct {
	Amount               string
	BeneficiaryParty     PaymentParty
	ChargesInformation   PaymentChargesInformation
	Currency             string
	DebtorParty          PaymentParty
	EndToEndReference    string
	Fx                   PaymentFx
	NumericReference     string
	PaymentId            string
	PaymentPurpose       string
	PaymentScheme        string
	PaymentType          string
	ProcessingDate       string
	Reference            string
	SchemePaymentSubType string
	SchemePaymentType    string
	SponsorParty         PaymentSponsorParty
}

type Payment struct {
	Id             string
	OrganisationId string
	Version        int
	Attributes     PaymentAttributes
}

// A database abstraction responsible for all retrieval and modification of
// persistent storage.
type Db interface {
	// Retrieve a list of payments
	GetPayments(ctx context.Context, size int, after *string) (*[]PaymentSummary, error)

	// Retrieve a single payment
	GetPaymentById(ctx context.Context, id string) (*Payment, error)
}

type db struct {
}

func (*db) GetPayments(ctx context.Context, size int, after *string) (*[]PaymentSummary, error) {
	panic("implement me")
}

func (*db) GetPaymentById(ctx context.Context, id string) (*Payment, error) {
	panic("implement me")
}

func InitDb(config *Config) Db {
	return &db{}
}
