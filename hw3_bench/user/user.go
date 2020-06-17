package user

//easyjson:json
type User struct {
	Browsers  []string `json:"browsers"`
	Email     string   `json:"email"`
	Name      string   `json:"name"`
	IsMSIE    bool     `json:"-"`
	IsAndroid bool     `json:"-"`
}
