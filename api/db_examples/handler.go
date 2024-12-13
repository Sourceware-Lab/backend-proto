package dbexample

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

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
		Name:         input.Body.Name,
		Email:        nil,
		Birthday:     nil,
		MemberNumber: sql.NullString{},
		ActivatedAt:  sql.NullTime{},
		Age:          input.Body.Age,
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
	result := DBpostgres.DB.Create(&user) // NOTE. This is a POINTER!
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error creating user")
		return nil, result.Error
	}
	resp.Body.ID = strconv.Itoa(int(user.ID))
	return resp, nil
}
