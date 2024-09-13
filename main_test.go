package main

import (
	"acme/assertslibrary"
	"acme/db"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRootHandler(t *testing.T) {

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()

	// actual HTTP handler for the rootHandler function
	handler := http.HandlerFunc(rootHandler)

	// Serve the request
	// shorthand for "Hey there HTTP Handler! I would like you to process this httpRequest with this responseWriter."
	handler.ServeHTTP(rr, req)

	// Check status code
	// if status := rr.Code; status != http.StatusOK {
	// 	t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	// }

	// using asserts library
	assertslibrary.CheckStatusCode(rr.Code, http.StatusOK, t)

	// Check response body
	// expected := "Hello Learner"
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v, want %v", rr.Body.String(), expected)
	// }

	// using asserts library
	assertslibrary.CheckResponseBody(rr.Body.String(), "Hello Learner", t)

}

// Integration tests

func TestRootHandlerWithServer(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(rootHandler))

	// this will be the last thing that runs before the function completes
	// using the defer keyword
	defer server.Close()

	// Act
	response, err := http.Get(server.URL + "/")

	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}

	defer response.Body.Close()

	// Assert using assert library
	assertslibrary.CheckStatusCode(response.StatusCode, http.StatusOK, t)

	// Check response Body - read body into bodyBytes variable
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body %v", err)
	}

	assertslibrary.CheckResponseBody(string(bodyBytes), "Hello Learner", t)

}

func TestGetUsersHandler(t *testing.T) {
	// Arrange
	req, err := http.NewRequest("GET", "/api/users", nil)

	if err != nil {
		t.Fatal(err)
	}

	// New response recorder
	rr := httptest.NewRecorder()

	// handler
	handler := http.HandlerFunc(getUsers)

	// Arrange out expected result
	expected := []db.User{
		{ID: 1, Name: "Name 1"},
		{ID: 2, Name: "Name 2"},
		{ID: 3, Name: "Name 3"},
	}

	// Marshall/serialize the data
	expectedJSON, err := json.Marshal(expected)

	if err != nil {
		t.Fatalf("Failed to marshal expectedJSON: %v", err)
	}

	// Act

	//Serve request
	handler.ServeHTTP(rr, req)

	// Assert using assert library
	assertslibrary.CheckStatusCode(rr.Code, http.StatusOK, t)

	assertslibrary.CheckResponseBody(rr.Body.String(), string(expectedJSON), t)

	var actual []db.User // used to store our unmarshalled data

	// attempting to unmarshal/deserialize response body and then place it in actual
	if err := json.Unmarshal(rr.Body.Bytes(), &actual); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// This allows the JSON responses to be compared based on their logical content rather than their string representation.
	// if !reflect.DeepEqual(actual, expected) {
	// 	t.Errorf("handler returned unexpected body: got %v, want %v", actual, expected)
	// }

	assertslibrary.CheckActualJsonData(actual, expected, t)

}

func TestGetUsersHandlerWithServer(t *testing.T) {
	// Arrange
	server := httptest.NewServer(http.HandlerFunc(getUsers))

	expected := []db.User{
		{ID: 1, Name: "Name 1"},
		{ID: 2, Name: "Name 2"},
		{ID: 3, Name: "Name 3"},
	}

	// Marshall/serialize the data
	expectedJSON, err := json.Marshal(expected)

	if err != nil {
		t.Fatalf("Failed to marshal expectedJSON: %v", err)
	}

	defer server.Close()

	// Act
	response, err := http.Get(server.URL + "/api/users")

	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}

	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("Failed to read response body %v", err)
	}
	var actual []db.User // used to store our unmarshalled data

	// attempting to unmarshal/deserialize response body and then place it in actual
	if err := json.Unmarshal(bodyBytes, &actual); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}

	// Assert
	assertslibrary.CheckStatusCode(response.StatusCode, http.StatusOK, t)

	assertslibrary.CheckResponseBody(string(bodyBytes), string(expectedJSON), t)

	assertslibrary.CheckActualJsonData(actual, expected, t)

}

func TestCreateUserHandler(t *testing.T) {
	// Arrange
	expectedString := db.User{
		ID: 4, Name: "Eddie",
	}

	expectedJSON, err := json.Marshal(expectedString)

	if err != nil {
		t.Fatalf("Failed to marshal expectedJSON: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/users", bytes.NewReader(expectedJSON))

	if err != nil {
		t.Fatal(err)
	}

	// response recorder
	rr := httptest.NewRecorder()

	// handler
	handler := http.HandlerFunc(createUser)

	// Act

	// server request
	handler.ServeHTTP(rr, req)

	assertslibrary.CheckStatusCode(rr.Code, http.StatusCreated, t)

}
