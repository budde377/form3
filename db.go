package main

type PaymentSummary struct {
	Id string
}

type Payment struct {
	Id string
}

type Db interface {
	GetPayments(size int, after *string) (*[]PaymentSummary, error)
	GetPaymentById(id string) (*Payment, error)
}

type db struct {

}

func (*db) GetPayments(size int, after *string) (*[]PaymentSummary, error) {
	panic("implement me")
}

func (*db) GetPaymentById(id string) (*Payment, error) {
	panic("implement me")
}

func InitDb(config *Config) Db {
	return &db {

	}
}
