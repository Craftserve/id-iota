package idiota

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/akamensky/base58"
)

var ErrInvalidByteLength = errors.New("Invalid byte length")

type Id struct {
	ts   uint32 // Its important to know that this is a UNIX timestamp which will overflow in 2106 (https://en.wikipedia.org/wiki/Year_2038_problem#Solutions)
	rand uint32
}

var RandomFunc = rand.Uint32

func NewId(inputTime *time.Time, inputRandom *uint32) (id Id) {
	var idTime = time.Now().Unix()
	var idRandom uint32

	if inputTime != nil {
		idTime = inputTime.Unix()
	}

	if inputRandom != nil {
		idRandom = *inputRandom
	}

	if inputRandom == nil {
		idRandom = RandomFunc()
	}

	id = Id{
		ts:   uint32(idTime),
		rand: idRandom,
	}

	return id
}

func (id Id) MarshalBinary() (idBytes []byte, err error) {
	idBytes = make([]byte, 8)

	binary.BigEndian.PutUint32(idBytes[0:4], id.ts)
	binary.BigEndian.PutUint32(idBytes[4:8], id.rand)

	return idBytes, nil
}

func (id *Id) UnmarshalBinary(data []byte) error {
	if len(data) != 8 {
		return ErrInvalidByteLength
	}

	id.ts = binary.BigEndian.Uint32(data[0:4])
	id.rand = binary.BigEndian.Uint32(data[4:8])

	return nil
}

func (id Id) MarshalText() ([]byte, error) {
	bytes, err := id.MarshalBinary()

	return []byte(base58.Encode(bytes)), err
}

func (id *Id) UnmarshalText(data []byte) error {
	bytes, err := base58.Decode(string(data))
	if err != nil {
		return err
	}

	return id.UnmarshalBinary(bytes)
}

func (id Id) String() string {
	bytes, _ := id.MarshalBinary()

	return base58.Encode(bytes)
}

func (id Id) UInt64() uint64 {
	return (uint64(id.ts) << 32) | uint64(id.rand)
}

func (id Id) Time() time.Time {
	return time.Unix(int64(id.ts), 0)
}

func (id *Id) Scan(src interface{}) error {
	switch src.(type) {
	case nil:
		return fmt.Errorf("Scan: unable to scan nil into Id-Iota Id")
	case uint64:
		bytes := make([]byte, 8)
		binary.BigEndian.PutUint64(bytes, src.(uint64))

		return id.UnmarshalBinary(bytes)
	case []byte:
		if len := len(src.([]byte)); len != 8 {
			return fmt.Errorf("Scan: unable to scan []byte of length %d into Id-Iota Id", len)
		}

		return id.UnmarshalBinary(src.([]byte))
	case string:
		if src == nil {
			return nil
		}

		return id.UnmarshalText([]byte(src.(string)))
	default:
		return fmt.Errorf("Scan: unable to scan type %T into Id-Iota Id", src)
	}
}

func (id Id) Value() (driver.Value, error) {
	return id.String(), nil
}
