package dto

import "time"

type AdminInfoOutput struct {
	ID           int       `json:"id"`
	UserName     string    `json:"name"`
	LoginTime    time.Time `json:"login_time"`
	Avatar       string    `json:"avatar"`
	Introduction string    `json:"introduction"`
	Roles        []string  `json:"roles"`
}

type ChangePasswordInput struct {
	NewPassword string `json:"new_password" form:"new_password" validate:"required"`
}
