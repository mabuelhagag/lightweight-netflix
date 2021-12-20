package user

// User struct
type User struct {
	Name  string `json:"name" bson:"name"`
	Pages uint   `json:"pages" bson:"page_count"`
}
