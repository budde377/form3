// +build integration

package main_test

import (
	. "./"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DbIntegration", func() {
	BeforeEach(func() {
		_ = db.(TestDb).Drop(testCtx)
	})

	Describe("GetPayments", func() {
		BeforeEach(func() {
			_ = db.(TestDb).Drop(testCtx)
			_ = populateDatabase(100)
		})
		It("should list all payments", func() {
			res, err := db.GetPayments(testCtx, 150, nil)
			Expect(err).To(BeNil())
			Expect(*res).To(HaveLen(100))
		})
		It("should limit the output", func() {
			res, err := db.GetPayments(testCtx, 10, nil)
			Expect(err).To(BeNil())
			Expect(*res).To(HaveLen(10))
		})
		It("should fetch after", func() {
			res1, err := db.GetPayments(testCtx, 100, nil)
			Expect(err).To(BeNil())
			Expect(*res1).To(HaveLen(100))
			id := (*res1)[10].ID
			res2, err := db.GetPayments(testCtx, 2, &id)
			Expect(err).To(BeNil())
			Expect(*res2).To(Equal([]PaymentSummary{
				{
					ID: (*res1)[11].ID,
				},
				{
					ID: (*res1)[12].ID,
				},
			}))
		})
	})

	Describe("CreatePayment", func() {
		It("should create and return id", func() {
			orgId := "org"
			id, err := db.CreatePayment(testCtx, orgId, paymentSample.Attributes)
			Expect(err).To(BeNil())
			Expect(id).ToNot(BeNil())
		})
		It("should be possible to fetch newly created resource", func() {
			orgId := "org"
			id, err := db.CreatePayment(testCtx, orgId, paymentSample.Attributes)
			Expect(err).To(BeNil())
			Expect(id).ToNot(BeNil())
			payment, _ := db.GetPaymentByID(testCtx, *id)
			Expect(*payment).To(Equal(Payment{
				ID:             *id,
				OrganisationID: orgId,
				Version:        0,
				Attributes:     paymentSample.Attributes,
			}))
		})
	})
	Describe("UpdatePayment", func() {
		It("should update organization", func() {
			id, _ := db.CreatePayment(testCtx, paymentSample.OrganisationID, paymentSample.Attributes)
			err := db.UpdatePayment(testCtx, *id, "org", paymentSample.Attributes)
			Expect(err).To(BeNil())
			payment, _ := db.GetPaymentByID(testCtx, *id)
			Expect(*payment).To(Equal(Payment{
				Version:        1,
				OrganisationID: "org",
				Attributes:     paymentSample.Attributes,
				ID:             *id,
			}))
		})
		It("should update on non-existing-id", func() {
			id, _ := StringToID("aaaaaaaaaaaaaaaaaaaaaaaa")
			err := db.UpdatePayment(testCtx, *id, "org", paymentSample.Attributes)
			Expect(err).To(BeNil())
			payment, _ := db.GetPaymentByID(testCtx, *id)
			Expect(payment).To(BeNil())
		})
	})
	Describe("DeletePayment", func() {
		It("should delete an existing payment", func() {
			id, _ := db.CreatePayment(testCtx, paymentSample.OrganisationID, paymentSample.Attributes)
			err := db.DeletePayment(testCtx, *id)
			Expect(err).To(BeNil())
			payment, _ := db.GetPaymentByID(testCtx, *id)
			Expect(payment).To(BeNil())
		})
		It("should delete a non-existing-id without failure", func() {
			id, _ := StringToID("aaaaaaaaaaaaaaaaaaaaaaaa")
			err := db.UpdatePayment(testCtx, *id, "org", paymentSample.Attributes)
			Expect(err).To(BeNil())
			payment, _ := db.GetPaymentByID(testCtx, *id)
			Expect(payment).To(BeNil())
		})
	})
})
