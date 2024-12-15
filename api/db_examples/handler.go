package dbexample

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	DBpostgres "github.com/Sourceware-Lab/backend-proto/database/postgres"
)

func GetRawSql(ctx context.Context, input *GetInputDbExample) (resp *GetOutputDbExample, err error) {
	resp = &GetOutputDbExample{}
	id, err := strconv.Atoi(input.ID)

	if err != nil {
		log.Error().Err(err).Msg("Error parsing ID")
		return nil, err
	}

	DBpostgres.DB.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&resp.Body)
	resp.Format()
	return
}

//	func PostRawSql(ctx context.Context, input *PostBodyInput) (*PostOutput, error) {
//		resp := &PostOutput{}
//		resp.Body.ID = "0"
//		return resp, nil
//	}
func GetOrm(ctx context.Context, input *GetInputDbExample) (resp *GetOutputDbExample, err error) {
	resp = &GetOutputDbExample{}
	id, err := strconv.Atoi(input.ID)

	if err != nil {
		log.Error().Err(err).Msg("Error parsing ID")
		return nil, err
	}

	result := DBpostgres.DB.Model(DBpostgres.User{}).Where(DBpostgres.User{ID: uint(id)}).First(&resp.Body)
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error getting user")
		return nil, result.Error
	}
	resp.Format()
	return
}

func PostOrm(ctx context.Context, input *PostBodyInputDbExample) (resp *PostOutputDbExample, err error) {
	resp = &PostOutputDbExample{}
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
		birthday, err := time.Parse(time.DateOnly, *input.Body.Birthday)
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
	return
}
