package httphelpers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/adampresley/adamgokit/httphelpers"
	"github.com/stretchr/testify/assert"
)

func TestWriteJSON(t *testing.T) {
	type TestingType struct {
		Key1 string `json:"key1"`
		Key2 int    `json:"key2"`
	}

	input := TestingType{
		Key1: "Adam",
		Key2: 10,
	}

	recorder := httptest.NewRecorder()
	httphelpers.WriteJson(recorder, http.StatusOK, input)

	result := recorder.Result()

	want := []byte(`{"key1":"Adam","key2":10}`)
	got, err := io.ReadAll(result.Body)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestJsonOK(t *testing.T) {
	type TestingType struct {
		Key1 string `json:"key1"`
		Key2 int    `json:"key2"`
	}

	input := TestingType{
		Key1: "Adam",
		Key2: 10,
	}

	recorder := httptest.NewRecorder()
	httphelpers.JsonOK(recorder, input)

	result := recorder.Result()

	want := []byte(`{"key1":"Adam","key2":10}`)
	got, err := io.ReadAll(result.Body)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusOK, result.StatusCode)
}

func TestJsonBadRequest(t *testing.T) {
	type TestingType struct {
		Key1 string `json:"key1"`
		Key2 int    `json:"key2"`
	}

	input := TestingType{
		Key1: "Adam",
		Key2: 10,
	}

	recorder := httptest.NewRecorder()
	httphelpers.JsonBadRequest(recorder, input)

	result := recorder.Result()

	want := []byte(`{"key1":"Adam","key2":10}`)
	got, err := io.ReadAll(result.Body)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
}

func TestJsonInternalServerError(t *testing.T) {
	type TestingType struct {
		Key1 string `json:"key1"`
		Key2 int    `json:"key2"`
	}

	input := TestingType{
		Key1: "Adam",
		Key2: 10,
	}

	recorder := httptest.NewRecorder()
	httphelpers.JsonInternalServerError(recorder, input)

	result := recorder.Result()

	want := []byte(`{"key1":"Adam","key2":10}`)
	got, err := io.ReadAll(result.Body)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

func TestJsonErrorMessage(t *testing.T) {
	recorder := httptest.NewRecorder()
	httphelpers.JsonErrorMessage(recorder, http.StatusInternalServerError, "bad news")

	result := recorder.Result()

	want := []byte(`{"message":"bad news"}`)
	got, err := io.ReadAll(result.Body)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

func TestJsonUnauthorized(t *testing.T) {
	type TestingType struct {
		Key1 string `json:"key1"`
		Key2 int    `json:"key2"`
	}

	input := TestingType{
		Key1: "Adam",
		Key2: 10,
	}

	recorder := httptest.NewRecorder()
	httphelpers.JsonUnauthorized(recorder, input)

	result := recorder.Result()

	want := []byte(`{"key1":"Adam","key2":10}`)
	got, err := io.ReadAll(result.Body)

	assert.NoError(t, err)
	assert.Equal(t, want, got)
	assert.Equal(t, http.StatusUnauthorized, result.StatusCode)
}
