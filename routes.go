package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SelfLinks struct {
	Self string `json:"self"`
}

type PaymentSummaryRest struct {
	Id    string    `json:"id"`
	Links SelfLinks `json:"links"`
}

func summaryToRest(config *Config, summary PaymentSummary) PaymentSummaryRest {
	id := summary.Id
	return PaymentSummaryRest{
		Id: id,
		Links: SelfLinks{
			Self: fmt.Sprintf("%s/v1/payments/%s/", config.Host, id),
		},
	}

}

type PageLinks struct {
	Self string  `json:"self"`
	Next *string `json:"next"`
}

type PaymentsDataRest struct {
	Data  []PaymentSummaryRest `json:"data"`
	Links PageLinks            `json:"links"`
}

func listPaymentsEndpoint(config *Config, db Db) func(ctx *gin.Context) {
	return func(context *gin.Context) {
		// Extract size, default 10
		size := SafeStringToInt(context.Query("count"), 10)
		if size <= 0 {
			size = 10
		}
		// Extract after
		var after *string
		if v, ok := context.GetQuery("after"); ok {
			after = &v
		}

		// Fetch summaries
		summaries, err := db.GetPayments(size + 1, after)
		if err != nil {
			panic("failed to list payments")
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
			n := fmt.Sprintf("%s/v1/payments/?count=%d&after=%s", config.Host, size, mapped[len(mapped) - 1].Id)
			nextLink = &n
		}

		context.JSON(http.StatusOK, PaymentsDataRest{
			Data: mapped,
			Links: PageLinks{
				Self: selfLink,
				Next: nextLink,
			},
		})
	}
}

func v1(config *Config, db Db, router *gin.Engine) {
	v1 := router.Group("/v1")

	payments := v1.Group("/payments")
	{
		payments.GET("/", listPaymentsEndpoint(config, db))
	}

}

func Routes(config *Config, db Db) *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	v1(config, db, router)
	return router
}
