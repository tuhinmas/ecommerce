package helper

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"

	"ecommerce/internal/entity"
	"ecommerce/pkg/constant"
)

func CreateContext() context.Context {
	return context.Background()
}

func CreateContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

func CreateContextWithCustomTimeout(timeout int) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
}

func SetValueToContext(ctx context.Context, c *fiber.Ctx) context.Context {
	var userId, warehouseId string

	userId, ok := c.Locals("user-id").(string)
	if !ok {
		userId = "0"
	}
	warehouseId, ok = c.Locals("warehouse-id").(string)
	if !ok {
		warehouseId = "0"
	}

	ctx = context.WithValue(ctx, constant.HeaderContext, entity.ValueContext{
		UserId:      userId,
		WarehouseId: warehouseId,
	})

	return context.WithValue(ctx, constant.FiberContext, c)
}

func GetValueContext(ctx context.Context) (valueCtx entity.ValueContext) {
	return ctx.Value(constant.HeaderContext).(entity.ValueContext)
}

func GetValueFiberFromContext(ctx context.Context) *fiber.Ctx {
	return ctx.Value(constant.FiberContext).(*fiber.Ctx)
}

func SetContext(ctx context.Context, c *fiber.Ctx) context.Context {
	return context.WithValue(ctx, constant.FiberContext, c)
}
