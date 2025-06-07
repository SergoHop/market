package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"market/internal/database"
	"market/internal/handlers"
	"market/internal/repository"
	"net/http"
	"time"
	_ "fmt"
	"html/template"
	
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	loc, err := time.LoadLocation("Europe/Moscow")
	if err == nil {
		time.Local = loc
	}

	productRepo := repository.NewProductRepository(db)

	productHandler := handlers.NewProductHandler(productRepo)

	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"div": func(a, b float64) float64 { return a / b },
        "mod": func(a, b int) int { return a % b },
	})

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/sell", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sell.html", nil)
	})
	r.POST("/sell", productHandler.HandleSell)

	r.GET("/buy", productHandler.BuyPage)

	r.GET("/after-sell", func(c *gin.Context) {
		c.HTML(http.StatusOK, "after_sell.html", nil)
	})

	port := ":8080"
	log.Printf("Server started on http://localhost%s", port)
	log.Fatal(r.Run(port))
}

