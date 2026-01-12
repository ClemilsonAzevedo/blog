package pkg

import (
	"database/sql/driver"
	"errors"

	"go.bryk.io/pkg/ulid"
)

type ULID struct {
	internal ulid.ULID
}

func NewULID() (ULID, error) {
	id, err := ulid.New()
	if err != nil {
		return ULID{}, err
	}
	return ULID{internal: id}, nil
}

func MustNewULID() ULID {
	u, err := NewULID()
	if err != nil {
		panic(err)
	}
	return u
}

func ParseULID(s string) (ULID, error) {
	id, err := ulid.Parse(s)
	if err != nil {
		return ULID{}, err
	}
	return ULID{internal: id}, nil
}

func (u ULID) String() string {
	return u.internal.String()
}

func (u ULID) Value() (driver.Value, error) {
	return u.internal.String(), nil
}

func (u *ULID) Scan(src any) error {
	switch v := src.(type) {
	case string:
		id, err := ulid.Parse(v)
		if err != nil {
			return err
		}
		u.internal = id
		return nil
	case []byte:
		id, err := ulid.Parse(string(v))
		if err != nil {
			return err
		}
		u.internal = id
		return nil
	default:
		return errors.New("unsupported type for ULID Scan")
	}
}
