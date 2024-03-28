package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golangdemo/config"
	"html/template"
	"log"
	"net/http"
	"time"
)

type PageData struct {
	ClientID    string
	RedirectURI string
}

func main() {
	config.LoadEnv()
	config.ConnectDB()
	http.HandleFunc("/initializ/v1/home", func(w http.ResponseWriter, r *http.Request) {
		var redirectURI string
		values := r.URL.Query()
		state := values.Get("state")
		clientID := "Iv1.1e998f2844d38483"
		if state == "YXV0aG9yaXplY2xp" {
			redirectURI = "http://localhost:8088/create-app"
		} else {
			redirectURI = "http://localhost:3000/create-app"
		}

		data := PageData{
			ClientID:    clientID,
			RedirectURI: redirectURI,
		}

		tmpl := template.Must(template.New("index").Parse(`
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>GitHub App Installation</title>
				<style>
					body {
						font-family: Arial, sans-serif;
						text-align: center;
						margin: 0;
						background-color: #f4f7fc;
						height: 100vh;
						display: flex;
						justify-content: center;
						align-items: center;
					}

					.container {
						background-color: #ffffff;
						padding: 20px;
						border-radius: 10px;
						box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
						max-width: 400px;
						margin: auto;
					}

					h1 {
						color: #333;
					}

					p {
						color: #555;
					}

					ul {
						list-style-type: none;
						padding: 0;
					}

					li {
						display: inline;
						margin: 0 10px;
					}

					a {
						text-decoration: none;
						padding: 10px 20px;
						border: 1px solid #4CAF50;
						background-color: #4CAF50;
						color: #fff;
						border-radius: 5px;
						transition: background-color 0.3s;
					}

					a:hover {
						background-color: #45a049;
					}

					.footer {
						background-color: #4CAF50;
						color: #fff;
						padding: 10px;
						position: fixed;
						bottom: 0;
						width: 100%;
						text-align: center;
					}
				</style>
			</head>
			<body>
				<div class="container">
					<h1>Successfully installed GitHub app</h1>
					<p>Click the button below to authorize the app</p>
					<ul>
						<li><a href="https://github.com/login/oauth/authorize?client_id={{.ClientID}}&redirect_uri={{.RedirectURI}}">Authorize with GitHub</a></li>
					</ul>
				</div>
				<div class="footer">
					Made by Initializ Inc &copy; 2023
				</div>
			</body>
			</html>
		`))

		tmpl.Execute(w, data)
	})
	http.HandleFunc("/initializ/v1/users", func(w http.ResponseWriter, r *http.Request) {
		userrepo := config.GetCollection(config.DB, "User")
		filter := make(map[string]interface{})
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cursor, err := userrepo.Find(ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
		defer cursor.Close(ctx)

		var results []map[string]interface{} // You can define a struct to match the document structure

		// Iterate through the documents and append them to the results slice
		for cursor.Next(ctx) {
			var result map[string]interface{}
			if err := cursor.Decode(&result); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println(result)
			results = append(results, result)
		}

		if err := cursor.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convert the results to JSON
		jsonData, err := json.Marshal(results)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Set the content type and write the JSON data to the response
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})
	http.ListenAndServe(":9090", nil)
}
