package dtos

type UserRequest struct {
	RoleID int64  `json:"role_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}
