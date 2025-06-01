package pkg

import (
	"encoding/json"
	"fap-server/models"
	"fmt"
	"net/http"
)

// Same as http.Error, but it sets the Content-Type header to application/json instead of text/plain
func JsonError(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, err)
}

func GenericResponseJson(result string, message string) string {
	rawJson, _ := json.Marshal(models.GenericResponse{
		Result:  result,
		Message: message,
	})
	return string(rawJson)
}
