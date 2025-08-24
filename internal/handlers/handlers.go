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
	renderPage(c, "index.html", nil)
}

func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func SellPage(c *gin.Context){
	renderPage(c, "sell.html", nil)
}

func (h *ProductHandler) HandleSell(c *gin.Context) {
    name := c.PostForm("name")
    priceStr := c.PostForm("price")
    description := c.PostForm("description")

    price, err := strconv.ParseFloat(priceStr, 64)
    if err != nil {
        renderError(c, "Некорректная цена")
        return
    }

    imagePath := handleImageUpload(c)

    product := models.Product{
        Name:        name,
        Price:       price,
        Description: description,
        ImagePath:   imagePath,
    }

    if err := h.repo.CreateProduct(&product); err != nil {
        renderError(c, "Не удалось сохранить товар")
        return
    }

    c.Redirect(http.StatusFound, "/after-sell")
}


func (h *ProductHandler) BuyPage(c *gin.Context) {
    // Получаем ВСЕ товары без фильтрации по времени
    products, err := h.repo.GetActiveProducts()
    if err != nil {
        log.Printf("Ошибка получения товаров: %v", err)
        renderError(c, "Не удалось загрузить товары")
        return
    }

    // Логируем количество полученных товаров
    log.Printf("Получено товаров: %d", len(products))

    c.HTML(http.StatusOK, "buy.html", gin.H{
        "Products": products,
        "Now":      time.Now(),
    })
}

func AfterSellPage(c *gin.Context) {
    renderPage(c, "after_sell.html", gin.H{
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

func (h *ProductHandler) HandleBuy(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        c.HTML(http.StatusBadRequest, "error.html", gin.H{
            "Error": "Некорректный ID товара",
        })
        return
    }

    if err := h.repo.MarkAsSold(id); err != nil {
        log.Printf("Ошибка при покупке товара ID=%d: %v", id, err)
        c.HTML(http.StatusInternalServerError, "error.html", gin.H{
            "Error": "Не удалось купить товар",
        })
        return
    }

    c.Redirect(http.StatusFound, "/buy")
}



//вспомагательная футкция для ошибки
func renderError(c *gin.Context, message string) {
    c.HTML(http.StatusInternalServerError, "error.html", gin.H{
        "Error": message,
    })
}
//обработка картинки 
func handleImageUpload(c *gin.Context) *string {
    file, header, err := c.Request.FormFile("image")
    if err != nil {
        return nil
    }
    defer file.Close()

    path := "uploads/" + header.Filename
    if err := saveUploadedFile(file, path); err != nil {
        return nil
    }
    return &path
}
//общ функция для рендера
func renderPage(c *gin.Context, templateName string, data gin.H) {
	c.HTML(http.StatusOK, templateName, data)
}



