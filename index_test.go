package idiota_test

import (
	"encoding/binary"
	"testing"
	"time"

	idiota "github.com/Craftserve/id-iota"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	id := idiota.NewId(nil, nil)

	assert.Len(t, id.String(), 13)
}

func TestNewWithTime(t *testing.T) {
	yesterday := time.Now().AddDate(0, 0, -1).Round(time.Second)

	id := idiota.NewId(&yesterday, nil)
	assert.Len(t, id.String(), 13)
	assert.Equal(t, yesterday, id.Time())
}

func TestNewWithRandom(t *testing.T) {

	random := uint32(123456789)

	id := idiota.NewId(nil, &random)
	assert.Len(t, id.String(), 13)

	fullIdUint64 := id.UInt64()
	randomPart := uint32(fullIdUint64)

	assert.Equal(t, random, randomPart)
}

func TestNewWithRandomAndTime(t *testing.T) {
	yesterday := time.Now().AddDate(0, 0, -1).Round(time.Second)

	random := uint32(123456789)

	id := idiota.NewId(&yesterday, &random)
	assert.Len(t, id.String(), 13)

	fullIdUint64 := id.UInt64()

	timePart := uint32(fullIdUint64 >> 32)
	randomPart := uint32(fullIdUint64)

	assert.Equal(t, random, randomPart)

	assert.Equal(t, uint32(yesterday.Unix()), timePart)
}

func TestUnmarshalText(t *testing.T) {
	id := idiota.NewId(nil, nil)
	idString := id.String()

	var id2 idiota.Id
	err := id2.UnmarshalText([]byte(idString))
	assert.NoError(t, err)

	assert.Equal(t, id, id2)
}

func TestMarshalText(t *testing.T) {
	id := idiota.NewId(nil, nil)
	idString := id.String()

	idBytes, err := id.MarshalText()
	assert.NoError(t, err)

	assert.Equal(t, idString, string(idBytes))
}

func TestMarshalBinary(t *testing.T) {
	id := idiota.NewId(nil, nil)

	idBytes, err := id.MarshalBinary()
	assert.NoError(t, err)

	assert.Len(t, idBytes, 8)
	assert.Equal(t, uint32(id.UInt64()>>32), binary.BigEndian.Uint32(idBytes[0:4]))
	assert.Equal(t, uint32(id.UInt64()), binary.BigEndian.Uint32(idBytes[4:8]))
}

func TestUnmarshalBinary(t *testing.T) {
	id := idiota.NewId(nil, nil)

	idBytes, err := id.MarshalBinary()
	assert.NoError(t, err)

	var id2 idiota.Id
	err = id2.UnmarshalBinary(idBytes)
	assert.NoError(t, err)

	assert.Equal(t, id, id2)
}

func TestScanString(t *testing.T) {

	// string
	id := idiota.NewId(nil, nil)

	var id2 idiota.Id
	err := id2.Scan(id.String())
	assert.NoError(t, err)
	assert.Equal(t, id, id2)

	// bytes (text-based)
	idBytes, err := id.MarshalText()
	assert.NoError(t, err)

	err = id2.Scan(idBytes)
	assert.NoError(t, err)
	assert.Equal(t, id, id2)

}

func TestScanUint64(t *testing.T) {

	var id idiota.Id
	now := time.Now().Round(time.Second).Unix()

	caseValueRandomPart := uint32(123456789)
	caseValue := uint64(now)<<32 | uint64(caseValueRandomPart)
	err := id.Scan(caseValue)
	assert.NoError(t, err)

	assert.Equal(t, uint32(now), uint32(id.Time().Unix()))
	assert.Equal(t, caseValue, id.UInt64())

	randomPart := uint32(caseValue)
	assert.Equal(t, randomPart, uint32(id.UInt64()))
}

func TestScanNil(t *testing.T) {
	var id idiota.Id
	err := id.Scan(nil)
	assert.Error(t, err)
}

func TestScanInvalidType(t *testing.T) {
	var id idiota.Id
	err := id.Scan([]interface{}{})
	assert.Error(t, err)
}

func TestFromUint64(t *testing.T) {
	original := idiota.NewId(nil, nil)
	u := original.UInt64()

	id2, err := idiota.FromUint64(u)
	assert.NoError(t, err)
	assert.Equal(t, original, id2)
}

func TestMarshalUnmarshalJSON(t *testing.T) {
	original := idiota.NewId(nil, nil)

	data, err := original.MarshalJSON()
	assert.NoError(t, err)

	var id2 idiota.Id
	err = id2.UnmarshalJSON(data)
	assert.NoError(t, err)

	assert.Equal(t, original, id2)
}

func TestValue(t *testing.T) {
	id := idiota.NewId(nil, nil)

	val, err := id.Value()
	assert.NoError(t, err)

	strVal, ok := val.(string)
	assert.True(t, ok)
	assert.Equal(t, id.String(), strVal)
}
