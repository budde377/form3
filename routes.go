package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/logger"
	"net/http"
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
	var after *string
	if v := r.URL.Query().Get("after"); v != "" {
		after = &v
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
	mapped := make([]PaymentSummaryRest, resultLen)
	for i, v := range *summaries {
		if i >= resultLen {
			break
		}
		mapped[i] = summaryToRest(config, v)
	}

	// Create self link
	var afterStr = ""
	if after != nil {
		afterStr = fmt.Sprintf("&after=%s", *after)
	}
	var selfLink = fmt.Sprintf("%s/v1/payments/?count=%d%s", config.Host, size, afterStr)

	// Create next link
	var nextLink *string
	if len(*summaries) > size {
		n := fmt.Sprintf("%s/v1/payments/?count=%d&after=%s", config.Host, size, mapped[len(mapped)-1].Id)
		nextLink = &n
	}
	render.JSON(w, r, PaymentsDataRest{
		Data: mapped,
		Links: PageLinksRest{
			Self: selfLink,
			Next: nextLink,
		},
	})
}

func PaymentCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		db := ctx.Value(ContextDb).(Db)
		id := chi.URLParam(r, "paymentID")
		payment, err := db.GetPaymentById(ctx, id)
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

func PaymentRoute() http.Handler {
	r := chi.NewRouter()
	r.Use(PaymentCtx)
	r.Get("/", getPaymentEndpoint)
	return r
}

func PaymentsRoute() http.Handler {
	r := chi.NewRouter()
	r.Get("/", listPaymentsEndpoint)
	return r
}

func V1Route() http.Handler {
	r := chi.NewRouter()
	r.Mount("/v1/payments", PaymentsRoute())
	r.Mount("/v1/payments/{paymentID}", PaymentRoute())
	return r
}

func okEndpoint(w http.ResponseWriter, r *http.Request) {
	render.PlainText(w, r, "OK")
}

func RootRoute() http.Handler {
	r := chi.NewRouter()
	r.Get("/", okEndpoint)
	r.Mount("/", V1Route())
	return r
}
