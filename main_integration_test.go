// +build integration

package main_test

import (
	. "./"
	"context"
	. "github.com/onsi/ginkgo"
)

type TestDb interface {
	Db
	Drop(ctx context.Context) error
}

var db, _ = NewDb(&testConfig)

var _ = BeforeSuite(func() {
	_ = db.Connect(testCtx)
})

var _ = AfterSuite(func() {
	_ = db.Close(testCtx)
})

var testDbCtx = context.WithValue(testCtx, ContextDb, db)

func populateDatabase(n int) []ID {
	var ids = make([]ID, n)
	for i := 0; i < n; i++ {
		id, _ := db.CreatePayment(testCtx, paymentSample.OrganisationID, paymentSample.Attributes)
		ids[i] = *id
	}
	return ids
}
