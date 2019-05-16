package main

import (
	"context"
	"github.com/google/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ID is an unique identifier in the database
type ID = primitive.ObjectID

// StringToID converts a string to an ID
func StringToID(id string) (*ID, error) {
	v, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// IDToString converts an ID to a string
func IDToString(id ID) string {
	return id.Hex()
}

// PaymentSummary is a simplification of the payment resource
type PaymentSummary struct {
	ID ID `bson:"_id"`
}

// PaymentParty ...
type PaymentParty struct {
	AccountName       string `bson:"account_name"`
	AccountNumber     string `bson:"account_number"`
	AccountNumberCode string `bson:"account_number_code"`
	AccountType       int    `bson:"account_type"`
	Address           string `bson:"address"`
	BankID            string `bson:"bank_id"`
	BankIDCode        string `bson:"bank_id_code"`
	Name              string `bson:"name"`
}

// PaymentSenderCharge ...
type PaymentSenderCharge struct {
	Amount   string `bson:"amount"`
	Currency string `bson:"currency"`
}

// PaymentChargesInformation ...
type PaymentChargesInformation struct {
	BearerCode              string                `bson:"bearer_code"`
	ReceiverChargesAmount   string                `bson:"receiver_charges_amount"`
	ReceiverChargesCurrency string                `bson:"receiver_charges_currency"`
	SenderCharges           []PaymentSenderCharge `bson:"sender_charges"`
}

// PaymentFx ...
type PaymentFx struct {
	ContractReference string `bson:"contract_reference"`
	ExchangeRate      string `bson:"exchange_rate"`
	OriginalAmount    string `bson:"original_amount"`
	OriginalCurrency  string `bson:"original_currency"`
}

// PaymentSponsorParty ...
type PaymentSponsorParty struct {
	AccountNumber string `bson:"account_number"`
	BankID        string `bson:"bank_id"`
	BankIDCode    string `bson:"bank_id_code"`
}

// PaymentAttributes are attributes associated with a payment
type PaymentAttributes struct {
	Amount               string                    `bson:"amount"`
	BeneficiaryParty     PaymentParty              `bson:"beneficiary_party"`
	ChargesInformation   PaymentChargesInformation `bson:"charges_information"`
	Currency             string                    `bson:"currency"`
	DebtorParty          PaymentParty              `bson:"debtor_party"`
	EndToEndReference    string                    `bson:"end_to_end_reference"`
	Fx                   PaymentFx                 `bson:"fx"`
	NumericReference     string                    `bson:"numeric_reference"`
	PaymentID            string                    `bson:"payment_id"`
	PaymentPurpose       string                    `bson:"payment_purpose"`
	PaymentScheme        string                    `bson:"payment_scheme"`
	PaymentType          string                    `bson:"payment_type"`
	ProcessingDate       string                    `bson:"processing_date"`
	Reference            string                    `bson:"reference"`
	SchemePaymentSubType string                    `bson:"scheme_payment_sub_type"`
	SchemePaymentType    string                    `bson:"scheme_payment_type"`
	SponsorParty         PaymentSponsorParty       `bson:"sponsor_party"`
}

// Payment is a payment representation
type Payment struct {
	ID             ID                `bson:"_id"`
	OrganisationID string            `bson:"organisation_id"`
	Version        int               `bson:"version"`
	Attributes     PaymentAttributes `bson:"attributes"`
}

type cPayment struct {
	OrganisationID string            `bson:"organisation_id"`
	Version        int               `bson:"version"`
	Attributes     PaymentAttributes `bson:"attributes"`
}

// Db is an abstraction responsible for all retrieval and modification of
// persistent storage.
type Db interface {
	// Retrieve a list of payments
	GetPayments(ctx context.Context, size int, after *ID) (*[]PaymentSummary, error)

	// Retrieve a single payment
	GetPaymentByID(ctx context.Context, id ID) (*Payment, error)

	// Create a new payment
	CreatePayment(ctx context.Context, organizationID string, attributes PaymentAttributes) (*ID, error)

	// Update a new payment
	UpdatePayment(ctx context.Context, ID ID, organizationID string, attributes PaymentAttributes) error

	// Delete a payment for good
	DeletePayment(ctx context.Context, ID ID) error

	// Connect to database
	Connect(ctx context.Context) error

	// Close connection
	Close(ctx context.Context) error
}

type db struct {
	Client *mongo.Client
}

func (db *db) DeletePayment(ctx context.Context, id ID) error {
	_, err := db.paymentsCollection(ctx).DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (db *db) UpdatePayment(ctx context.Context, id ID, organisationID string, attributes PaymentAttributes) error {
	_, err := db.paymentsCollection(ctx).UpdateMany(ctx, bson.M{"_id": id}, bson.M{
		"$set": bson.M{
			"organisation_id": organisationID,
			"attributes":      attributes,
		},
		"$inc": bson.M{"version": 1},
	})
	return err
}

func (db *db) paymentsCollection(ctx context.Context) *mongo.Collection {
	conf := ctx.Value(ContextConfig).(*Config)
	return db.Client.Database(conf.MongoDbDatabase).Collection("payments")
}

func (db *db) CreatePayment(ctx context.Context, organizationID string, attributes PaymentAttributes) (*ID, error) {
	res, err := db.paymentsCollection(ctx).InsertOne(
		ctx, cPayment{
			OrganisationID: organizationID,
			Attributes:     attributes,
			Version:        0,
		},
	)
	if err != nil {
		logger.Fatal(err)
	}
	str := res.InsertedID.(primitive.ObjectID)
	return &str, nil
}

func (db *db) Connect(ctx context.Context) error {
	return db.Client.Connect(ctx)
}

func (db *db) Close(ctx context.Context) error {
	return db.Client.Disconnect(ctx)
}

func (db *db) Drop(ctx context.Context) error {
	return db.paymentsCollection(ctx).Drop(ctx)
}

func (db *db) GetPayments(ctx context.Context, size int, after *ID) (*[]PaymentSummary, error) {
	opts := options.FindOptions{
		Projection: bson.M{"_id": 1},
		Sort:       bson.M{"_id": 1},
	}
	filter := bson.M{}
	if after != nil {
		filter["_id"] = bson.M{
			"$gt": *after,
		}
	}
	cur, err := db.paymentsCollection(ctx).Find(ctx, filter, &opts)
	if err != nil {
		logger.Fatal(err)
	}
	defer cur.Close(ctx)
	var res []PaymentSummary
	for cur.Next(ctx) {
		var elm PaymentSummary
		err = cur.Decode(&elm)
		if err != nil {
			return nil, err
		}
		res = append(res, elm)
		if len(res) >= size {
			break
		}
	}
	return &res, nil
}

func (db *db) GetPaymentByID(ctx context.Context, id ID) (*Payment, error) {
	res := db.paymentsCollection(ctx).FindOne(ctx, bson.M{"_id": id})
	payment := Payment{}
	err := res.Decode(&payment)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &payment, nil
}

// NewDb constructs a new Db wrapper
func NewDb(config *Config) (Db, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoDbURI))
	return &db{Client: client}, err
}
