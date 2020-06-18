package user

//easyjson:json
type User struct {
	Browsers  []string `json:"browsers,nocopy"`
	Email     string   `json:"email,nocopy"`
	Name      string   `json:"name,nocopy"`
	IsMSIE    bool     `json:"-"`
	IsAndroid bool     `json:"-"`
}
