package main

// All you need is models. And this file defines them.

import "github.com/cdfmlr/crud/orm"

type Host struct {
	orm.BasicModel
	Hostname string `json:"hostname" gorm:"uniqueIndex;not null;size:256"`
	IP       string `json:"ip" gorm:"not null"`
	Port     string `json:"port" gorm:"not null"`
}

type Session struct {
	orm.BasicModel

	// [Why Host and HostID]
	// A GORM trick to make session links to one host,
	// while a host can be linked to multiple sessions.
	// I wonder if there are any elegant way to write
	// this one-to-many relationship in GROM, since it's
	// none of the belongs-to / has-one / has-many / many-to-many
	// associations introduced in GORM docs.
	Host   *Host `json:"host" gorm:"foreignKey:host_id;references:id"`
	HostID uint  `json:"host_id"` // the foreign key to Host

	Username   string `json:"username"`
	PrivateKey string `json:"private_key"`
}

type User struct {
	orm.BasicModel
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Role     int        `json:"role"`

	// A many-to-many relationship.
	Sessions []*Session `json:"sessions" gorm:"many2many:user_sessions"`
}

const (
	RoleUser  = 1
	RoleAdmin = 3
)
