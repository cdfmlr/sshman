package main

// This file contains examples of how to custom handlers aside from the
// out-of-the-box router.Crud. And with different parts of the crud package,
// we'll show you how curd package can help you to write your own codes
// at any level, from orm to services to controllers.

// Functions below are extremely vobose. And we are just showing how to
// use features of curd package. A real world implementation should
// be more concise.

import (
	"errors"
	"github.com/cdfmlr/crud/controller"
	"github.com/cdfmlr/crud/log"
	"github.com/cdfmlr/crud/orm"
	"github.com/cdfmlr/crud/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// UserGetSelf is an example to custom handler for
//    GET /users/:UserID.
//
// We hand-write it because we want to promise the user can only get
// his/her own info. Crud package does not provide this feature, but
// it does help u to make it.
//
// And this is also an example of how to use curd/controller, curd/service
// and curd/orm to handle problems at different levels.
func UserGetSelf(c *gin.Context) {
	// instead of using the global logger (defined in log.go), we can get
	// a local logger instance (that is, a entry of the global logger)
	// from crud/log.
	logger := log.ZoneLogger("sshman/controller")

	// user can only get his/her own info
	if !isUserSelf(c) {
		err := errors.New("you can only get your own information")
		logger.WithError(err).Error("mustUserSelf forbidden")
		controller.ResponseError(c, http.StatusForbidden, err)
		return
	}

	// U can use crud/controller to finish this job:
	if iPreferCrudController := false; iPreferCrudController {
		controller.GetByIDHandler[User]("UserID")(c)
	}
	// but, here, we decide to write codes:

	id := c.Param("UserID") // isUserSelf ensured that id is existent

	var user User
	var err error
	// crud/service defined lots of methods to do CRUD operations with orm.
	// The service.GetByID[T] is matching our need here. So, we just use it:
	if iPreferCrudService := true; iPreferCrudService {
		err = service.GetByID[User](c, id, &user)
	}
	// or, if you like, you can go even deeper, playing the crud/orm, which is
	// eventually the GROM:
	if iPreferWritingGORM := false; iPreferWritingGORM {
		// with orm.DB, u can even write raw SQL here, (see GORM docs for help:
		// https://gorm.io/docs/sql_builder.html), but it is not my style.
		result := orm.DB.WithContext(c).
			Where("id = ?", id).
			First(&user)
		err = result.Error
	}

	if err != nil {
		logger.WithContext(c).WithError(err).
			Warn("GetByIDHandler: getModelByID failed")
		controller.ResponseError(c, controller.CodeProcessFailed, err)
		return
	}

	controller.ResponseSuccess(c, user, gin.H{"message": "welcome to sshman"})
}

// UserGetSessions is another example to custom handler for
//    GET /users/:UserID/sessions
func UserGetSessions(c *gin.Context) {
	logger := log.ZoneLogger("sshman/controller")

	if !isUserSelf(c) {
		err := errors.New("you can only get your own sessions")
		logger.WithError(err).Error("mustUserSelf forbidden")
		controller.ResponseError(c, http.StatusForbidden, err)
		return
	}

	// Use a crud/controller helps us do remained works.
	controller.GetFieldHandler[User]("UserID", "sessions")(c)

}

// isUserSelf shows you how to get values from gin.Context.
func isUserSelf(c *gin.Context) bool {
	id, _ := strconv.ParseUint(c.Param("UserID"), 10, 32)
	userInfo := c.MustGet("userInfo").(*Payload)

	return userInfo.Id == uint(id)
}
