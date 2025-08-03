package users

type (
	UserResponse struct {
		UUID      string `json:"uuid"`
		Name      string `json:"name"`
		Username  string `json:"username"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
)
