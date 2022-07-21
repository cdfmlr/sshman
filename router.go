package main

// This file is a example of how to customize the router.
// The crud router is gin.Engine. So all your gin skills work here.

import (
	"github.com/cdfmlr/crud/router"
	"github.com/gin-gonic/gin"
)

// Router will be customized later with middlewares.
var Router *gin.Engine

var jwtAuth JwtAuthenticator

func initJwtAuth() {
	jwtAuth = NewHmacSecretAuthenticator(GlobalConfig.JwtSecret)
}

func InitRouters() {
	initJwtAuth()

	Router = router.NewRouter(
		router.WithMiddleware(JwtAuthMiddleware(jwtAuth)))
	// equivalent to:
	// Router = router.NewRouter()
	// Router.Use(JwtAuthMiddleware(jwtAuth))

	initAdminRouter()
	initUserRouter()

	logger.Info("routers are ready")
}

// InitAdminRouter initializes the routers for admin.
func initAdminRouter() {
	groupAdmin := Router.Group("/admin")
	groupAdmin.Use(JwtAuthMiddleware(jwtAuth))
	groupAdmin.Use(RoleAuthMiddleware(RoleAdmin))

	// use the out-of-the-box crud router.
	router.Crud[Host](groupAdmin, "/hosts")
	router.Crud[Session](groupAdmin, "/sessions")
	router.Crud[User](groupAdmin, "/users",
		router.CrudNested[User, Session]("sessions"))
}

// InitUserRouter initializes the routers for users.
func initUserRouter() {
	groupUser := Router.Group("/user")
	groupUser.Use(JwtAuthMiddleware(jwtAuth))
	groupUser.Use(RoleAuthMiddleware(RoleUser))

	// this is exactly what you do with the pure gin.
	groupUser.GET("/:UserID", UserGetSelf)
	groupUser.GET("/:UserID/sessions", UserGetSessions)
}
