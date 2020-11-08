package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"bufio"
	"time"
)

// Examples represents the examples loaded from examples.json.
type Examples []*Example

var (
	errListExamples  = errors.New("failed to list examples (please run in the examples folder)")
	errParseExamples = errors.New("failed to parse examples")
	token_stream = ""
	token_connect = ""
	sync time.Time
)

// Example represents an example loaded from examples.json.
type Example struct {
	Title       string `json:"title"`
	Link        string `json:"link"`
	Description string `json:"description"`
	Type        string `json:"type"`
	IsJS        bool
	IsWASM      bool
}

func main() {
	addr := flag.String("address", ":80", "Address to host the HTTP server on.")
	flag.Parse()

	log.Println("Listening on", *addr)
	err := serve(*addr)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func serve(addr string) error {
	// Load the examples
	examples, err := getExamples()
	if err != nil {
		return err
	}



	// Load the templates
	homeTemplate := template.Must(template.ParseFiles("index.html"))

	//send stream from camera
	http.HandleFunc("/sendtokenstreamer", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		//user := req.Form.Get("user")
		//id_camera := req.Form.Get("id_camera")
		token := req.Form.Get("token")

		if token_stream == ""{
			token_stream = token
		}

		//fmt.Printf(token_stream)

		// await for token_connect
		endTime := time.Now().Add(time.Second * 61)
		for time.Now().Before(endTime) {
			if token_connect != "" {
				tokenResponse := token_connect
				token_connect = ""
				fmt.Fprintln(w, tokenResponse)
				return
			}
			//fmt.Println("Esperando")
		}

		//fmt.Println("SALIENDO")
		fmt.Fprintln(w, "")
		return
	})

	//Receiver token from consumer for connect to streamer
	http.HandleFunc("/sendtokenconnect", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		token_connect = req.Form.Get("token")
		//fmt.Printf(token_stream)

		token := token_stream
		if token != ""{
			fmt.Fprintln(w,token_connect)
			return
		}

		fmt.Printf(token_connect)
		return
	})

	//Check if exists stream from user and id_camera
	http.HandleFunc("/checkstream", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		//user := req.Form.Get("user")
		//id_camera := req.Form.Get("id_camera")

		if token_stream != ""{
			fmt.Fprintln(w, token_stream)
			return
		}

		fmt.Fprintln(w, "token not found")
		return
	})

	//Syncronize clients, streamer check every minute if the timing is! = 0. Yes! = 0 will require the connection token in sync + 120
	http.HandleFunc("/checksyncronize", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		//synchronization for user and consumer before connection connection

		//send time for request sync
		fmt.Fprintln(w,sync.Format("03:04:05")+" "+time.Now().Format("03:04:05"))
		return

	})

	//Get token from consumer
	http.HandleFunc("/gettokenconsumer", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		//synchronization for user and consumer before connection connection

		//send time for request sync
		sync = time.Time{}
		token_stream = ""
		token := token_connect
		token_connect = ""
		fmt.Fprintln(w,token)
		return

	})

	//Add camera
	http.HandleFunc("/addcamera", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		user := req.Form.Get("user")
		loc := req.Form.Get("loc")

		fmt.Printf(user+" "+loc)
		return
		//TODO agregate new camera in DB
	})

	//Delete camera
	http.HandleFunc("/deletecamera", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		user := req.Form.Get("user")
		loc := req.Form.Get("id_camera")

		fmt.Printf(user+" "+loc)
		return
		//TODO delete camera in DB
	})

	//Active/desactive camera
	http.HandleFunc("/activecamera", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		user := req.Form.Get("user")
		loc := req.Form.Get("id_camera")

		fmt.Printf(user+" "+loc)
		return
		//TODO active/desactive camera in DB
	})

	//Login user
	http.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		email := req.Form.Get("email")
		loc := req.Form.Get("password")

		fmt.Printf(email+" "+loc)
		return
	})

	//Generate session token
	http.HandleFunc("/gentoken", func(w http.ResponseWriter, req *http.Request) {

		req.ParseForm()		// Parses the request body
		token := req.Form.Get("token") // x will be "" if parameter is not set

//		token := "eyJ0eXBlIjoib2ZmZXIiLCJzZHAiOiJ2PTBcclxubz0tIDc5MTkzNjY5MzczNDQxNjk3ODUgMiBJTiBJUDQgMTI3LjAuMC4xXHJcbnM9LVxyXG50PTAgMFxyXG5hPWdyb3VwOkJVTkRMRSAwIDFcclxuYT1tc2lkLXNlbWFudGljOiBXTVMgSUV3eExWcFBiVHEwUENvdGVIaG1WaUxkUTZpUU1OOFNQVHJkXHJcbm09YXVkaW8gNTczOTYgVURQL1RMUy9SVFAvU0FWUEYgMTExIDEwMyAxMDQgOSAwIDggMTA2IDEwNSAxMyAxMTAgMTEyIDExMyAxMjZcclxuYz1JTiBJUDQgMTUyLjE3MC4xLjE2M1xyXG5hPXJ0Y3A6OSBJTiBJUDQgMC4wLjAuMFxyXG5hPWNhbmRpZGF0ZTo0MjUyODc2MjU2IDEgdWRwIDIxMjIyNjAyMjMgMTkyLjE2OC4wLjE5NCA1NzM5NiB0eXAgaG9zdCBnZW5lcmF0aW9uIDAgbmV0d29yay1pZCAxIG5ldHdvcmstY29zdCAxMFxyXG5hPWNhbmRpZGF0ZTozMDE5Nzg0NDY0IDEgdGNwIDE1MTgyODA0NDcgMTkyLjE2OC4wLjE5NCA5IHR5cCBob3N0IHRjcHR5cGUgYWN0aXZlIGdlbmVyYXRpb24gMCBuZXR3b3JrLWlkIDEgbmV0d29yay1jb3N0IDEwXHJcbmE9Y2FuZGlkYXRlOjIwODM4OTYxNDggMSB1ZHAgMTY4NjA1MjYwNyAxNTIuMTcwLjEuMTYzIDU3Mzk2IHR5cCBzcmZseCByYWRkciAxOTIuMTY4LjAuMTk0IHJwb3J0IDU3Mzk2IGdlbmVyYXRpb24gMCBuZXR3b3JrLWlkIDEgbmV0d29yay1jb3N0IDEwXHJcbmE9aWNlLXVmcmFnOnYvb05cclxuYT1pY2UtcHdkOnNhNSsvSGd5WFFyaUpsR29ocmRDU09iV1xyXG5hPWljZS1vcHRpb25zOnRyaWNrbGVcclxuYT1maW5nZXJwcmludDpzaGEtMjU2IDgyOjcwOjlEOkI2OkI5OkIwOkM0OkFFOjIxOjg5OjY1OjFDOjM5OjM5OkE0OkJGOjZCOkJBOjMyOjhFOkU0OjY3OjgxOjM1OkQ5OjYzOkQyOkQ4OkQ0OjEwOjYyOkNCXHJcbmE9c2V0dXA6YWN0cGFzc1xyXG5hPW1pZDowXHJcbmE9ZXh0bWFwOjEgdXJuOmlldGY6cGFyYW1zOnJ0cC1oZHJleHQ6c3NyYy1hdWRpby1sZXZlbFxyXG5hPWV4dG1hcDoyIGh0dHA6Ly93d3cud2VicnRjLm9yZy9leHBlcmltZW50cy9ydHAtaGRyZXh0L2Ficy1zZW5kLXRpbWVcclxuYT1leHRtYXA6MyBodHRwOi8vd3d3LmlldGYub3JnL2lkL2RyYWZ0LWhvbG1lci1ybWNhdC10cmFuc3BvcnQtd2lkZS1jYy1leHRlbnNpb25zLTAxXHJcbmE9ZXh0bWFwOjQgdXJuOmlldGY6cGFyYW1zOnJ0cC1oZHJleHQ6c2RlczptaWRcclxuYT1leHRtYXA6NSB1cm46aWV0ZjpwYXJhbXM6cnRwLWhkcmV4dDpzZGVzOnJ0cC1zdHJlYW0taWRcclxuYT1leHRtYXA6NiB1cm46aWV0ZjpwYXJhbXM6cnRwLWhkcmV4dDpzZGVzOnJlcGFpcmVkLXJ0cC1zdHJlYW0taWRcclxuYT1zZW5kcmVjdlxyXG5hPW1zaWQ6SUV3eExWcFBiVHEwUENvdGVIaG1WaUxkUTZpUU1OOFNQVHJkIDIzZmVlMTQ2LWEyM2QtNDQ3NC05OWYwLTI2MjQxOTFmNzBmYVxyXG5hPXJ0Y3AtbXV4XHJcbmE9cnRwbWFwOjExMSBvcHVzLzQ4MDAwLzJcclxuYT1ydGNwLWZiOjExMSB0cmFuc3BvcnQtY2NcclxuYT1mbXRwOjExMSBtaW5wdGltZT0xMDt1c2VpbmJhbmRmZWM9MVxyXG5hPXJ0cG1hcDoxMDMgSVNBQy8xNjAwMFxyXG5hPXJ0cG1hcDoxMDQgSVNBQy8zMjAwMFxyXG5hPXJ0cG1hcDo5IEc3MjIvODAwMFxyXG5hPXJ0cG1hcDowIFBDTVUvODAwMFxyXG5hPXJ0cG1hcDo4IFBDTUEvODAwMFxyXG5hPXJ0cG1hcDoxMDYgQ04vMzIwMDBcclxuYT1ydHBtYXA6MTA1IENOLzE2MDAwXHJcbmE9cnRwbWFwOjEzIENOLzgwMDBcclxuYT1ydHBtYXA6MTEwIHRlbGVwaG9uZS1ldmVudC80ODAwMFxyXG5hPXJ0cG1hcDoxMTIgdGVsZXBob25lLWV2ZW50LzMyMDAwXHJcbmE9cnRwbWFwOjExMyB0ZWxlcGhvbmUtZXZlbnQvMTYwMDBcclxuYT1ydHBtYXA6MTI2IHRlbGVwaG9uZS1ldmVudC84MDAwXHJcbmE9c3NyYzoxNzQ3NzgyNjk0IGNuYW1lOnlOSjAxdHBBUDFMdkY3R09cclxuYT1zc3JjOjE3NDc3ODI2OTQgbXNpZDpJRXd4TFZwUGJUcTBQQ290ZUhobVZpTGRRNmlRTU44U1BUcmQgMjNmZWUxNDYtYTIzZC00NDc0LTk5ZjAtMjYyNDE5MWY3MGZhXHJcbmE9c3NyYzoxNzQ3NzgyNjk0IG1zbGFiZWw6SUV3eExWcFBiVHEwUENvdGVIaG1WaUxkUTZpUU1OOFNQVHJkXHJcbmE9c3NyYzoxNzQ3NzgyNjk0IGxhYmVsOjIzZmVlMTQ2LWEyM2QtNDQ3NC05OWYwLTI2MjQxOTFmNzBmYVxyXG5tPXZpZGVvIDQ4NTc2IFVEUC9UTFMvUlRQL1NBVlBGIDk2IDk3IDk4IDk5IDEwMCAxMDEgMTAyIDEyMSAxMjcgMTIwIDEyNSAxMDcgMTA4IDEwOSAxMjQgMTE5IDEyM1xyXG5jPUlOIElQNCAxNTIuMTcwLjEuMTYzXHJcbmE9cnRjcDo5IElOIElQNCAwLjAuMC4wXHJcbmE9Y2FuZGlkYXRlOjQyNTI4NzYyNTYgMSB1ZHAgMjEyMjI2MDIyMyAxOTIuMTY4LjAuMTk0IDQ4NTc2IHR5cCBob3N0IGdlbmVyYXRpb24gMCBuZXR3b3JrLWlkIDEgbmV0d29yay1jb3N0IDEwXHJcbmE9Y2FuZGlkYXRlOjMwMTk3ODQ0NjQgMSB0Y3AgMTUxODI4MDQ0NyAxOTIuMTY4LjAuMTk0IDkgdHlwIGhvc3QgdGNwdHlwZSBhY3RpdmUgZ2VuZXJhdGlvbiAwIG5ldHdvcmstaWQgMSBuZXR3b3JrLWNvc3QgMTBcclxuYT1jYW5kaWRhdGU6MjA4Mzg5NjE0OCAxIHVkcCAxNjg2MDUyNjA3IDE1Mi4xNzAuMS4xNjMgNDg1NzYgdHlwIHNyZmx4IHJhZGRyIDE5Mi4xNjguMC4xOTQgcnBvcnQgNDg1NzYgZ2VuZXJhdGlvbiAwIG5ldHdvcmstaWQgMSBuZXR3b3JrLWNvc3QgMTBcclxuYT1pY2UtdWZyYWc6di9vTlxyXG5hPWljZS1wd2Q6c2E1Ky9IZ3lYUXJpSmxHb2hyZENTT2JXXHJcbmE9aWNlLW9wdGlvbnM6dHJpY2tsZVxyXG5hPWZpbmdlcnByaW50OnNoYS0yNTYgODI6NzA6OUQ6QjY6Qjk6QjA6QzQ6QUU6MjE6ODk6NjU6MUM6Mzk6Mzk6QTQ6QkY6NkI6QkE6MzI6OEU6RTQ6Njc6ODE6MzU6RDk6NjM6RDI6RDg6RDQ6MTA6NjI6Q0JcclxuYT1zZXR1cDphY3RwYXNzXHJcbmE9bWlkOjFcclxuYT1leHRtYXA6MTQgdXJuOmlldGY6cGFyYW1zOnJ0cC1oZHJleHQ6dG9mZnNldFxyXG5hPWV4dG1hcDoyIGh0dHA6Ly93d3cud2VicnRjLm9yZy9leHBlcmltZW50cy9ydHAtaGRyZXh0L2Ficy1zZW5kLXRpbWVcclxuYT1leHRtYXA6MTMgdXJuOjNncHA6dmlkZW8tb3JpZW50YXRpb25cclxuYT1leHRtYXA6MyBodHRwOi8vd3d3LmlldGYub3JnL2lkL2RyYWZ0LWhvbG1lci1ybWNhdC10cmFuc3BvcnQtd2lkZS1jYy1leHRlbnNpb25zLTAxXHJcbmE9ZXh0bWFwOjEyIGh0dHA6Ly93d3cud2VicnRjLm9yZy9leHBlcmltZW50cy9ydHAtaGRyZXh0L3BsYXlvdXQtZGVsYXlcclxuYT1leHRtYXA6MTEgaHR0cDovL3d3dy53ZWJydGMub3JnL2V4cGVyaW1lbnRzL3J0cC1oZHJleHQvdmlkZW8tY29udGVudC10eXBlXHJcbmE9ZXh0bWFwOjcgaHR0cDovL3d3dy53ZWJydGMub3JnL2V4cGVyaW1lbnRzL3J0cC1oZHJleHQvdmlkZW8tdGltaW5nXHJcbmE9ZXh0bWFwOjggaHR0cDovL3d3dy53ZWJydGMub3JnL2V4cGVyaW1lbnRzL3J0cC1oZHJleHQvY29sb3Itc3BhY2VcclxuYT1leHRtYXA6NCB1cm46aWV0ZjpwYXJhbXM6cnRwLWhkcmV4dDpzZGVzOm1pZFxyXG5hPWV4dG1hcDo1IHVybjppZXRmOnBhcmFtczpydHAtaGRyZXh0OnNkZXM6cnRwLXN0cmVhbS1pZFxyXG5hPWV4dG1hcDo2IHVybjppZXRmOnBhcmFtczpydHAtaGRyZXh0OnNkZXM6cmVwYWlyZWQtcnRwLXN0cmVhbS1pZFxyXG5hPXNlbmRyZWN2XHJcbmE9bXNpZDpJRXd4TFZwUGJUcTBQQ290ZUhobVZpTGRRNmlRTU44U1BUcmQgMzQ4M2E5NmEtZTAzNy00ZGU1LWFkMzQtOWYyYWMzNWQ0Y2E1XHJcbmE9cnRjcC1tdXhcclxuYT1ydGNwLXJzaXplXHJcbmE9cnRwbWFwOjk2IFZQOC85MDAwMFxyXG5hPXJ0Y3AtZmI6OTYgZ29vZy1yZW1iXHJcbmE9cnRjcC1mYjo5NiB0cmFuc3BvcnQtY2NcclxuYT1ydGNwLWZiOjk2IGNjbSBmaXJcclxuYT1ydGNwLWZiOjk2IG5hY2tcclxuYT1ydGNwLWZiOjk2IG5hY2sgcGxpXHJcbmE9cnRwbWFwOjk3IHJ0eC85MDAwMFxyXG5hPWZtdHA6OTcgYXB0PTk2XHJcbmE9cnRwbWFwOjk4IFZQOS85MDAwMFxyXG5hPXJ0Y3AtZmI6OTggZ29vZy1yZW1iXHJcbmE9cnRjcC1mYjo5OCB0cmFuc3BvcnQtY2NcclxuYT1ydGNwLWZiOjk4IGNjbSBmaXJcclxuYT1ydGNwLWZiOjk4IG5hY2tcclxuYT1ydGNwLWZiOjk4IG5hY2sgcGxpXHJcbmE9Zm10cDo5OCBwcm9maWxlLWlkPTBcclxuYT1ydHBtYXA6OTkgcnR4LzkwMDAwXHJcbmE9Zm10cDo5OSBhcHQ9OThcclxuYT1ydHBtYXA6MTAwIFZQOS85MDAwMFxyXG5hPXJ0Y3AtZmI6MTAwIGdvb2ctcmVtYlxyXG5hPXJ0Y3AtZmI6MTAwIHRyYW5zcG9ydC1jY1xyXG5hPXJ0Y3AtZmI6MTAwIGNjbSBmaXJcclxuYT1ydGNwLWZiOjEwMCBuYWNrXHJcbmE9cnRjcC1mYjoxMDAgbmFjayBwbGlcclxuYT1mbXRwOjEwMCBwcm9maWxlLWlkPTJcclxuYT1ydHBtYXA6MTAxIHJ0eC85MDAwMFxyXG5hPWZtdHA6MTAxIGFwdD0xMDBcclxuYT1ydHBtYXA6MTAyIEgyNjQvOTAwMDBcclxuYT1ydGNwLWZiOjEwMiBnb29nLXJlbWJcclxuYT1ydGNwLWZiOjEwMiB0cmFuc3BvcnQtY2NcclxuYT1ydGNwLWZiOjEwMiBjY20gZmlyXHJcbmE9cnRjcC1mYjoxMDIgbmFja1xyXG5hPXJ0Y3AtZmI6MTAyIG5hY2sgcGxpXHJcbmE9Zm10cDoxMDIgbGV2ZWwtYXN5bW1ldHJ5LWFsbG93ZWQ9MTtwYWNrZXRpemF0aW9uLW1vZGU9MTtwcm9maWxlLWxldmVsLWlkPTQyMDAxZlxyXG5hPXJ0cG1hcDoxMjEgcnR4LzkwMDAwXHJcbmE9Zm10cDoxMjEgYXB0PTEwMlxyXG5hPXJ0cG1hcDoxMjcgSDI2NC85MDAwMFxyXG5hPXJ0Y3AtZmI6MTI3IGdvb2ctcmVtYlxyXG5hPXJ0Y3AtZmI6MTI3IHRyYW5zcG9ydC1jY1xyXG5hPXJ0Y3AtZmI6MTI3IGNjbSBmaXJcclxuYT1ydGNwLWZiOjEyNyBuYWNrXHJcbmE9cnRjcC1mYjoxMjcgbmFjayBwbGlcclxuYT1mbXRwOjEyNyBsZXZlbC1hc3ltbWV0cnktYWxsb3dlZD0xO3BhY2tldGl6YXRpb24tbW9kZT0wO3Byb2ZpbGUtbGV2ZWwtaWQ9NDIwMDFmXHJcbmE9cnRwbWFwOjEyMCBydHgvOTAwMDBcclxuYT1mbXRwOjEyMCBhcHQ9MTI3XHJcbmE9cnRwbWFwOjEyNSBIMjY0LzkwMDAwXHJcbmE9cnRjcC1mYjoxMjUgZ29vZy1yZW1iXHJcbmE9cnRjcC1mYjoxMjUgdHJhbnNwb3J0LWNjXHJcbmE9cnRjcC1mYjoxMjUgY2NtIGZpclxyXG5hPXJ0Y3AtZmI6MTI1IG5hY2tcclxuYT1ydGNwLWZiOjEyNSBuYWNrIHBsaVxyXG5hPWZtdHA6MTI1IGxldmVsLWFzeW1tZXRyeS1hbGxvd2VkPTE7cGFja2V0aXphdGlvbi1tb2RlPTE7cHJvZmlsZS1sZXZlbC1pZD00MmUwMWZcclxuYT1ydHBtYXA6MTA3IHJ0eC85MDAwMFxyXG5hPWZtdHA6MTA3IGFwdD0xMjVcclxuYT1ydHBtYXA6MTA4IEgyNjQvOTAwMDBcclxuYT1ydGNwLWZiOjEwOCBnb29nLXJlbWJcclxuYT1ydGNwLWZiOjEwOCB0cmFuc3BvcnQtY2NcclxuYT1ydGNwLWZiOjEwOCBjY20gZmlyXHJcbmE9cnRjcC1mYjoxMDggbmFja1xyXG5hPXJ0Y3AtZmI6MTA4IG5hY2sgcGxpXHJcbmE9Zm10cDoxMDggbGV2ZWwtYXN5bW1ldHJ5LWFsbG93ZWQ9MTtwYWNrZXRpemF0aW9uLW1vZGU9MDtwcm9maWxlLWxldmVsLWlkPTQyZTAxZlxyXG5hPXJ0cG1hcDoxMDkgcnR4LzkwMDAwXHJcbmE9Zm10cDoxMDkgYXB0PTEwOFxyXG5hPXJ0cG1hcDoxMjQgcmVkLzkwMDAwXHJcbmE9cnRwbWFwOjExOSBydHgvOTAwMDBcclxuYT1mbXRwOjExOSBhcHQ9MTI0XHJcbmE9cnRwbWFwOjEyMyB1bHBmZWMvOTAwMDBcclxuYT1zc3JjLWdyb3VwOkZJRCAyMjA2MTA0MjM4IDIyNjk0ODE1OTZcclxuYT1zc3JjOjIyMDYxMDQyMzggY25hbWU6eU5KMDF0cEFQMUx2RjdHT1xyXG5hPXNzcmM6MjIwNjEwNDIzOCBtc2lkOklFd3hMVnBQYlRxMFBDb3RlSGhtVmlMZFE2aVFNTjhTUFRyZCAzNDgzYTk2YS1lMDM3LTRkZTUtYWQzNC05ZjJhYzM1ZDRjYTVcclxuYT1zc3JjOjIyMDYxMDQyMzggbXNsYWJlbDpJRXd4TFZwUGJUcTBQQ290ZUhobVZpTGRRNmlRTU44U1BUcmRcclxuYT1zc3JjOjIyMDYxMDQyMzggbGFiZWw6MzQ4M2E5NmEtZTAzNy00ZGU1LWFkMzQtOWYyYWMzNWQ0Y2E1XHJcbmE9c3NyYzoyMjY5NDgxNTk2IGNuYW1lOnlOSjAxdHBBUDFMdkY3R09cclxuYT1zc3JjOjIyNjk0ODE1OTYgbXNpZDpJRXd4TFZwUGJUcTBQQ290ZUhobVZpTGRRNmlRTU44U1BUcmQgMzQ4M2E5NmEtZTAzNy00ZGU1LWFkMzQtOWYyYWMzNWQ0Y2E1XHJcbmE9c3NyYzoyMjY5NDgxNTk2IG1zbGFiZWw6SUV3eExWcFBiVHEwUENvdGVIaG1WaUxkUTZpUU1OOFNQVHJkXHJcbmE9c3NyYzoyMjY5NDgxNTk2IGxhYmVsOjM0ODNhOTZhLWUwMzctNGRlNS1hZDM0LTlmMmFjMzVkNGNhNVxyXG4ifQ=="
		
		cmd := exec.Command("/src/github.com/pion/webrtc/examples/rtp-forwarder/jsfiddle/script.sh",token)
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()
		oneByte := make([]byte, 1)
		num := 1
		for {
			_, err := stdout.Read(oneByte)
			if err != nil {
				fmt.Printf(err.Error())
				break
			}
			r := bufio.NewReader(stdout)
			line, _, _ := r.ReadLine()
			fmt.Println(string(line))
			num = num + 1
			if num == 3 {
				fmt.Fprintln(w,"e"+string(line))
				return
			}
		}
		cmd.Wait()
		return
	})

	// Serve the required pages
	// DIY 'mux' to avoid additional dependencies
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		if url == "/wasm_exec.js" {
			http.FileServer(http.Dir("./vendor-wasm/golang.org/misc/wasm/")).ServeHTTP(w, r)
			return
		}

		// Split up the URL. Expected parts:
		// 1: Base url
		// 2: "example"
		// 3: Example type: js or wasm
		// 4: Example folder, e.g.: data-channels
		// 5: Static file as part of the example
		parts := strings.Split(url, "/")
		if len(parts) > 4 &&
			parts[1] == "example" {
			exampleType := parts[2]
			exampleLink := parts[3]
			for _, example := range *examples {
				if example.Link != exampleLink {
					continue
				}
				fiddle := filepath.Join(exampleLink, "jsfiddle")
				if len(parts[4]) != 0 {
					http.StripPrefix("/example/"+exampleType+"/"+exampleLink+"/", http.FileServer(http.Dir(fiddle))).ServeHTTP(w, r)
					return
				}

				temp := template.Must(template.ParseFiles("example.html"))
				_, err = temp.ParseFiles(filepath.Join(fiddle, "demo.html"))
				if err != nil {
					panic(err)
				}

				data := struct {
					*Example
					JS bool
				}{
					example,
					exampleType == "js",
				}

				err = temp.Execute(w, data)
				if err != nil {
					panic(err)
				}
				return
			}
		}	

		// Serve the main page
		err = homeTemplate.Execute(w, examples)
		if err != nil {
			panic(err)
		}
	})

	// Start the server
	return http.ListenAndServe(addr, nil)
}

// getExamples loads the examples from the examples.json file.
func getExamples() (*Examples, error) {
	file, err := os.Open("./examples.json")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errListExamples, err)
	}
	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			panic(closeErr)
		}
	}()

	var examples Examples
	err = json.NewDecoder(file).Decode(&examples)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errParseExamples, err)
	}

	for _, example := range examples {
		fiddle := filepath.Join(example.Link, "jsfiddle")
		js := filepath.Join(fiddle, "demo.js")
		if _, err := os.Stat(js); !os.IsNotExist(err) {
			example.IsJS = true
		}
		wasm := filepath.Join(fiddle, "demo.wasm")
		if _, err := os.Stat(wasm); !os.IsNotExist(err) {
			example.IsWASM = true
		}
	}

	return &examples, nil
}
