package models

type Collections string

const (
	UsersCollection     Collections = "users"
	CampaignsCollection Collections = "campaigns"
)

type User struct {
	ID       string `json:"id" bson:"_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone_number"`
	Address  string `json:"address"`
	Password string `json:"-"`
}
