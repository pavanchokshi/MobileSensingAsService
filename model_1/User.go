package model

type (
	User struct {
		UserId   string `json:"userid,omitempty"`
		Username string `json:"username"`
		EmailId  string `json:"emailid"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	FormUserLogin struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
)
