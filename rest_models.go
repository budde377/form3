package main

import "fmt"

type selfLinksRest struct {
	Self string `json:"self"`
}

type paymentSummaryRest struct {
	ID    string        `json:"id"`
	Links selfLinksRest `json:"links"`
}

type paymentPartyRest struct {
	AccountName       string `json:"account_name"`
	AccountNumber     string `json:"account_number"`
	AccountNumberCode string `json:"account_number_code"`
	AccountType       int    `json:"account_type"`
	Address           string `json:"address"`
	BankID            string `json:"bank_id"`
	BankIDCode        string `json:"bank_id_code"`
	Name              string `json:"name"`
}

type paymentChargeRest struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

type paymentChargesInformationRest struct {
	BearerCode              string              `json:"bearer_code"`
	ReceiverChargesAmount   string              `json:"receiver_charges_amount"`
	ReceiverChargesCurrency string              `json:"receiver_charges_currency"`
	SenderCharges           []paymentChargeRest `json:"sender_charges"`
}

type paymentFxRest struct {
	ContractReference string `json:"contract_reference"`
	ExchangeRate      string `json:"exchange_rate"`
	OriginalAmount    string `json:"original_amount"`
	OriginalCurrency  string `json:"original_currency"`
}

type paymentSponsorPartyRest struct {
	AccountNumber string `json:"account_number"`
	BankID        string `json:"bank_id"`
	BankIDCode    string `json:"bank_id_code"`
}

type paymentAttributesRest struct {
	Amount               string                        `json:"amount"`
	BeneficiaryParty     paymentPartyRest              `json:"beneficiary_party"`
	ChargesInformation   paymentChargesInformationRest `json:"charges_information"`
	Currency             string                        `json:"currency"`
	DebtorParty          paymentPartyRest              `json:"debtor_party"`
	EndToEndReference    string                        `json:"end_to_end_reference"`
	Fx                   paymentFxRest                 `json:"fx"`
	NumericReference     string                        `json:"numeric_reference"`
	PaymentID            string                        `json:"payment_id"`
	PaymentPurpose       string                        `json:"payment_purpose"`
	PaymentScheme        string                        `json:"payment_scheme"`
	PaymentType          string                        `json:"payment_type"`
	ProcessingDate       string                        `json:"processing_date"`
	Reference            string                        `json:"reference"`
	SchemePaymentSubType string                        `json:"scheme_payment_sub_type"`
	SchemePaymentType    string                        `json:"scheme_payment_type"`
	SponsorParty         paymentSponsorPartyRest       `json:"sponsor_party"`
}

type paymentRest struct {
	ID             string                `json:"id"`
	OrganisationID string                `json:"organisation_id"`
	Version        int                   `json:"version"`
	Attributes     paymentAttributesRest `json:"attributes"`
	Links          selfLinksRest         `json:"links"`
	Type           string                `json:"type"`
}

type pageLinksRest struct {
	Self string  `json:"self"`
	Next *string `json:"next"`
}

type paymentsDataRest struct {
	Data  []paymentSummaryRest `json:"data"`
	Links pageLinksRest        `json:"links"`
}

func summaryToRest(config *Config, summary PaymentSummary) paymentSummaryRest {
	id := summary.ID
	return summaryIDToRest(config, id)
}

func summaryIDToRest(config *Config, id ID) paymentSummaryRest {
	return paymentSummaryRest{
		ID: IDToString(id),
		Links: selfLinksRest{
			Self: fmt.Sprintf("%s/v1/payments/%s/", config.Host, IDToString(id)),
		},
	}

}

func partyToRest(party PaymentParty) paymentPartyRest {
	return paymentPartyRest{
		AccountNumber:     party.AccountNumber,
		AccountName:       party.AccountName,
		AccountNumberCode: party.AccountNumberCode,
		AccountType:       party.AccountType,
		Address:           party.Address,
		BankID:            party.BankID,
		BankIDCode:        party.BankIDCode,
		Name:              party.Name,
	}
}
func partyFromRest(party paymentPartyRest) PaymentParty {
	return PaymentParty{
		AccountNumber:     party.AccountNumber,
		AccountName:       party.AccountName,
		AccountNumberCode: party.AccountNumberCode,
		AccountType:       party.AccountType,
		Address:           party.Address,
		BankID:            party.BankID,
		BankIDCode:        party.BankIDCode,
		Name:              party.Name,
	}
}

