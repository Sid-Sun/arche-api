package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/nsnikhil/erx"
	"github.com/sid-sun/arche-api/app/custom_errors"
	"github.com/sid-sun/arche-api/app/types"
	"go.uber.org/zap"
)

type UsersTable interface {
	Get(emailID string) (types.User, *erx.Erx)
	GetVerificationStatus(emailID string) (bool, string, *erx.Erx)
	UpdateVerificationToken(emailID string, vetkn string) *erx.Erx
	Create(emailID string, encryptionKey string, keyHash string, vetkn string) (types.UserID, *erx.Erx)
	VerifyUser(vetkn string) *erx.Erx
}

type users struct {
	lgr *zap.Logger
	db  *sql.DB
}

func (u *users) VerifyUser(vetkn string) *erx.Erx {
	query := `UPDATE users SET verified = 1 WHERE verification_key=@veKey;`

	res, err := u.db.Exec(query, sql.Named("veKey", vetkn))
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			u.lgr.Error(fmt.Sprintf("[Database] [Users] [VerifyUser] [Scan] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return errx
		}
		if errors.Is(err, sql.ErrNoRows) {
			errx = erx.WithArgs(errx, erx.SeverityInfo, custom_errors.NoRowsInResultSet)
			u.lgr.Info(fmt.Sprintf("[Database] [Users] [VerifyUser] [Scan] [ErrSQLNoResultsInSet] %s", errx.String()))
			return errx
		}
		u.lgr.Debug(fmt.Sprintf("[Database] [Users] [VerifyUser] [Scan] %s", errx.Error()))
		return errx
	}

	var count int64
	if count, err = res.RowsAffected(); err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			u.lgr.Error(fmt.Sprintf("[Database] [Users] [VerifyUser] [RowsAffected] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return errx
		}
		u.lgr.Debug(fmt.Sprintf("[Database] [Users] [VerifyUser] [RowsAffected] %s", err.Error()))
		return errx
	}

	if count == 0 {
		return erx.WithArgs(custom_errors.NoRowsAffected, erx.SeverityInfo)
	}

	return nil
}

func (u *users) UpdateVerificationToken(emailID string, vetkn string) *erx.Erx {
	query := `UPDATE users SET verification_key = @veKey WHERE email=@email`

	_, err := u.db.Exec(query, sql.Named("email", emailID), sql.Named("veKey", vetkn))
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			u.lgr.Error(fmt.Sprintf("[Database] [Users] [UpdateVerificationToken] [Scan] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return errx
		}
		if errors.Is(err, sql.ErrNoRows) {
			errx = erx.WithArgs(errx, erx.SeverityInfo, custom_errors.NoRowsInResultSet)
			u.lgr.Info(fmt.Sprintf("[Database] [Users] [UpdateVerificationToken] [Scan] [ErrSQLNoResultsInSet] %s", errx.String()))
			return errx
		}
		u.lgr.Debug(fmt.Sprintf("[Database] [Users] [UpdateVerificationToken] [Scan] %s", errx.Error()))
		return errx
	}

	return nil
}

func (u *users) GetVerificationStatus(emailID string) (bool, string, *erx.Erx) {
	query := `SELECT verified, verification_key FROM users WHERE email=@email;`

	var verified bool
	var verificationKey sql.NullString

	row := u.db.QueryRow(query, sql.Named("email", emailID))
	err := row.Scan(&verified, &verificationKey)
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			u.lgr.Error(fmt.Sprintf("[Database] [Users] [GetVerificationStatus] [Scan] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return false, "", errx
		}
		if errors.Is(err, sql.ErrNoRows) {
			errx = erx.WithArgs(errx, erx.SeverityInfo, custom_errors.NoRowsInResultSet)
			u.lgr.Info(fmt.Sprintf("[Database] [Users] [GetVerificationStatus] [Scan] [ErrSQLNoResultsInSet] %s", errx.String()))
			return false, "", errx
		}
		u.lgr.Debug(fmt.Sprintf("[Database] [Users] [GetVerificationStatus] [Scan] %s", errx.Error()))
		return false, "", errx
	}

	return verified, verificationKey.String, nil
}

func (u *users) Get(emailID string) (types.User, *erx.Erx) {
	query := `SELECT user_id, encryption_key, key_hash, verification_key, verified FROM users WHERE email=@email;`

	var userID types.UserID
	var encryptionKey, keyHash string
	var verificationKey sql.NullString
	var verificationStatus bool

	row := u.db.QueryRow(query, sql.Named("email", emailID))
	err := row.Scan(&userID, &encryptionKey, &keyHash, &verificationKey, &verificationStatus)
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			u.lgr.Error(fmt.Sprintf("[Database] [Users] [Get] [Scan] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return types.User{}, errx
		}
		if errors.Is(err, sql.ErrNoRows) {
			errx = erx.WithArgs(errx, erx.SeverityInfo, custom_errors.NoRowsInResultSet)
			u.lgr.Info(fmt.Sprintf("[Database] [Users] [Get] [Scan] [ErrSQLNoResultsInSet] %s", errx.String()))
			return types.User{}, errx
		}
		u.lgr.Debug(fmt.Sprintf("[Database] [Users] [Get] [Scan] %s", errx.Error()))
		return types.User{}, errx
	}

	return types.User{
		ID:              userID,
		Email:           emailID,
		EncryptionKey:   encryptionKey,
		KeyHash:         keyHash,
		VerificationKey: verificationKey.String,
		Verified:        verificationStatus,
	}, nil
}

func (u *users) Create(emailID string, encryptionKey string, keyHash string, vetkn string) (types.UserID, *erx.Erx) {
	query := `INSERT INTO users (email, encryption_key, key_hash, verification_key) OUTPUT inserted.user_id
	VALUES (@email, @key, @hash, @veKey);`

	var userID types.UserID

	row := u.db.QueryRow(query, sql.Named("email", emailID), sql.Named("key", encryptionKey),
		sql.Named("hash", keyHash), sql.Named("veKey", vetkn))
	err := row.Err()
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			errx = erx.WithArgs(errx, erx.SeverityError)
			u.lgr.Error(fmt.Sprintf("[Database] [Users] [Create] [QueryRow] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			return 0, errx
		}
		u.lgr.Debug(fmt.Sprintf("[Database] [Users] [Create] [QueryRow] %s", err.Error()))
		return 0, errx
	}

	err = row.Scan(&userID)
	if err != nil {
		sqlErr, errx := checkForSQLError(err)
		if sqlErr != nil {
			switch sqlErr.Number {
			case 2601:
				errx = erx.WithArgs(errx, custom_errors.DuplicateRecordInsertion, erx.SeverityInfo)
				u.lgr.Info(fmt.Sprintf("[Database] [Users] [Create] [Scan] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			default:
				errx = erx.WithArgs(errx, erx.SeverityError)
				u.lgr.Error(fmt.Sprintf("[Database] [Users] [Create] [Scan] [sqlErr] %d : %s", sqlErr.Number, sqlErr.Error()))
			}
			return 0, errx
		}
		u.lgr.Debug(fmt.Sprintf("[Database] [Users] [Create] [Scan] %s", errx.Error()))
		return 0, errx
	}

	return userID, nil
}
