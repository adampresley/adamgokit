package httphelpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
WriteHtml writes content to the response write with a text/html header.
*/
func WriteHtml(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "text/html")

	if status > 299 {
		w.WriteHeader(status)
	}

	_, _ = fmt.Fprintf(w, "%v", value)
}

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
HtmlOK is a convenience wrapper to send a 200 with an arbitrary HTML body
*/
func HtmlOK(w http.ResponseWriter, value any) {
	WriteHtml(w, http.StatusOK, value)
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

/*
JsonUnauthorized is a convenience wrapper to send a 401 with an
arbitrary value.
*/
func JsonUnauthorized(w http.ResponseWriter, value any) {
	WriteJson(w, http.StatusUnauthorized, value)
}
