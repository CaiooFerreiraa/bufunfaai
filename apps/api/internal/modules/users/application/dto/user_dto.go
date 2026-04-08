package dto

type UpdateCurrentUserRequest struct {
	FullName string `json:"full_name" validate:"required,min=3,max=120"`
	Phone    string `json:"phone" validate:"omitempty,min=8,max=20"`
}

type UserOutput struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone,omitempty"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
