package middleware

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/mstrYoda/maasanketi.co/pkg/logger"
	"github.com/mstrYoda/maasanketi.co/pkg/rescode"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusBadRequest
	if res, ok := err.(*rescode.RC); ok {
		log.Logger().Error(res.OriginalError().Error())
		return c.Status(res.HttpCode).JSON(res.JSON(res.Message))
	} else {
		log.Logger().Error(err.Error())
	}
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.SendStatus(code)
}
