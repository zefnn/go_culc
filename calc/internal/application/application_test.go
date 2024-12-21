package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCalc_ValidExpression(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleCalc))
	defer server.Close()

	reqBody := RequestBody{
		Expression: "2+3*4",
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(server.URL+"/api/v1/calculate", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var respBody ResponseBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	expected := float64(14)
	if respBody.Result != expected {
		t.Errorf("Expected result %f, got %f", expected, respBody.Result)
	}
}

func TestHandleCalc_DivisionByZero(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleCalc))
	defer server.Close()

	reqBody := RequestBody{
		Expression: "5/0",
	}
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(server.URL+"/api/v1/calculate", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}

	var respBody ResponseBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	expectedError := "Expression is not valid: division by zero"
	if respBody.Error != expectedError {
		t.Errorf("Expected error message %s, got %s", expectedError, respBody.Error)
	}
}

func TestHandleCalc_UnsupportedSymbol(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleCalc))
	defer server.Close()

	reqBody := RequestBody{
		Expression: "10a",
	}
	BodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.Post(server.URL+"/api/v1/calculate", "application/json", bytes.NewBuffer(BodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, resp.StatusCode)
	}

	var respBody ResponseBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	expectedError := "Internal server error: expression is not valid: unsupported symbol"
	if respBody.Error != expectedError {
		t.Errorf("Expected error message %s, got %s", expectedError, respBody.Error)
	}

}

func TestHandleCalc_InvalidMethod(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/api/v1/calculate", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handleCalc(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, rr.Code)
	}

	var respBody ResponseBody
	err = json.NewDecoder(rr.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	expectedError := "Метод не поддерживается!!!"
	if respBody.Error != expectedError {
		t.Errorf("Expected error message %s, got %s", expectedError, respBody.Error)
	}
}

func TestHandleCalc_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleCalc))
	defer server.Close()

	bodyBytes := []byte(`{ "expression": "2+3"*4}`)

	resp, err := http.Post(server.URL+"/api/v1/calculate", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	var respBody ResponseBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	expectedError := "Недопустимый запрос!!!"
	if respBody.Error != expectedError {
		t.Errorf("Expected error message %s, got %s", expectedError, respBody.Error)
	}
}

func TestHandleCalc_MissingExpression(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(handleCalc))
	defer server.Close()

	bodyBytes := []byte(`{}`)

	resp, err := http.Post(server.URL+"/api/v1/calculate", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("Expected status code %d, got %d", http.StatusUnprocessableEntity, resp.StatusCode)
	}

	var respBody ResponseBody
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		t.Fatal(err)
	}

	expectedError := "Expression is not valid"
	if respBody.Error != expectedError {
		t.Errorf("Expected error message %s, got %s", expectedError, respBody.Error)
	}
}
