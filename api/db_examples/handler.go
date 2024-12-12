package dbexample

import (
	"context"
	"database/sql"
	"time"

	DBpostgres "github.com/Sourceware-Lab/backend-proto/database/postgres"
)

//func PostRawSql(ctx context.Context, input *PostBodyInput) (*PostOutput, error) {
//	resp := &PostOutput{}
//	resp.Body.ID = "0"
//	return resp, nil
//}

func PostOrm(ctx context.Context, input *PostBodyInputDbExample) (*PostOutputDbExample, error) {
	resp := &PostOutputDbExample{}
	user := DBpostgres.User{
		Name: input.Body.Name,
		Age:  input.Body.Age,
	}
	if input.Body.Email != "" {
		user.Email = &input.Body.Email
	}
	if input.Body.Birthday != nil {
		birthday, err := time.Parse("2006-01-01", *input.Body.Birthday)
		if err != nil {
			return nil, err
		}
		user.Birthday = &birthday
	}
	if input.Body.MemberNumber != nil {
		user.MemberNumber = sql.NullString{
			String: *input.Body.MemberNumber,
			Valid:  true,
		}
	}
	DBpostgres.DB.Create(user)
	resp.Body.ID = "0"
	return resp, nil
}
