package model

type UserProfile struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Bio      string `json:"bio"`
}
