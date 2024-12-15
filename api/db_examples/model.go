package dbexample

import (
	"time"

	"github.com/Sourceware-Lab/backend-proto/internal/utils"
)

type GetInputDbExample struct {
	ID string `path:"id" example:"999" doc:"Id for the user you want to get"`
}
type GetOutputDbExample struct {
	PostBodyInputDbExample
}

type PostOutputDbExample struct {
	Body struct {
		ID string `json:"id" example:"999" doc:"Id for new user"`
	}
}

type PostBodyInputDbExample struct {
	Body PostBodyInputDbExampleBody
}
type PostBodyInputDbExampleBody struct {
	Name string `path:"name" maxLength:"100" example:"Jo" doc:"Name for new user"`
	Age  uint8  `path:"age" example:"25" doc:"Age for new user"`

	// Optional
	Email        string  `path:"email" maxLength:"100" example:"jo@example.com" doc:"Email for new user" required:"false"`
	Birthday     *string `path:"birthday" example:"2006-01-02" doc:"Birthday for new user" required:"false" format:"date"`
	MemberNumber *string `path:"member_number" example:"123456" doc:"Member number for new user" required:"false"`
}

func (p *PostBodyInputDbExample) Format() *PostBodyInputDbExample {
	if p.Body.Birthday != nil {
		birthday, err := utils.ParseAnyDatetime(*p.Body.Birthday)
		if err != nil {
			return p
		}
		marshaledBirthday := birthday.Format(time.DateOnly)

		p.Body.Birthday = &marshaledBirthday
	}
	return p
}
