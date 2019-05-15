package main_test

import (
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	. "./"
)

var config = OpenConfigFile("./test/config.json")

type mockDb struct {
	Payments map[string]Payment
}

func (d mockDb) GetPayments(size int, after *string) (*[]PaymentSummary, error) {
	summaries := make([]PaymentSummary, 0)
	for k := range d.Payments {
		summaries = append(summaries, PaymentSummary{
			Id: k,
		})
	}
	return &summaries, nil
}

func (mockDb) GetPaymentById(id string) (*Payment, error) {
	panic("implement me")
}

func performRequest(db Db, method, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	Routes(config, db).ServeHTTP(w, req)
	return w
}

var _ = Describe("Routes", func() {
	var db Db = mockDb{
	}

	BeforeSuite(func() {
		gin.SetMode(gin.TestMode)
	})

	Describe("GET /", func() {
		It("should succeed with OK", func() {
			w := performRequest(db, "GET", "/")
			Expect(w.Code).To(Equal(http.StatusOK))
			r, _ := ioutil.ReadAll(w.Body)
			Expect(string(r)).To(Equal("OK"))
		})
		It("should fail with invalid method", func() {
			w := performRequest(db, "POST", "/")
			Expect(w.Code).To(Equal(http.StatusNotFound))
			r, _ := ioutil.ReadAll(w.Body)
			Expect(string(r)).To(Equal("404 page not found"))
		})
	})
	Describe("GET /v1/payments/", func() {
		It("should list payments", func() {
			w := performRequest(db, "GET", "/v1/payments/")
			Expect(w.Code).To(Equal(http.StatusOK))
			r, _ := ioutil.ReadAll(w.Body)
			Expect(r).To(MatchJSON(`{"data": [], "links": {"self": "http://example.com/v1/payments/?count=10", "next": null}}`))
		})
		It("should default zero count to 10", func() {
			w := performRequest(db, "GET", "/v1/payments/?count=0")
			Expect(w.Code).To(Equal(http.StatusOK))
			r, _ := ioutil.ReadAll(w.Body)
			Expect(r).To(MatchJSON(`{"data": [], "links": {"self": "http://example.com/v1/payments/?count=10", "next": null}}`))
		})
		It("should default negative count size to 10", func() {
			w := performRequest(db, "GET", "/v1/payments/?count=-1")
			Expect(w.Code).To(Equal(http.StatusOK))
			r, _ := ioutil.ReadAll(w.Body)
			Expect(r).To(MatchJSON(`{"data": [], "links": {"self": "http://example.com/v1/payments/?count=10", "next": null}}`))
		})
		It("should list payments with updated count", func() {
			w := performRequest(db, "GET", "/v1/payments/?count=20")
			Expect(w.Code).To(Equal(http.StatusOK))
			r, _ := ioutil.ReadAll(w.Body)
			Expect(r).To(MatchJSON(`{"data": [], "links": {"self": "http://example.com/v1/payments/?count=20", "next": null}}`))
		})
		It("should list payments with updated after", func() {
			w := performRequest(db, "GET", "/v1/payments/?count=20&after=1")
			Expect(w.Code).To(Equal(http.StatusOK))
			r, _ := ioutil.ReadAll(w.Body)
			Expect(r).To(MatchJSON(`{"data": [], "links": {"self": "http://example.com/v1/payments/?count=20&after=1", "next": null}}`))
		})
		It("should list payments from database", func() {
			w := performRequest(mockDb{
				Payments: map[string]Payment{
					"1": {Id: "1"},
				},
			}, "GET", "/v1/payments/?count=10")
			Expect(w.Code).To(Equal(http.StatusOK))
			r, _ := ioutil.ReadAll(w.Body)
			Expect(r).To(MatchJSON(`{
					"data": [
						{"id": "1", "links": {"self": "http://example.com/v1/payments/1/"}}], 
					"links": {"self": "http://example.com/v1/payments/?count=10", "next": null}}`))
		})
		It("should have next link if there are too many items", func() {
			w := performRequest(mockDb{
				Payments: map[string]Payment{
					"a": {Id: "a"},
					"b": {Id: "b"},
				},
			}, "GET", "/v1/payments/?count=1")
			Expect(w.Code).To(Equal(http.StatusOK))
			r, _ := ioutil.ReadAll(w.Body)
			Expect(r).To(MatchJSON(`{
					"data": [
						{"id": "a", "links": {"self": "http://example.com/v1/payments/a/"}}], 
					"links": {"self": "http://example.com/v1/payments/?count=1", "next": "http://example.com/v1/payments/?count=1&after=a"}}`))
		})
	})
})
