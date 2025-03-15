package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

type LogMiddleware struct {
	Handler http.Handler
}

func (l *LogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s %s from %s\n", r.Method, r.URL.Path, r.RemoteAddr)
	body, _ := io.ReadAll(r.Body)
	log.Printf("Request Body: %s\n", string(body))
	r.Body = io.NopCloser(io.MultiReader(bytes.NewReader(body), r.Body))
	l.Handler.ServeHTTP(w, r)
}

func calculate(w http.ResponseWriter, r *http.Request, _ httprouter.Params, op func(float64, float64) (float64, error)) {
	var requestData RequestData
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, `{"error":"must provide number1 (int) and number2 (int)"}`, http.StatusBadRequest)
		return
	}
	result, err := op(requestData.Number1, requestData.Number2)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ResponseData{Result: result})
}

func add(a, b float64) (float64, error)      { return a + b, nil }
func subtract(a, b float64) (float64, error) { return a - b, nil }
func multiply(a, b float64) (float64, error) { return a * b, nil }
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return a / b, nil
}

func main() {
	router := httprouter.New()
	router.POST("/add", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		calculate(w, r, p, add)
	})
	router.POST("/subtract", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		calculate(w, r, p, subtract)
	})
	router.POST("/multiply", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		calculate(w, r, p, multiply)
	})
	router.POST("/divide", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		calculate(w, r, p, divide)
	})

	handler := cors.Default().Handler(&LogMiddleware{router})
	http.ListenAndServe(":8080", handler)
}
