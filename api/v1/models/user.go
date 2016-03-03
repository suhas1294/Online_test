package models

import (
	"time"
)

// Registration struct [account/sign_up]
type Register struct {
	Id                    int    `valid:"numeric"`
	Firstname             string `valid:"alphanum,required"`
	Lastname              string `valid:"alphanum,required"`
	Email                 string `valid:"email,required"`
	Password              string `valid:"alphanum,required"`
	Password_confirmation string `valid:"alphanum,required"`
	Branch                string `valid:"alphanum",required`
	Year_of_passing				string `valid:"alphanum",required`
	Phone_number					string `valid:"alphanum,required"`
	Auth_token          	string `valid:"alphanum,required"`
}


type ErrorMessage struct {
	Success string
	Error   string
}

type SignIn struct {
	Success string
	Message string
	User    Register
	Session Session
}
// Sign_up struct end




type UserDetails struct {
	Id                 int
	Firstname          string
	Lastname           string
	Email              string
	User_thumbnail     string
	User_thumbnail_web string
}

type Notification struct {
	SenderId   int `valid:"numeric,required"`
	RecieverId int `valid:"numeric,required"`
}

type InviteEmail struct {
	SenderId      int    `valid:"numeric,required"`
	RecieverEmail string `valid:"email,required"`
}

// Session struct [account/session]
type Session struct {
	SessionId int
	StartTime time.Time
}


// Message struct [controllers/account]
// Common for sign_up, session and password
type Message struct {
	Success string
	Message string
	User    Register
}


type EmailMessage struct {
	Success string
	Message string
	User    InviteEmail
}

type EmailErrorMessage struct {
	Success string
	Error   string
}

// User profile Struct
type Profile struct {
	Id                    int    `valid:"numeric,required"`
	Firstname             string `valid:"alphanum,required"`
	Lastname              string `valid:"alphanum,required"`
	Password              string `valid:"alphanum,required"`
	Password_confirmation string `valid:"alphanum,required"`
	City                  string `valid:"alphanum"`
	State                 string `valid:"alphanum"`
	Country               string `valid:"alphanum"`
	User_thumbnail        string `valid:"alphanum"`
	Description           string `valid:"alphanum"`
}

type ProfileMessage struct {
	Success string
	Message string
	User    Profile
}

type ProfileErrorMessage struct {
	Success string
	Error   string
}

type UserListMessage struct {
	Success  string
	Message  string
	User_ids []int
}

type UserList struct {
	Success      string
	No_Of_Users  int
	User_Details []UserDetails
}