package model

type UserType string

const (
	UserTypeGender UserType = "SUCCESS"
	UserTypeFuck   UserType = "aaaa"
)

// UserType2
// @Description 1(成功),2(失败)
type UserType2 int

const (
	UserType2Gender UserType2 = 1
	UserType2Fuck   UserType2 = 2
)

// Account model info
// @Description User account information
// @Description with user id and username
type Account struct {
	// userId
	ID int `json:"id" validate:"required,min=5" example:"1"`
	// user name
	Name string `json:"name" example:"account name"`

	T UserType `json:"userType"`

	T2 UserType2 `json:"userType2"`
}
