package main

import (
    "log"
    "market/internal/handlers"
    "market/internal/repository"
    "market/internal/database"
	"html/template"
    "github.com/gin-gonic/gin"
)

func main() {
    db, err := database.ConnectDB() // подключение к базе
    if err != nil {
        log.Fatal("Ошибка подключения к БД:", err)
    }

    repo := repository.NewProductRepository(db)
    handler := handlers.NewProductHandler(repo)

    r := gin.Default()

	funcMap := template.FuncMap{
        "div": func(a, b float64) float64 { return a / b },
        "mod": func(a, b float64) float64 { return float64(int(a) % int(b)) },
    }

    r.SetFuncMap(funcMap)
    r.LoadHTMLGlob("templates/*")

    // Настройка шаблонов и статики
    r.LoadHTMLGlob("templates/*")
    r.Static("/static", "./static")
    r.Static("/uploads", "./uploads")

    // Маршруты
    r.GET("/", handlers.HomePage)
    r.GET("/sell", handlers.SellPage)
    r.POST("/sell", handler.HandleSell)
    r.GET("/buy", handler.BuyPage)
    r.GET("/after-sell", handlers.AfterSellPage)
    r.POST("/buy/:id", handler.HandleBuy)

    // Старт сервера
    if err := r.Run(":8080"); err != nil {
        log.Fatal("Ошибка запуска сервера:", err)
    }
}
