package main

import (
	"github.com/cdfmlr/crud/orm"
)

func InitDB() {
	orm.ConnectDB(
		orm.DBDriver(GlobalConfig.DB.Driver),
		GlobalConfig.DB.DSN)

	registerModels()

	logger.Info("db is ready")

	if GlobalConfig.CreateRootUser {
		createRootUser()
	}
}

func registerModels() {
	orm.RegisterModel(
		&User{},
		&Session{},
		&Host{},
	)
}

func createRootUser() {
	user := User{
		Name:  "root",
		Email: "root@sshman.example",
		Role:  RoleAdmin,
	}
	orm.DB.Create(&user)
	logger.Infof("root user created: id=%v", user.ID)
}
