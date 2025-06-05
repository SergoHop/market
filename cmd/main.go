package main

import (
	"log"
	"net/http"
	
	"github.com/gin-gonic/gin"
)

func main(){
	r := gin.Default

	r.GET("/" , func(c *gin.Context){
		c.HTML(http.StatusOK, "index.html", nil)
	})

	log.Fatal(r.Run(":8080"))
}