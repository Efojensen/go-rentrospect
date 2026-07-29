package types

import "time"

type Client struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	ProfilePic  string `json:"profilePic"`
	PhoneNumber string `json:"phoneNumber"`
}

type Vendor struct {
	Client
	Calls      bool      `json:"calls"`
	NatId      string    `json:"natId"`
	Meetups    bool      `json:"meetups"`
	EndTime    time.Time `json:"endTime"`
	Location   string    `json:"location"`
	StartTime  time.Time `json:"startTime"`
	Deliveries bool      `json:"deliveries"`
}

type ClientBal struct {
	Name         string `json:"name"`
	TotalBal     int64  `json:"totalBal"`
	EscrowBal    int64  `json:"escrowBal"`
	AvailableBal int64  `json:"availableBal"`
}
