package dbexample

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	DBpostgres "github.com/Sourceware-Lab/go-huma-gin-postgres-template/database/postgres"
)

func GetRawSQL(_ context.Context, input *GetInputDBExample) (*GetOutputDBExample, error) {
	resp := &GetOutputDBExample{}

	id, err := strconv.Atoi(input.ID)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing ID")

		return nil, fmt.Errorf("error parsing id: %w", err)
	}

	DBpostgres.DB.Raw("SELECT * FROM users WHERE id = ?", id).Scan(&resp.Body)
	resp.Format()

	return resp, nil
}

func PostRawSQL(_ context.Context, input *PostInputDBExample) (*PostOutputDBExample, error) {
	resp := &PostOutputDBExample{}

	var email *string
	if input.Body.Email != "" {
		email = &input.Body.Email
	}

	var birthday *time.Time

	if input.Body.Birthday != nil {
		parsedBirthday, err := time.Parse(time.DateOnly, *input.Body.Birthday)
		if err != nil {
			log.Error().Err(err).Msg("Error parsing birthday")

			return nil, fmt.Errorf("error parsing birthday: %w", err)
		}

		birthday = &parsedBirthday
	}

	memberNumber := sql.NullString{}
	if input.Body.MemberNumber != nil {
		memberNumber = sql.NullString{
			String: *input.Body.MemberNumber,
			Valid:  true,
		}
	}

	activatedAt := sql.NullTime{}

	result := DBpostgres.DB.Raw(`
        INSERT INTO users (name, email, birthday, member_number, activated_at, age)
        VALUES (?, ?, ?, ?, ?, ?)
        RETURNING id`,
		input.Body.Name, email, birthday, memberNumber.String, activatedAt.Time, input.Body.Age).Scan(&resp.Body.ID)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error inserting user")

		return nil, result.Error
	}

	return resp, nil
}

func GetOrm(_ context.Context, input *GetInputDBExample) (*GetOutputDBExample, error) {
	resp := &GetOutputDBExample{}

	id, err := strconv.Atoi(input.ID)
	if err != nil {
		log.Error().Err(err).Msg("Error parsing ID")

		return nil, fmt.Errorf("error parsing id: %w", err)
	}

	result := DBpostgres.DB.Model(DBpostgres.User{}).Where(DBpostgres.User{ID: uint(id)}).First(&resp.Body) //nolint:gosec
	if result.Error != nil {
		log.Error().Err(result.Error).Msg("Error getting user")

		return nil, result.Error
	}

	resp.Format()

	return resp, nil
}

func PostOrm(_ context.Context, input *PostInputDBExample) (*PostOutputDBExample, error) {
	resp := &PostOutputDBExample{}
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
			return nil, fmt.Errorf("error parsing birthday: %w", err)
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

	resp.Body.ID = strconv.Itoa(int(user.ID)) //nolint:gosec

	return resp, nil
}
