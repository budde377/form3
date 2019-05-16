package main

import "fmt"

type SelfLinksRest struct {
	Self string `json:"self"`
}

type PaymentSummaryRest struct {
	Id    string        `json:"id"`
	Links SelfLinksRest `json:"links"`
}

type PaymentPartyRest struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	AccountType       int    `json:"account_type"`
	Address           string `json:"address"`
	BankId            string `json:"bank_id"`
	BankIdCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}

type PaymentChargeRest struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type PaymentChargesInformationRest struct {
	BearerCode              string              `json:"bearer_code"`
	ReceiverChargesAmount   string              `json:"receiver_charges_amount"`
	ReceiverChargesCurrency string              `json:"receiver_charges_currency"`
	SenderCharges           []PaymentChargeRest `json:"sender_charges"`
}

type PaymentFxRest struct {
	ContractReference string `json:"contract_reference"`
	ExchangeRate      string `json:"exchange_rate"`
	OriginalAmount    string `json:"original_amount"`
	OriginalCurrency  string `json:"original_currency"`
}

type PaymentSponsorPartyRest struct {
	AccountNumber string `json:"account_number"`
	BankId        string `json:"bank_id"`
	BankIdCode    string `json:"bank_id_code"`
}

type PaymentAttributesRest struct {
	Amount               string                        `json:"amount"`
	BeneficiaryParty     PaymentPartyRest              `json:"beneficiary_party"`
	ChargesInformation   PaymentChargesInformationRest `json:"charges_information"`
	Currency             string                        `json:"currency"`
	DebtorParty          PaymentPartyRest              `json:"debtor_party"`
	EndToEndReference    string                        `json:"end_to_end_reference"`
	Fx                   PaymentFxRest                 `json:"fx"`
	NumericReference     string                        `json:"numeric_reference"`
	PaymentId            string                        `json:"payment_id"`
	PaymentPurpose       string                        `json:"payment_purpose"`
	PaymentScheme        string                        `json:"payment_scheme"`
	PaymentType          string                        `json:"payment_type"`
	ProcessingDate       string                        `json:"processing_date"`
	Reference            string                        `json:"reference"`
	SchemePaymentSubType string                        `json:"scheme_payment_sub_type"`
	SchemePaymentType    string                        `json:"scheme_payment_type"`
	SponsorParty         PaymentSponsorPartyRest       `json:"sponsor_party"`
}

type PaymentRest struct {
	Id             string                `json:"id"`
	OrganisationId string                `json:"organisation_id"`
	Version        int                   `json:"version"`
	Attributes     PaymentAttributesRest `json:"attributes"`
	Links          SelfLinksRest         `json:"links"`
	Type           string                `json:"type"`
}

type PageLinksRest struct {
	Self string  `json:"self"`
	Next *string `json:"next"`
}

type PaymentsDataRest struct {
	Data  []PaymentSummaryRest `json:"data"`
	Links PageLinksRest        `json:"links"`
}

func summaryToRest(config *Config, summary PaymentSummary) PaymentSummaryRest {
	id := summary.Id
	return PaymentSummaryRest{
		Id: id,
		Links: SelfLinksRest{
			Self: fmt.Sprintf("%s/v1/payments/%s/", config.Host, id),
		},
	}

}

func partyToRest(party PaymentParty) PaymentPartyRest {
	return PaymentPartyRest{
		AccountNumber:     party.AccountNumber,
		AccountName:       party.AccountName,
		AccountNumberCode: party.AccountNumberCode,
		AccountType:       party.AccountType,
		Address:           party.Address,
		BankId:            party.BankId,
		BankIdCode:        party.BankIdCode,
		Name:              party.Name,
	}
}

func paymentChargesInformationToRest(info PaymentChargesInformation) PaymentChargesInformationRest {
	mappedCharges := make([]PaymentChargeRest, len(info.SenderCharges))
	for i, v := range info.SenderCharges {
		mappedCharges[i] = PaymentChargeRest{
			Amount:   v.Amount,
			Currency: v.Currency,
		}
	}
	return PaymentChargesInformationRest{
		SenderCharges:           mappedCharges,
		BearerCode:              info.BearerCode,
		ReceiverChargesCurrency: info.ReceiverChargesCurrency,
		ReceiverChargesAmount:   info.ReceiverChargesAmount,
	}
}

func paymentAttributesToRest(attributes PaymentAttributes) PaymentAttributesRest {
	return PaymentAttributesRest{
		Amount:             attributes.Amount,
		BeneficiaryParty:   partyToRest(attributes.BeneficiaryParty),
		ChargesInformation: paymentChargesInformationToRest(attributes.ChargesInformation),
		Currency:           attributes.Currency,
		DebtorParty:        partyToRest(attributes.DebtorParty),
		EndToEndReference:  attributes.EndToEndReference,
		Fx: PaymentFxRest{
			ContractReference: attributes.Fx.ContractReference,
			ExchangeRate:      attributes.Fx.ExchangeRate,
			OriginalAmount:    attributes.Fx.OriginalAmount,
			OriginalCurrency:  attributes.Fx.OriginalCurrency,
		},
		NumericReference:     attributes.NumericReference,
		PaymentId:            attributes.PaymentId,
		PaymentPurpose:       attributes.PaymentPurpose,
		PaymentScheme:        attributes.PaymentScheme,
		PaymentType:          attributes.PaymentType,
		ProcessingDate:       attributes.ProcessingDate,
		Reference:            attributes.Reference,
		SchemePaymentSubType: attributes.SchemePaymentSubType,
		SchemePaymentType:    attributes.SchemePaymentType,
		SponsorParty: PaymentSponsorPartyRest{
			AccountNumber: attributes.SponsorParty.AccountNumber,
			BankId:        attributes.SponsorParty.BankId,
			BankIdCode:    attributes.SponsorParty.BankIdCode,
		},
	}
}

func paymentToRest(config *Config, payment Payment) PaymentRest {
	id := payment.Id
	return PaymentRest{
		Id:             id,
		OrganisationId: payment.OrganisationId,
		Attributes:     paymentAttributesToRest(payment.Attributes),
		Type:           "Payment",
		Links: SelfLinksRest{
			Self: fmt.Sprintf("%s/v1/payments/%s/", config.Host, id),
		},
	}

}
