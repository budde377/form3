package main

import (
	"fmt"
	"github.com/google/logger"
	"io/ioutil"
)

func main() {
	defer logger.Init("Form3 API", true, false, ioutil.Discard).Close()
	c := OpenConfig()
	db := InitDb(c)
	router := Routes(c, db)
	if err := router.Run(fmt.Sprintf(":%d", c.Port)); err != nil {
		logger.Fatal("Failed to start server: ", err)
	}
}
