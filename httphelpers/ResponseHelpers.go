package httphelpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
WriteJson writes JSON content to the response writer.
*/
func WriteJson(w http.ResponseWriter, status int, value any) {
	var (
		err error
		b   []byte
	)

	w.Header().Set("Content-Type", "application/json")

	if b, err = json.Marshal(value); err != nil {
		b, _ = json.Marshal(struct {
			Message    string `json:"message"`
			Suggestion string `json:"suggestion"`
		}{
			Message:    "Error marshaling value for writing",
			Suggestion: "See error log for more information",
		})

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = fmt.Fprintf(w, "%s", string(b))
		return
	}

	if status > 299 {
		w.WriteHeader(status)
	}

	_, _ = fmt.Fprintf(w, "%s", string(b))
}

/*
JsonOK is a convenience wrapper to send a 200 with an
arbitrary structure.
*/
func JsonOK(w http.ResponseWriter, value any) {
	WriteJson(w, http.StatusOK, value)
}

/*
JsonBadRequest is a convenience wrapper to send a 400 with an
arbitrary structure.
*/
func JsonBadRequest(w http.ResponseWriter, value any) {
	WriteJson(w, http.StatusBadRequest, value)
}

/*
JsonInternalServerError is a convenience wrapper to send a 500 with an
arbitrary structure.
*/
func JsonInternalServerError(w http.ResponseWriter, value any) {
	WriteJson(w, http.StatusInternalServerError, value)
}

/*
JsonErrorMessage is a convenience wrapper to send a JSON body with
the specified status code, and a body that looks like this:

	{"message": "<message> goes here"}
*/
func JsonErrorMessage(w http.ResponseWriter, status int, message string) {
	result := make(map[string]string)
	result["message"] = message

	WriteJson(w, status, result)
}
