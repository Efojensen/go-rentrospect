package types

type Client struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	ProfilePic  string `json:"profilePic"`
	PhoneNumber string `json:"phoneNumber"`
}

type Vendor struct {
	Client
	NatId string `json:"natId"`
}
