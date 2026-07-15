package types

type Renter struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	ProfilePic  string `json:"profilePic"`
	PhoneNumber string `json:"phoneNumber"`
}

type Vendor struct {
	Renter
	NatId string `json:"natId"`
}
