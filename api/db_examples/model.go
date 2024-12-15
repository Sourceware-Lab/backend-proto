package dbexample

import (
	"time"

	"github.com/Sourceware-Lab/backend-proto/internal/utils"
)

type GetInputDbExample struct {
	ID string `doc:"Id for the user you want to get" example:"999" path:"id"`
}
type GetOutputDbExample struct {
	PostInputDbExample
}

type PostInputDbExample struct {
	Body PostBodyInputDbExampleBody `json:"body"`
}

func (p *PostInputDbExample) Format() *PostInputDbExample {
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

type PostBodyInputDbExampleBody struct {
	Name string `json:"name" doc:"Name for new user" example:"Jo" maxLength:"100" path:"name"`
	Age  uint8  `json:"age" doc:"Age for new user"  example:"25" path:"age"`

	// Optional
	Email        string  `json:"email" doc:"Email for new user"         example:"jo@example.com" maxLength:"100"      path:"email"     required:"false"`    //nolint:lll
	Birthday     *string `json:"birthday" doc:"Birthday for new user"      example:"2006-01-02"     format:"date"        path:"birthday"  required:"false"` //nolint:lll
	MemberNumber *string `json:"member_number" doc:"Member number for new user" example:"123456"         path:"member_number" required:"false"`             //nolint:lll
}

type PostOutputDbExample struct {
	Body struct {
		ID string `doc:"Id for new user" example:"999" json:"id"`
	}
}
