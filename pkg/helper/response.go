package helper

import (
	"context"
	"crypto/sha512"
	"ecommerce/pkg/logger"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"ecommerce/pkg/constant"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

type ResponseWithPagination struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination interface{} `json:"pagination"`
}

func EncryptPassword(text string) string {
	passwordHash := sha512.Sum512([]byte(text))
	return hex.EncodeToString(passwordHash[:])
}

func ResponseError(ctx context.Context, err error) error {
	var (
		statusCode    int
		customError   string
		originalError string
	)

	statusCode, customError, originalError = TrimMesssage(err)
	route := ctx.Value(constant.FiberContext).(*fiber.Ctx).Path()
	body := ctx.Value(constant.FiberContext).(*fiber.Ctx).Body()
	method := ctx.Value(constant.FiberContext).(*fiber.Ctx).Method()
	var jsonBody map[string]interface{}
	if strings.ToLower(method) != "get" {
		if err := json.Unmarshal(body, &jsonBody); err != nil {
			log.Printf("Error unmarshaling body request to JSON: %v", err)
			jsonBody = map[string]interface{}{"body": string(body)}
		}
	}

	logger.LogError(constant.RESPONSE, originalError, jsonBody, route, method)

	response := Response{
		Message: customError,
		Data:    nil,
	}

	c := GetValueFiberFromContext(ctx)
	return c.Status(statusCode).JSON(response)
}

func ResponseOK(c *fiber.Ctx, msg string, data interface{}) error {
	route := c.Path()
	body := c.Body()
	method := c.Method()
	var jsonBody map[string]interface{}
	if strings.ToLower(method) != "get" {
		if err := json.Unmarshal(body, &jsonBody); err != nil {
			log.Printf("Error unmarshaling body request to JSON: %v", err)
			jsonBody = map[string]interface{}{"body": string(body)}
		}
	}
	logger.LogInfoWithData(constant.RESPONSE, jsonBody, method, route, msg)
	response := Response{
		Message: msg,
		Data:    data,
	}

	return c.Status(http.StatusOK).JSON(response)
}

func ResponseCreatedOK(c *fiber.Ctx, msg string, data interface{}) error {
	route := c.Path()
	body := c.Body()
	method := c.Method()
	var jsonBody map[string]interface{}
	if strings.ToLower(method) != "get" {
		if err := json.Unmarshal(body, &jsonBody); err != nil {
			log.Printf("Error unmarshaling body request to JSON: %v", err)
			jsonBody = map[string]interface{}{"body": string(body)}
		}
	}

	logger.LogInfoWithData(constant.RESPONSE, jsonBody, method, route, msg)
	response := Response{
		Message: msg,
		Data:    data,
	}

	return c.Status(http.StatusCreated).JSON(response)
}

func ResponseOkWithPagination(c *fiber.Ctx, msg string, data interface{}, pagination interface{}) error {
	route := c.Path()
	body := c.Body()
	method := c.Method()
	logger.LogInfoWithData(constant.RESPONSE, body, method, route, msg)
	response := ResponseWithPagination{
		Message:    msg,
		Data:       data,
		Pagination: pagination,
	}

	return c.Status(http.StatusOK).JSON(response)
}
