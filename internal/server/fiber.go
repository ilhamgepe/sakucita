package server

import (
	"errors"
	"fmt"

	"sakucita/internal/domain"
	"sakucita/internal/dto"
	"sakucita/internal/shared/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func fiberErrorHandler(c fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	message := any("internal server error")

	// error validation
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		status = fiber.StatusUnprocessableEntity
		message = utils.GenerateMessageValidation(err)

		return c.Status(status).JSON(dto.ErrorResponse{
			Message: fiber.ErrUnprocessableEntity.Message,
			Errors:  message,
		})
	}

	// handle errro dari domain
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		// fmt.Printf("%+v\n", appErr.Err)
		// fmt.Printf("%+v\n", appErr.Message)
		// fmt.Printf("%+v\n", appErr.Error())
		// fmt.Printf("%+v\n", appErr.Err.Error())
		// fmt.Printf("%+v\n", appErr.Code)
		return c.Status(appErr.Code).JSON(dto.ErrorResponse{
			Message: appErr.Message,
			Errors:  appErr.Err.Error(),
		})
	}

	// fiber error
	var fe *fiber.Error
	if errors.As(err, &fe) {
		status = fe.Code
		message = fe.Message
		return c.Status(status).JSON(dto.ErrorResponse{
			Errors: message,
		})
	}

	fmt.Printf("error fallback: %v", err)
	// fallback
	return c.Status(status).JSON(dto.ErrorResponse{
		Message: fiber.ErrInternalServerError.Message,
		Errors:  message,
	})
}
