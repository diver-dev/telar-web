// Copyright (c) 2021 Amirhossein Movahedi (@qolzam)
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/middleware/authcookie"
	"github.com/red-gold/telar-core/middleware/authhmac"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-web/micros/actions/handlers"
)

// @title Actions micro API
// @version 1.0
// @description This is an API to handle web socket server actions dispatch
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email dev@telar.dev
// @license.name MIT
// @license.url https://github.com/red-gold/telar-web/blob/master/LICENSE
// @host social.faas.telar.dev
// @BasePath /actions
func SetupRoutes(app *fiber.App) {

	// Middleware
	authHMACMiddleware := func(hmacWithCookie bool) func(*fiber.Ctx) error {
		var Next func(c *fiber.Ctx) bool
		if hmacWithCookie {
			Next = func(c *fiber.Ctx) bool {
				if c.Get(types.HeaderHMACAuthenticate) != "" {
					return false
				}
				return true
			}
		}
		return authhmac.New(authhmac.Config{
			Next:          Next,
			PayloadSecret: *config.AppConfig.PayloadSecret,
		})
	}

	authCookieMiddleware := func(hmacWithCookie bool) func(*fiber.Ctx) error {
		var Next func(c *fiber.Ctx) bool
		if hmacWithCookie {
			Next = func(c *fiber.Ctx) bool {
				if c.Get(types.HeaderHMACAuthenticate) != "" {
					return true
				}
				return false
			}
		}
		return authcookie.New(authcookie.Config{
			Next:         Next,
			JWTSecretKey: []byte(*config.AppConfig.PublicKey),
		})
	}

	hmacCookieHandlers := []func(*fiber.Ctx) error{authHMACMiddleware(true), authCookieMiddleware(true)}

	// Router
	app.Post("/room", authHMACMiddleware(false), handlers.CreateActionRoomHandle)
	app.Post("/dispatch/:roomId", authHMACMiddleware(false), handlers.DispatchHandle)
	app.Put("/room", append(hmacCookieHandlers, handlers.UpdateActionRoomHandle)...)
	app.Put("/room/access-key", append(hmacCookieHandlers, handlers.SetAccessKeyHandle)...)
	app.Delete("/room/:roomId", authHMACMiddleware(false), handlers.DeleteActionRoomHandle)
	app.Get("/room/access-key", append(hmacCookieHandlers, handlers.GetAccessKeyHandle)...)
	app.Post("/room/verify", append(hmacCookieHandlers, handlers.VerifyAccessKeyHandle)...)
}
