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
	Body PostBodyInputDbExampleBody
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
	Name string `doc:"Name for new user" example:"Jo" maxLength:"100" path:"name"`
	Age  uint8  `doc:"Age for new user"  example:"25" path:"age"`

	// Optional
	Email        string  `doc:"Email for new user"         example:"jo@example.com" maxLength:"100"      path:"email"     required:"false"`
	Birthday     *string `doc:"Birthday for new user"      example:"2006-01-02"     format:"date"        path:"birthday"  required:"false"`
	MemberNumber *string `doc:"Member number for new user" example:"123456"         path:"member_number" required:"false"`
}

type PostOutputDbExample struct {
	Body struct {
		ID string `doc:"Id for new user" example:"999" json:"id"`
	}
}
