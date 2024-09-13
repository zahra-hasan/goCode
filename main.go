package main

import (
	"acme/db"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func rootHandler(writer http.ResponseWriter, request *http.Request) {
	io.WriteString(writer, "Hello Learner")
}

func createUser(writer http.ResponseWriter, request *http.Request) {
	// Make sure that the JSON we have been sent is a valid user object.
	var user db.User
	err := json.NewDecoder(request.Body).Decode(&user)

	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}

	// Call adduser and returns the id once added
	id := db.AddUser(user)

	writer.WriteHeader(http.StatusCreated)
	fmt.Fprintf(writer, "User created successfully: %d", id)

}

func getUsers(writer http.ResponseWriter, request *http.Request) {

	// output for the console
	fmt.Printf("got /api/users request.\n")

	// using our in memory database from the db package
	users := db.GetUsers()

	// Marshal or serialize users slice into JSON
	// this will do the job of turning the above into a byte array
	usersJSON, errMarshal := json.Marshal(users)

	if errMarshal != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// set Content-Type header of the respsonse which we want to be JSON (application/json)
	// indicates our response contains JSON
	writer.Header().Set("Content-Type", "application/json")

	// write marshalled data to the writer - does not need []byte() as marshalling has been done
	_, err := writer.Write(usersJSON)

	if err != nil {
		// if an error occurs use the http.Error(writer, a string to output, and error code)
		http.Error(writer, "Internal server error", http.StatusInternalServerError)
		return
	}

	io.WriteString(writer, "")
}

func getSingleUser(writer http.ResponseWriter, request *http.Request) {
	// this path value comes from the custom part of the path
	idStr := request.PathValue("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		fmt.Println("Error parsing ID:", err)
		http.Error(writer, "Bad request", http.StatusBadRequest)
		return
	}

	user := db.GetUser(id)

	// Converts user into a format that can be seen in HTML
	json.NewEncoder(writer).Encode(user)

	// fmt.Fprintf(writer, "get user with id=%v\n", idStr)
}

func deleteSingleUser(writer http.ResponseWriter, request *http.Request) {
	idString := request.PathValue("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		fmt.Println("Error parsing ID", err)
		http.Error(writer, "Bad Request", http.StatusBadRequest)
	}

	users := db.DeleteUser(id)

	json.NewEncoder(writer).Encode(users)

	// fmt.Fprintf(writer, "The user you want to delete is %v", idString)
}

func updateSingleUser(writer http.ResponseWriter, request *http.Request) {
	var userData db.User
	err := json.NewDecoder(request.Body).Decode(&userData)

	if err != nil {
		fmt.Println("Error decoding request body:", err)
		http.Error(writer, "Bad Request", http.StatusBadRequest)
		return
	}

	idString := request.PathValue("id")
	id, err := strconv.Atoi(idString)

	if err != nil {
		fmt.Println("Error parsing ID", err)
		http.Error(writer, "Bad Request", http.StatusBadRequest)
	}

	user := db.UpdateUser(id, userData)

	writer.WriteHeader(http.StatusCreated)
	fmt.Fprintf(writer, "User updated successfully: %v", user)
}

// Function to stop any Cors errors
func CorsMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(writer, request)
	})
}

func main() {

	//set up our multiplexer - we then use the router.HandleFunc to handle
	router := http.NewServeMux()

	// router.HandlFunc can handle multiple routes
	router.HandleFunc("GET /", rootHandler)
	router.HandleFunc("GET /api/users", getUsers)
	router.HandleFunc("POST /api/users", createUser)
	router.HandleFunc("GET /api/users/{id}", getSingleUser)
	router.HandleFunc("DELETE /api/users/{id}", deleteSingleUser)
	router.HandleFunc("PUT /api/users/{id}", updateSingleUser)

	// Start Server here
	fmt.Println("Server listening on port 8080")

	// make sure to pass "router" to ListenAndServe
	err := http.ListenAndServe(":8080", CorsMiddleWare(router))

	if err != nil {
		fmt.Println("Error starting", err)
	}
}
