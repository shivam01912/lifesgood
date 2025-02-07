package llm

import (
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"strconv"
)

// readRequestJSON expects req to have a JSON content type with a body that
// contains a JSON-encoded value complying with the underlying type of target.
// It populates target, or returns an error.
func readRequestJSON(req *http.Request, target any) error {
	contentType := req.Header.Get("Content-Type")
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}
	if mediaType != "application/json" {
		return fmt.Errorf("expect application/json Content-Type, got %s", mediaType)
	}

	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(target)
}

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v any) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	respTest, err := strconv.Unquote(string(js))
	if err != nil {
		log.Printf("Error un-escaping the LLM response error=%s reponseText=%s\n", err.Error(), string(js))
	}

	_, err = w.Write([]byte(respTest))
	if err != nil {
		return
	}
}
