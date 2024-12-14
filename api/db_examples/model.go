package dbexample

import (
	"time"

	DBpostgres "github.com/Sourceware-Lab/backend-proto/database/postgres"
)

type GetInputDbExample struct {
	ID string `path:"id" example:"999" doc:"Id for the user you want to get"`
}
type GetOutputDbExample struct {
	PostBodyInputDbExample
}

func (g *GetOutputDbExample) fromUserORM(user DBpostgres.User) *GetOutputDbExample {
	var memberNumber *string

	birthday := user.Birthday.Format(time.DateOnly)
	if user.MemberNumber.Valid {
		memberNumber = &user.MemberNumber.String
	}
	g.Body.Name = user.Name
	g.Body.Age = user.Age
	g.Body.Email = *user.Email
	g.Body.Birthday = &birthday
	g.Body.MemberNumber = memberNumber
	return g
}

type PostOutputDbExample struct {
	Body struct {
		ID string `json:"id" example:"999" doc:"Id for new user"`
	}
}

type PostBodyInputDbExample struct {
	Body struct {
		Name string `path:"name" maxLength:"100" example:"Jo" doc:"Name for new user"`
		Age  uint8  `path:"age" example:"25" doc:"Age for new user"`

		// Optional
		Email        string  `path:"email" maxLength:"100" example:"jo@example.com" doc:"Email for new user" required:"false"`
		Birthday     *string `path:"birthday" example:"2006-01-02" doc:"Birthday for new user" required:"false"`
		MemberNumber *string `path:"member_number" example:"123456" doc:"Member number for new user" required:"false"`
	}
}
