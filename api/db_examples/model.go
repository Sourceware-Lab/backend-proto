package dbexample

import (
	"time"

	"github.com/Sourceware-Lab/go-huma-gin-postgres-template/internal/utils"
)

type GetInputDBExample struct {
	ID string `doc:"Id for the user you want to get" example:"999" path:"id"`
}
type GetOutputDBExample struct {
	PostInputDBExample
}

type PostInputDBExample struct {
	Body PostBodyInputDBExampleBody `json:"body"`
}

func (p *PostInputDBExample) Format() *PostInputDBExample {
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

type PostBodyInputDBExampleBody struct {
	Name string `doc:"Name for new user" example:"Jo" json:"name" maxLength:"100" path:"name"`
	Age  uint8  `doc:"Age for new user"  example:"25" json:"age"  path:"age"`

	// Optional
	Email        string  `doc:"Email for new user"         example:"jo@example.com" json:"email"        maxLength:"100"     path:"email"     required:"false"` //nolint:lll
	Birthday     *string `doc:"Birthday for new user"      example:"2006-01-02"     format:"date"       json:"birthday"     path:"birthday"  required:"false"` //nolint:lll
	MemberNumber *string `doc:"Member number for new user" example:"123456"         json:"memberNumber" path:"memberNumber" required:"false"`                  //nolint:lll
}

type PostOutputDBExample struct {
	Body struct {
		ID string `doc:"Id for new user" example:"999" json:"id"`
	}
}
