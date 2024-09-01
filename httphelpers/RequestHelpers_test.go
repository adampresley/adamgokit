package httphelpers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/adampresley/adamgokit/httphelpers"
	"github.com/stretchr/testify/assert"
)

func TestGetFromRequest(t *testing.T) {
	wantString := "Adam"
	wantStringSlice := []string{"Adam", "Maryanne"}
	wantInt := 11
	wantIntSlice := []int{1, 2, 12}
	wantInt32 := int32(13)
	wantInt32Slice := []int32{14, 15}
	wantInt64 := int64(-16)
	wantInt64Slice := []int64{17, -18, 19, 20}
	wantUInt := uint(21)
	wantUIntSlice := []uint{22, 23}
	wantUInt32 := uint32(24)
	wantUInt32Slice := []uint32{25, 26}
	wantUInt64 := uint64(27)
	wantUInt64Slice := []uint64{28, 29}
	wantFloat32 := float32(30.1)
	wantFloat32Slice := []float32{30.2, 30.3}
	wantFloat64 := float64(31.1)
	wantFloat64Slice := []float64{31.2, 31.3}
	wantBool := true

	u := `/?string=Adam&`
	u += `stringslice=Adam&stringslice=Maryanne&`
	u += `int=11&`
	u += `intslice=1&intslice=2&intslice=12&`
	u += `int32=13&`
	u += `int32slice=14&int32slice=15&`
	u += `int64=-16&`
	u += `int64slice=17&int64slice=-18&int64slice=19&int64slice=20&`
	u += `uint=21&`
	u += `uintslice=22&uintslice=23&`
	u += `uint32=24&`
	u += `uint32slice=25&uint32slice=26&`
	u += `uint64=27&`
	u += `uint64slice=28&uint64slice=29&`
	u += `float32=30.1&`
	u += `float32slice=30.2&float32slice=30.3&`
	u += `float64=31.1&`
	u += `float64slice=31.2&float64slice=31.3&`
	u += `bool=true&`

	r := httptest.NewRequest(http.MethodGet, u, nil)

	gotString := httphelpers.GetFromRequest[string](r, "string")
	gotStringSlice := httphelpers.GetFromRequest[[]string](r, "stringslice")
	gotInt := httphelpers.GetFromRequest[int](r, "int")
	gotIntSlice := httphelpers.GetFromRequest[[]int](r, "intslice")
	gotInt32 := httphelpers.GetFromRequest[int32](r, "int32")
	gotInt32Slice := httphelpers.GetFromRequest[[]int32](r, "int32slice")
	gotInt64 := httphelpers.GetFromRequest[int64](r, "int64")
	gotInt64Slice := httphelpers.GetFromRequest[[]int64](r, "int64slice")
	gotUInt := httphelpers.GetFromRequest[uint](r, "uint")
	gotUIntSlice := httphelpers.GetFromRequest[[]uint](r, "uintslice")
	gotUInt32 := httphelpers.GetFromRequest[uint32](r, "uint32")
	gotUInt32Slice := httphelpers.GetFromRequest[[]uint32](r, "uint32slice")
	gotUInt64 := httphelpers.GetFromRequest[uint64](r, "uint64")
	gotUInt64Slice := httphelpers.GetFromRequest[[]uint64](r, "uint64slice")
	gotFloat32 := httphelpers.GetFromRequest[float32](r, "float32")
	gotFloat32Slice := httphelpers.GetFromRequest[[]float32](r, "float32slice")
	gotFloat64 := httphelpers.GetFromRequest[float64](r, "float64")
	gotFloat64Slice := httphelpers.GetFromRequest[[]float64](r, "float64slice")
	gotBool := httphelpers.GetFromRequest[bool](r, "bool")

	assert.Equal(t, wantString, gotString)
	assert.Equal(t, wantStringSlice, gotStringSlice)
	assert.Equal(t, wantInt, gotInt)
	assert.Equal(t, wantIntSlice, gotIntSlice)
	assert.Equal(t, wantInt32, gotInt32)
	assert.Equal(t, wantInt32Slice, gotInt32Slice)
	assert.Equal(t, wantInt64, gotInt64)
	assert.Equal(t, wantInt64Slice, gotInt64Slice)
	assert.Equal(t, wantUInt, gotUInt)
	assert.Equal(t, wantUIntSlice, gotUIntSlice)
	assert.Equal(t, wantUInt32, gotUInt32)
	assert.Equal(t, wantUInt32Slice, gotUInt32Slice)
	assert.Equal(t, wantUInt64, gotUInt64)
	assert.Equal(t, wantUInt64Slice, gotUInt64Slice)
	assert.Equal(t, wantFloat32, gotFloat32)
	assert.Equal(t, wantFloat32Slice, gotFloat32Slice)
	assert.Equal(t, wantFloat64, gotFloat64)
	assert.Equal(t, wantFloat64Slice, gotFloat64Slice)
	assert.Equal(t, wantBool, gotBool)
}

func TestGetStringListFromRequest(t *testing.T) {
	want := []string{"1", "5", "10"}
	r := httptest.NewRequest(http.MethodGet, "/?input=1,5,10", nil)

	got := httphelpers.GetStringListFromRequest(r, "input", ",")

	assert.Equal(t, want, got)
}

func TestReadJSONBody(t *testing.T) {
	type TestingType struct {
		Key1 string `json:"key1"`
		Key2 int    `json:"key2"`
	}

	input := `{"key1": "Adam", "key2": 10}`

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(input))
	r.Header.Set("Content-Type", "application/json")

	want := TestingType{
		Key1: "Adam",
		Key2: 10,
	}

	got := TestingType{}
	err := httphelpers.ReadJSONBody(r, &got)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
}
