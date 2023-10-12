package server

import (
	"encoding/json"
	"strings"

	"tunes-service/cache"
	"tunes-service/data"
	"tunes-service/server/handlers"
	"tunes-service/server/middleware"

	"github.com/gofiber/fiber/v2"
)

func RunServer(db data.DB, adb data.AuthDB, cache cache.Cache) {
	app := fiber.New()

	app.Use("/api/spin", middleware.JWTMiddleware())

	app.Post("/api/register", func(c *fiber.Ctx) error {
		payload := struct {
			Username string
			Email    string
			Password string
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		if ok, _ := handlers.HandleRegistration(payload.Username, payload.Email, payload.Password, db); !ok {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/api/login", func(c *fiber.Ctx) error {
		payload := struct {
			Username string `json:"username,omitempty"`
			Email    string `json:"email,omitempty"`
			Password string
		}{}

		if err := c.BodyParser(&payload); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		usernameOrEmail := ""
		if payload.Username == "" && payload.Email == "" {
			return c.SendStatus(fiber.StatusBadRequest)
		} else if payload.Username != "" {
			usernameOrEmail = payload.Username
		} else {
			usernameOrEmail = payload.Email
		}

		if at, rt, err := handlers.HandleLogin(usernameOrEmail, payload.Password, db, adb); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		} else {
			c.Cookie(&fiber.Cookie{
				Name:     "tunes-refresh-token",
				Value:    rt,
				Path:     "/api/refresh",
				Secure:   true,
				HTTPOnly: true,
				SameSite: "strict",
			})

			j, _ := json.Marshal(struct {
				Token string
			}{
				at,
			})
			c.Status(fiber.StatusOK).SendString(string(j))
		}
		return c.SendStatus(fiber.StatusNotFound)
	})

	app.Get("/api/refresh", func(c *fiber.Ctx) error {
		at := c.GetReqHeaders()[fiber.HeaderAuthorization]
		rt := c.Cookies("tunes-refresh-token")
		if at == "" || rt == "" {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		if newAccessToken, err := handlers.HandleRefresh(strings.Replace(at, "Bearer ", "", 1), rt, adb); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		} else {
			j, _ := json.Marshal(struct {
				Token string
			}{
				newAccessToken,
			})
			c.Status(fiber.StatusOK).SendString(string(j))
		}
		return c.SendStatus(fiber.StatusNotFound)
	})

	app.Get("/api/health", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	app.Post("/api/spin", func(c *fiber.Ctx) error {
		req := handlers.SpinRequest{}

		err := c.BodyParser(&req)
		if err != nil {
			return err
		}

		handlers.HandleSpin(req, db, cache)
		return c.SendStatus(fiber.StatusOK)
	})

	app.Listen(":8080")
}
