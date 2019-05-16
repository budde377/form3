package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/logger"
)

type key int

const (
	// ContextConfig key used to fetch config from context
	ContextConfig key = iota
	// ContextDb key used to fetch db from context
	ContextDb key = iota
	// ContextPayment key used to fetch current payment from context
	ContextPayment key = iota
)

func main() {
	defer logger.Init("Form3 API", true, false, ioutil.Discard).Close()
	c := OpenConfigFile("./config/config.json")
	db, err := NewDb(c)
	if err != nil {
		logger.Fatal("failed to initialize database: ", err)
	}
	r := chi.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), ContextDb, db)
			ctx = context.WithValue(ctx, ContextConfig, c)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})
	r.Mount("/", RootRoute())
	if err := http.ListenAndServe(fmt.Sprintf(":%d", c.Port), r); err != nil {
		logger.Fatal("failed to start server ", err)
	}

}
