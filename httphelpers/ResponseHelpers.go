package httphelpers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func writePlain(w http.ResponseWriter, contentType string, status int, value any) {
	w.Header().Set("Content-Type", contentType)

	if status > 299 {
		w.WriteHeader(status)
	}

	_, _ = fmt.Fprintf(w, "%v", value)
}

/*
WriteHtml writes content to the response write with a text/html header.
*/
func WriteHtml(w http.ResponseWriter, status int, value any) {
	writePlain(w, "text/html", status, value)
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
WriteText writes content to the response writer with a text/plain header.
*/
func WriteText(w http.ResponseWriter, status int, value any) {
	writePlain(w, "text/plain", status, value)
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

/*
TextOK is a convenience wrapper to send a 200 with an
arbitrary text body.
*/
func TextOK(w http.ResponseWriter, value any) {
	WriteText(w, http.StatusOK, value)
}

/*
TextBadRequest is a convenience wrapper to send a 400 with an
arbitrary text body.
*/
func TextBadRequest(w http.ResponseWriter, value any) {
	WriteText(w, http.StatusBadRequest, value)
}

/*
TextBadRequest is a convenience wrapper to send a 500 with an
arbitrary text body.
*/
func TextInternalServerError(w http.ResponseWriter, value any) {
	WriteText(w, http.StatusInternalServerError, value)
}

/*
TextBadRequest is a convenience wrapper to send a 401 with an
arbitrary text body.
*/
func TextUnauthorized(w http.ResponseWriter, value any) {
	WriteText(w, http.StatusUnauthorized, value)
}

/*
DownloadCSV writes bytes as a CSV file to the response writer and downloads it.
*/
func DownloadCSV(w http.ResponseWriter, filename string, csvContent []byte) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", strconv.Itoa(len(csvContent)))

	w.Write(csvContent)
}

func IsSuccessRange(status int) bool {
	return status >= 200 && status < 300
}
