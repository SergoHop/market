package handlers_test

import (
	"html/template"
	"market/internal/handlers"
	"market/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type fakeRepo struct{}

func (f *fakeRepo) CreateProduct(product *models.Product) error {
	return nil
}

func (f *fakeRepo) GetActiveProducts() ([]models.Product, error) {
	return []models.Product{
		{
			ID:        1,
			Name:      "Товар 1",
			Price:     100,
			ExpiresAt: time.Now().Add(1 * time.Hour),
		},
	}, nil
}

func TestBuyPage(t *testing.T){
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	tmpl := template.Must(template.New("buy.html").Parse(`
		<html>
		<body>
		{{range .Products}}
		<p>{{.Name}}</p>
		{{end}}
		</body>
		</html>`))
	router.SetHTMLTemplate(tmpl)
	h := handlers.NewProductHandler(&fakeRepo{})
	router.GET("/buy", h.BuyPage)

	req, _ := http.NewRequest("GET", "/buy", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w ,req)

	assert.Equal(t,http.StatusOK, w.Code)

	assert.Contains(t, w.Body.String(), "Товар 1")
}





