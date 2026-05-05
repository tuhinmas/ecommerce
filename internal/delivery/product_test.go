package delivery

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ecommerce/internal/entity"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductListAPI(t *testing.T) {
	app := fiber.New()

	// Test endpoint
	app.Get("/api/product", func(c *fiber.Ctx) error {
		// Simulate getting products
		products := []entity.GetProductListResponse{
			{
				Id:   "product-1",
				Name: "Laptop",
				Sku: []entity.Sku{
					{
						Id:      "sku-1",
						Variant: "Silver",
						Price:   1000000,
						Stock:   50,
						Uom:     "pcs",
						Image:   "laptop-silver.jpg",
					},
				},
			},
		}
		return c.JSON(fiber.Map{
			"status": "success",
			"data":   products,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/product", nil)
	resp, err := app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, "success", response["status"])
	assert.NotNil(t, response["data"])
}

func TestCreateProductAPI(t *testing.T) {
	app := fiber.New()

	app.Post("/api/product", func(c *fiber.Ctx) error {
		request := entity.CreateProductRequest{}
		if err := c.BodyParser(&request); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Validate required fields
		if request.ShopId == "" || request.Name == "" || len(request.Sku) == 0 {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": "shop_id, name, and sku are required",
			})
		}

		return c.Status(http.StatusCreated).JSON(fiber.Map{
			"status": "success",
		})
	})

	tests := []struct {
		name           string
		payload        entity.CreateProductRequest
		expectedStatus int
	}{
		{
			name: "Valid product creation",
			payload: entity.CreateProductRequest{
				ShopId: "shop-1",
				Name:   "Laptop",
				Sku: []entity.CreateSkuRequest{
					{
						Variant: "Silver",
						Price:   1000000,
						Uom:     "pcs",
						Image:   "laptop.jpg",
					},
				},
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Missing shop_id",
			payload: entity.CreateProductRequest{
				Name: "Laptop",
				Sku: []entity.CreateSkuRequest{
					{
						Variant: "Silver",
						Price:   1000000,
						Uom:     "pcs",
						Image:   "laptop.jpg",
					},
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payload, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/product", bytes.NewBuffer(payload))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
