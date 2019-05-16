package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/logger"
	"io/ioutil"
	"net/http"
)

type key int

const (
	ContextConfig  key = iota
	ContextDb      key = iota
	ContextPayment key = iota
)

func main() {
	defer logger.Init("Form3 API", true, false, ioutil.Discard).Close()
	c := OpenConfig()
	db := InitDb(c)
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
