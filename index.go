package idiota

import (
	"database/sql/driver"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/Craftserve/id-iota/pkg/base36"
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

// Not usual function so it does need to be part of the standard (because of case sensitivity we tend to use sometimes numbers)
func FromUint64(input uint64) (id Id, err error) {
	b := make([]byte, 8)

	binary.BigEndian.PutUint64(b, input)

	err = id.UnmarshalBinary(b)
	if err != nil {
		return id, err
	}

	return id, nil
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
	return []byte(base36.Encode(uint64((uint64(id.ts) << 32) + uint64(id.rand)))), nil
}

func (id *Id) UnmarshalText(data []byte) error {
	if len(data) > 13 {
		return ErrInvalidByteLength
	}

	bytes := base36.DecodeToBytes(string(data))

	return id.UnmarshalBinary(bytes)
}

func (id Id) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, id.String())), nil
}

func (id *Id) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	d := []byte(s)

	return id.UnmarshalText(d)
}

func (id Id) String() string {
	return base36.Encode(uint64((uint64(id.ts) << 32) + uint64(id.rand)))
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
	case []byte:
		if len := len(src.([]byte)); len > 13 {
			return fmt.Errorf("Scan: unable to scan []byte of length %d into Id-Iota Id", len)
		}

		err := id.UnmarshalText(src.([]byte))
		if err != nil {
			return fmt.Errorf("Scan: unable to scan while unmarshalling []byte %s into Id-Iota Id", src)
		}
		return nil
	case string:
		if src == nil {
			return nil
		}

		err := id.UnmarshalText(src.([]byte))
		if err != nil {
			return fmt.Errorf("Scan: unable to scan while unmarshalling string %s into Id-Iota Id", src)
		}
		return nil
	default:
		return fmt.Errorf("Scan: unable to scan type %T into Id-Iota Id", src)
	}
}

func (id Id) Value() (driver.Value, error) {
	return id.String(), nil
}
