package handlers

import (
	"net/http"
	"market/internal/repository" 
	_ "market/internal/database"
	"github.com/gin-gonic/gin"
	"strconv"
	"log"
	"market/internal/models"
    "time"
	"os"
	"io"
	"mime/multipart"
    _ "github.com/jmoiron/sqlx"
)

type ProductHandler struct {
	repo repository.ProductRepository
}

func HomePage(c *gin.Context){
	c.HTML(http.StatusOK, "index.html", nil)
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func SellPage(c *gin.Context){
	c.HTML(http.StatusOK, "sell.html",nil)
}

func (h *ProductHandler) HandleSell(c *gin.Context) {
	name := c.PostForm("name")
	priceStr := c.PostForm("price")
	description := c.PostForm("description")
	
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Error": "Некорректная цена",
		})
		return
	}
    var imagePath *string
    file, header, err := c.Request.FormFile("image")
    if err == nil {
        defer file.Close()
        path := "uploads/" + header.Filename
        if err := saveUploadedFile(file, path); err == nil {
            imagePath = &path
        }
    }

	product := models.Product{
		Name:        name,
		Price:       price,
		Description: description,
        ImagePath:   imagePath,
	}

	if err := h.repo.CreateProduct(&product); err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", gin.H{
			"Error": "Не удалось сохранить товар",
		})
		return
	}

	c.Redirect(http.StatusFound, "/after-sell")
}

func (h *ProductHandler) BuyPage(c *gin.Context) {
    products, err := h.repo.GetActiveProducts()
    if err != nil {
        log.Printf("DB error: %v", err)
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "Error": "Не удалось загрузить список товаров",
        })
        return
    }

    log.Printf("Loaded %d products", len(products))
    
    c.HTML(http.StatusOK, "buy.html", gin.H{
        "Products": products,
        "Now":      time.Now(),
    })
}

func AfterSellPage(c *gin.Context) {
    c.HTML(http.StatusOK, "after_sell.html", gin.H{
        "Title": "Что вы хотите сделать?",
    })
}

func saveUploadedFile(file multipart.File, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	return err
}

