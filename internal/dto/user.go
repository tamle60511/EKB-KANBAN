package dto

type UserCreateReq struct {
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	Password     string `json:"password"`
	Email        string `json:"email"`
	DepartmentID int64  `json:"department_id,omitempty"`
}

type UserUpdateReq struct {
	Username     *string `json:"username"`
	FullName     *string `json:"full_name"`
	Password     *string `json:"password"`
	Email        *string `json:"email"`
	DepartmentID *int64  `json:"department_id"`
}

type UserDetailReq struct {
	UserID int64 `json:"id"`
}

type UserDetailRes struct {
	ID           int64  `json:"id"`
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	DepartmentID int64  `json:"department_id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}
