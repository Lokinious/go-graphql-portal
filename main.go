package main

import (
	"bytes"
	"encoding/json"
	"fmt"  // Import fmt package
	"html/template"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

var client = &http.Client{}
var tmpl *template.Template

func init() {
	// Load HTML template
	tmpl = template.Must(template.New("index").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>GraphQL Query Portal</title>
	</head>
	<body>
		<h1>GraphQL Query Portal</h1>
		<form action="/query" method="post">
			<label for="query">GraphQL Query:</label>
			<textarea id="query" name="query" rows="5" cols="50"></textarea><br>
			<input type="submit" value="Submit">
		</form>
	</body>
	</html>
	`))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// Serve the HTML form
	tmpl.Execute(w, nil)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	// Handle GraphQL query
	if r.Method == http.MethodPost {
		query := r.FormValue("query")
		logrus.Infof("Received GraphQL query: %s", query)

		// Send the query to the go-graphql-app
		result, err := sendQueryToGraphQLApp(query)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error sending query to go-graphql-app: %s", err), http.StatusInternalServerError)
			return
		}

		// For simplicity, in this example, we just print the result to the response.
		w.Write([]byte(fmt.Sprintf("Query received. Result from go-graphql-app:\n%s", result)))
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func sendQueryToGraphQLApp(query string) (string, error) {
	// Define the GraphQL request payload
	requestPayload := map[string]interface{}{
		"query": query,
	}

	// Serialize the payload to JSON
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return "", err
	}

	// Create the GraphQL request
	req, err := http.NewRequest("POST", "http://go-graphql-app-service:8080/graphql", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the request to the go-graphql-app
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}

func main() {
	// Set up logging
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)

	// Set up HTTP routes
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/query", queryHandler)

	// Start the server
	logrus.Fatal(http.ListenAndServe(":8081", nil))
}
