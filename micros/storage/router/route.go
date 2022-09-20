// Copyright (c) 2021 Amirhossein Movahedi (@qolzam)
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/middleware/authcookie"
	appConfig "github.com/red-gold/telar-web/micros/storage/config"
	"github.com/red-gold/telar-web/micros/storage/handlers"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App) {

	if appConfig.StorageConfig.ProxyBalancer != "" {
		app.Use(proxy.Balancer(proxy.Config{
			Servers: []string{
				appConfig.StorageConfig.ProxyBalancer,
			},
		}))
	}

	// Middleware
	authCookieMiddleware := authcookie.New(authcookie.Config{
		JWTSecretKey: []byte(*config.AppConfig.PublicKey),
	})

	// Router
	app.Post("/:uid/:dir", authCookieMiddleware, handlers.UploadeHandle)
	app.Get("/:uid/:dir/:name", authCookieMiddleware, handlers.GetFileHandle)

}
