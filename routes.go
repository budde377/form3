package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/logger"
)

func getPaymentEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	conf := ctx.Value(ContextConfig).(*Config)
	payment := ctx.Value(ContextPayment).(*Payment)
	render.JSON(w, r, paymentToRest(conf, *payment))
}

func listPaymentsEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value(ContextConfig).(*Config)
	db := ctx.Value(ContextDb).(Db)
	// Extract size, default 10
	size := SafeStringToInt(r.URL.Query().Get("count"), 10)
	if size <= 0 {
		size = 10
	}
	// Extract after
	var after *ID
	if v, err := StringToID(r.URL.Query().Get("after")); err == nil {
		after = v
	}

	// Fetch summaries
	summaries, err := db.GetPayments(ctx, size+1, after)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		logger.Error("failed to list payments: ", err)
		return
	}

	// Transform summaries
	resultLen := IntMin(size, len(*summaries))
	mapped := make([]paymentSummaryRest, resultLen)
	for i, v := range *summaries {
		if i >= resultLen {
			break
		}
		mapped[i] = summaryToRest(config, v)
	}

	// Create self link
	var afterStr = ""
	if after != nil {
		afterStr = fmt.Sprintf("&after=%s", IDToString(*after))
	}
	var selfLink = fmt.Sprintf("%s/v1/payments/?count=%d%s", config.Host, size, afterStr)

	// Create next link
	var nextLink *string
	if len(*summaries) > size {
		n := fmt.Sprintf("%s/v1/payments/?count=%d&after=%s", config.Host, size, mapped[len(mapped)-1].ID)
		nextLink = &n
	}
	render.JSON(w, r, paymentsDataRest{
		Data: mapped,
		Links: pageLinksRest{
			Self: selfLink,
			Next: nextLink,
		},
	})
}

func paymentCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		db := ctx.Value(ContextDb).(Db)
		queryID := chi.URLParam(r, "paymentID")
		cID, err := StringToID(queryID)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		payment, err := db.GetPaymentByID(ctx, *cID)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			logger.Error("failed to fetch payment: ", err)
			return
		}
		if payment == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		ctx = context.WithValue(r.Context(), ContextPayment, payment)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type paymentRequest struct {
	Attributes     paymentAttributesRest `json:"attributes"`
	OrganisationID string                `json:"organisation_id"`
}

func (u *paymentRequest) Bind(r *http.Request) error {
	// TODO create custom validation logic
	return nil
}

func updatePaymentEndpoint(w http.ResponseWriter, r *http.Request) {
	data := &paymentRequest{}
	if err := render.Bind(r, data); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	payment := ctx.Value(ContextPayment).(*Payment)
	db := ctx.Value(ContextDb).(Db)
	err := db.UpdatePayment(ctx, payment.ID, data.OrganisationID, paymentAttributesFromRest(data.Attributes))
	if err != nil {
		logger.Error("failed to update payment: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	render.NoContent(w, r)

}

func deletePaymentEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	payment := ctx.Value(ContextPayment).(*Payment)
	db := ctx.Value(ContextDb).(Db)
	err := db.DeletePayment(ctx, payment.ID)
	if err != nil {
		logger.Error("failed to delete payment: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	render.NoContent(w, r)

}

func createPaymentEndpoint(w http.ResponseWriter, r *http.Request) {
	data := &paymentRequest{}
	if err := render.Bind(r, data); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	c := ctx.Value(ContextConfig).(*Config)
	db := ctx.Value(ContextDb).(Db)
	id, err := db.CreatePayment(ctx, data.OrganisationID, paymentAttributesFromRest(data.Attributes))
	if err != nil {
		logger.Error("failed to create payment: ", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, summaryIDToRest(c, *id))
}

func paymentRoute() http.Handler {
	r := chi.NewRouter()
	r.Use(paymentCtx)
	r.Get("/", getPaymentEndpoint)
	r.Put("/", updatePaymentEndpoint)
	r.Delete("/", deletePaymentEndpoint)
	return r
}

func paymentsRoute() http.Handler {
	r := chi.NewRouter()
	r.Get("/", listPaymentsEndpoint)
	r.Post("/", createPaymentEndpoint)
	return r
}

func v1Route() http.Handler {
	r := chi.NewRouter()
	r.Mount("/v1/payments", paymentsRoute())
	r.Mount("/v1/payments/{paymentID}", paymentRoute())
	return r
}

func okEndpoint(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "OK")
}

// RootRoute construcs a route for the API
func RootRoute() http.Handler {
	r := chi.NewRouter()
	r.Get("/", okEndpoint)
	r.Mount("/", v1Route())
	return r
}
