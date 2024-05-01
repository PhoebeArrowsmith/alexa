package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

type Input struct {
	Speech string `json:"speech"`
}

type Answer struct {
	Text string `json:"text"`
}

type Result struct {
	DisplayText string `json:"DisplayText"`
}

func main() {
	http.ListenAndServe(":3002", Router())
}

func Router() http.Handler {
	r := mux.NewRouter()
	/* controller */
	r.HandleFunc("/stt", readInput).Methods("POST")
	return r
}

func readInput(w http.ResponseWriter, r *http.Request) {
	var i Input
	err := json.NewDecoder(r.Body).Decode(&i)
	if err == nil {
		reply, err, StatusCode := callAPI(i)
		if err == 0 {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			response := &Answer{Text: reply}
			json.NewEncoder(w).Encode(response)
		} else {
			http.Error(w, reply, StatusCode)
		}
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	return
}

func callAPI(i Input) (string, int64, int) {

	MS_KEY := "2d982d27947348929b4e14648f4d16c1"

	// Take the input, which is base64 encoded and convert it to binary
	data := make([]byte, base64.StdEncoding.DecodedLen(len(i.Speech)))
	n, err := base64.StdEncoding.Decode(data, []byte(i.Speech))
	if err != nil {
		return "Internal Server Error", -1, 500
	}

	// set up the call to the MS speech recognition
	client := &http.Client{}
	r, _ := http.NewRequest(http.MethodPost, "https://uksouth.stt.speech.microsoft.com/speech/recognition/conversation/cognitiveservices/v1?language=en-US", bytes.NewReader(data[:n]))
	r.Header.Add("Content-Type", "audio/wav;samplerate=16000")
	r.Header.Add("Ocp-Apim-Subscription-Key", MS_KEY)

	// call API and store result in response
	response, err := client.Do(r)

	if err != nil {
		return "Service Unavailable", -1, 503
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "Internal Server Error", -1, 500
	}

	//Unmarshall DisplayText into a string
	var responses Result
	err1 := json.Unmarshal([]byte(body), &responses)
	if err1 != nil {
		return "Internal Server Error", -1, 500
	}

	api_resp := responses.DisplayText

	return api_resp, 0, response.StatusCode

}