func paymentChargesInformationToRest(info PaymentChargesInformation) paymentChargesInformationRest {
	mappedCharges := make([]paymentChargeRest, len(info.SenderCharges))
	for i, v := range info.SenderCharges {
		mappedCharges[i] = paymentChargeRest{
			Amount:   v.Amount,
			Currency: v.Currency,
		}
	}
	return paymentChargesInformationRest{
		SenderCharges:           mappedCharges,
		BearerCode:              info.BearerCode,
		ReceiverChargesCurrency: info.ReceiverChargesCurrency,
		ReceiverChargesAmount:   info.ReceiverChargesAmount,
	}
}

func paymentChargesInformationFromRest(info paymentChargesInformationRest) PaymentChargesInformation {
	mappedCharges := make([]PaymentSenderCharge, len(info.SenderCharges))
	for i, v := range info.SenderCharges {
		mappedCharges[i] = PaymentSenderCharge{
			Amount:   v.Amount,
			Currency: v.Currency,
		}
	}
	return PaymentChargesInformation{
		SenderCharges:           mappedCharges,
		BearerCode:              info.BearerCode,
		ReceiverChargesCurrency: info.ReceiverChargesCurrency,
		ReceiverChargesAmount:   info.ReceiverChargesAmount,
	}
}

func paymentAttributesToRest(attributes PaymentAttributes) paymentAttributesRest {
	return paymentAttributesRest{
		Amount:             attributes.Amount,
		BeneficiaryParty:   partyToRest(attributes.BeneficiaryParty),
		ChargesInformation: paymentChargesInformationToRest(attributes.ChargesInformation),
		Currency:           attributes.Currency,
		DebtorParty:        partyToRest(attributes.DebtorParty),
		EndToEndReference:  attributes.EndToEndReference,
		Fx: paymentFxRest{
			ContractReference: attributes.Fx.ContractReference,
			ExchangeRate:      attributes.Fx.ExchangeRate,
			OriginalAmount:    attributes.Fx.OriginalAmount,
			OriginalCurrency:  attributes.Fx.OriginalCurrency,
		},
		NumericReference:     attributes.NumericReference,
		PaymentID:            attributes.PaymentID,
		PaymentPurpose:       attributes.PaymentPurpose,
		PaymentScheme:        attributes.PaymentScheme,
		PaymentType:          attributes.PaymentType,
		ProcessingDate:       attributes.ProcessingDate,
		Reference:            attributes.Reference,
		SchemePaymentSubType: attributes.SchemePaymentSubType,
		SchemePaymentType:    attributes.SchemePaymentType,
		SponsorParty: paymentSponsorPartyRest{
			AccountNumber: attributes.SponsorParty.AccountNumber,
			BankID:        attributes.SponsorParty.BankID,
			BankIDCode:    attributes.SponsorParty.BankIDCode,
		},
	}
}

func paymentAttributesFromRest(attributes paymentAttributesRest) PaymentAttributes {
	return PaymentAttributes{
		Amount:             attributes.Amount,
		BeneficiaryParty:   partyFromRest(attributes.BeneficiaryParty),
		ChargesInformation: paymentChargesInformationFromRest(attributes.ChargesInformation),
		Currency:           attributes.Currency,
		DebtorParty:        partyFromRest(attributes.DebtorParty),
		EndToEndReference:  attributes.EndToEndReference,
		Fx: PaymentFx{
			ContractReference: attributes.Fx.ContractReference,
			ExchangeRate:      attributes.Fx.ExchangeRate,
			OriginalAmount:    attributes.Fx.OriginalAmount,
			OriginalCurrency:  attributes.Fx.OriginalCurrency,
		},
		NumericReference:     attributes.NumericReference,
		PaymentID:            attributes.PaymentID,
		PaymentPurpose:       attributes.PaymentPurpose,
		PaymentScheme:        attributes.PaymentScheme,
		PaymentType:          attributes.PaymentType,
		ProcessingDate:       attributes.ProcessingDate,
		Reference:            attributes.Reference,
		SchemePaymentSubType: attributes.SchemePaymentSubType,
		SchemePaymentType:    attributes.SchemePaymentType,
		SponsorParty: PaymentSponsorParty{
			AccountNumber: attributes.SponsorParty.AccountNumber,
			BankID:        attributes.SponsorParty.BankID,
			BankIDCode:    attributes.SponsorParty.BankIDCode,
		},
	}
}

func paymentToRest(config *Config, payment Payment) paymentRest {
	id := payment.ID
	return paymentRest{
		ID:             IDToString(id),
		OrganisationID: payment.OrganisationID,
		Attributes:     paymentAttributesToRest(payment.Attributes),
		Version:        payment.Version,
		Type:           "Payment",
		Links: selfLinksRest{
			Self: fmt.Sprintf("%s/v1/payments/%s/", config.Host, IDToString(id)),
		},
	}

}
