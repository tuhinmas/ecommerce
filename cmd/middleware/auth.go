package middleware

import (
	"errors"
	"net/http"
	"time"

	"ecommerce/pkg/config"
	"ecommerce/pkg/constant"
	"ecommerce/pkg/helper"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

// Define the claims structure for the JWT
type Claims struct {
	Id          string `json:"id"`
	WarehouseId string `json:"warehouse_id"`
	jwt.StandardClaims
}

// Define a function for generating a new JWT
func GenerateToken(id string, warehouseId string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(24))
	claims := &Claims{
		Id:          id,
		WarehouseId: warehouseId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateTokenAdmin(id string) (string, error) {
	expirationTime := time.Now().Add(time.Hour * time.Duration(24))
	claims := &Claims{
		Id: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(config.GetEnv("JWT_SECRET_ADMIN")))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Define a middleware for verifying JWT authentication and expiration
func AuthUser(c *fiber.Ctx) error {
	ctx := helper.CreateContext()
	ctx = helper.SetContext(ctx, c)

	// Get the Authorization header from the request
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		err := helper.Error(http.StatusUnauthorized, constant.MsgAuthorizationHeaderNotFound, errors.New(constant.MsgAuthorizationHeaderNotFound))
		return helper.ResponseError(ctx, err)
	}

	// Verify that the Authorization header starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		err := helper.Error(http.StatusUnauthorized, constant.MsgInvalidFormatAuthorization, errors.New(constant.MsgInvalidFormatAuthorization))
		return helper.ResponseError(ctx, err)
	}

	// Parse the JWT from the Authorization header
	tokenString := authHeader[len("Bearer "):]
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			err = helper.Error(http.StatusUnauthorized, constant.MsgInvalidSignature, err)
			return helper.ResponseError(ctx, err)
		}
		err = helper.Error(http.StatusUnauthorized, constant.MsgInvalidToken, err)
		return helper.ResponseError(ctx, err)
	}

	// Check if the JWT has expired
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		err = helper.Error(http.StatusUnauthorized, constant.MsgInvalidToken, errors.New(constant.MsgInvalidToken))
		return helper.ResponseError(ctx, err)
	}
	if claims.ExpiresAt < time.Now().Unix() {
		err = helper.Error(http.StatusUnauthorized, constant.MsgExpiredToken, errors.New(constant.MsgExpiredToken))
		return helper.ResponseError(ctx, err)
	}

	// Set the user ID in the context for future requests
	c.Locals("warehouse-id", claims.WarehouseId)
	c.Locals("user-id", claims.Id)

	// Call the next middleware in the chain
	return c.Next()
}

func AuthAdmin(c *fiber.Ctx) error {
	ctx := helper.CreateContext()
	ctx = helper.SetContext(ctx, c)

	// Get the Authorization header from the request
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		err := helper.Error(http.StatusUnauthorized, constant.MsgAuthorizationHeaderNotFound, errors.New(constant.MsgAuthorizationHeaderNotFound))
		return helper.ResponseError(ctx, err)
	}

	// Verify that the Authorization header starts with "Bearer "
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		err := helper.Error(http.StatusUnauthorized, constant.MsgInvalidFormatAuthorization, errors.New(constant.MsgInvalidFormatAuthorization))
		return helper.ResponseError(ctx, err)
	}

	// Parse the JWT from the Authorization header
	tokenString := authHeader[len("Bearer "):]
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetEnv("JWT_SECRET_ADMIN")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			err = helper.Error(http.StatusUnauthorized, constant.MsgInvalidSignature, errors.New(constant.MsgInvalidSignature))
			return helper.ResponseError(ctx, err)
		}
		err = helper.Error(http.StatusUnauthorized, constant.MsgInvalidToken, errors.New(constant.MsgInvalidToken))
		return helper.ResponseError(ctx, err)
	}

	// Check if the JWT has expired
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		err = helper.Error(http.StatusUnauthorized, constant.MsgInvalidToken, errors.New(constant.MsgInvalidToken))
		return helper.ResponseError(ctx, err)
	}
	if claims.ExpiresAt < time.Now().Unix() {
		err = helper.Error(http.StatusUnauthorized, constant.MsgExpiredToken, errors.New(constant.MsgExpiredToken))
		return helper.ResponseError(ctx, err)
	}

	c.Locals("user-id", claims.Id)

	// Call the next middleware in the chain
	return c.Next()
}
