package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

type RequestData struct {
	Number1 float64 `json:"number1"`
	Number2 float64 `json:"number2"`
}

type ResponseData struct {
	Result float64 `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type LogMiddleware struct {
	Handler http.Handler
}

func (l *LogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
	body, _ := io.ReadAll(r.Body)
	log.Printf("Request Body: %s\n", string(body))

	r.Body = io.NopCloser(io.MultiReader(bytes.NewReader(body), r.Body)) // Reset the body so it can be read again

	l.Handler.ServeHTTP(w, r)
}

func add(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "must provide number1 (int) and number2 (int)"})
		return
	}

	responseData := ResponseData{Result: requestData.Number1 + requestData.Number2}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responseData)
}

func subtract(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "must provide number1 (int) and number2 (int)"})
		return
	}

	responseData := ResponseData{Result: requestData.Number1 - requestData.Number2}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responseData)
}

func multiply(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "must provide number1 (int) and number2 (int)"})
		return
	}

	responseData := ResponseData{Result: requestData.Number1 * requestData.Number2}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responseData)
}

func divide(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var requestData RequestData
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "must provide number1 (int) and number2 (int)"})
		return
	}

	if requestData.Number2 == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "cannot divide by zero"})
		return
	}

	responseData := ResponseData{Result: requestData.Number1 / requestData.Number2}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(responseData)
}

func main() {
	router := httprouter.New()

	router.POST("/add", add)
	router.POST("/subtract", subtract)
	router.POST("/multiply", multiply)
	router.POST("/divide", divide)

	handler := cors.Default().Handler(&LogMiddleware{router})

	http.ListenAndServe(":8080", handler)
}
