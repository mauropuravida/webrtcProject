package main

import (
	"db"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"models"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	PUBLIC_KEY = "/etc/letsencrypt/live/testingwebrtc.ddns.net/fullchain.pem"
	PRIV_KEY   = "/etc/letsencrypt/live/testingwebrtc.ddns.net/privkey.pem"
)

// Examples represents the examples loaded from examples.json.
type Examples []*Example

var (
	errListExamples  = errors.New("failed to list examples (please run in the examples folder)")
	errParseExamples = errors.New("failed to parse examples")
	token_stream     = ""
	token_connect    = ""
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

//redirect request poty 80 to port 443
func redirect(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "https://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		http.StatusTemporaryRedirect)
}

//redirect request poty 443 to port 80
func redirect443to80(w http.ResponseWriter, req *http.Request) {
	// remove/add not default ports from req.Host
	target := "http://" + req.Host + req.URL.Path
	if len(req.URL.RawQuery) > 0 {
		target += "?" + req.URL.RawQuery
	}
	log.Printf("redirect to: %s", target)
	http.Redirect(w, req, target,
		http.StatusTemporaryRedirect)
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

		id_cam := req.Form.Get("id")
		user := req.Form.Get("user")
		token := req.Form.Get("token")
		cam_id, err := strconv.Atoi(id_cam)
		user_id, err := strconv.Atoi(user)

		//Reset vars connection for test
		token_stream = ""
		token_connect = ""

		if err == nil {
			db.UpdateTokenCam(cam_id, user_id, token)
		} else {
			fmt.Println(err)
		}

		if token_stream == "" {
			token_stream = token
		}

		// await for token_connect
		endTime := time.Now().Add(time.Second * 61)
		for time.Now().Before(endTime) { //listening response
			//consultar que la base no devuelva vacio
			if token_connect != "" {
				tokenResponse := token_connect
				token_connect = ""
				fmt.Fprintln(w, tokenResponse)
				return
			}
		}

		fmt.Fprintln(w, "")
		return
	})

	//Receiver token from consumer for connect to streamer
	http.HandleFunc("/sendtokenconnect", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		token_connect := req.Form.Get("token")
		id := req.Form.Get("id_camera")
		user := req.Form.Get("id_user")
		cam_id, err := strconv.Atoi(id)
		user_id, err := strconv.Atoi(user)

		if err == nil {
			db.UpdateTokenCon(cam_id, token_connect, user_id)
		}

		if token_connect != "" {
			fmt.Printf(token_connect)
			fmt.Fprintln(w, token_connect)
			return
		}

		fmt.Printf(" token_connect not found")
		return
	})

	//Delete vars only for *testing*
	http.HandleFunc("/reset", func(w http.ResponseWriter, req *http.Request) {
		token_stream = ""
		token_connect = ""

		fmt.Fprintln(w, "reset")
		return
	})

	//Check if exists stream from user and id_camera
	http.HandleFunc("/checkstream", func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		idCam := req.Form.Get("cam")
		user := req.Form.Get("user")

		user_id, err := strconv.Atoi(user)

		cam_id, err := strconv.Atoi(idCam)

		var token_stream string

		if err == nil {
			token_stream = db.GetTokenCam(cam_id, user_id)

		}

		if token_stream != "" {
			fmt.Fprintln(w, token_stream)
			return
		}

		fmt.Fprintln(w, "token not found")
		return
	})

	//Add camera
	http.HandleFunc("/addcamera", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		user := req.Form.Get("user")
		loc := req.Form.Get("loc")
		url := req.Form.Get("url")
		idcam := req.Form.Get("idcam")

		user_id, err := strconv.Atoi(user)

		cam_id, err := strconv.Atoi(idcam)

		if err != nil {
			fmt.Println(err)
		}

		if err == nil {
			err = db.InsertCam(user_id, loc, url, cam_id)

			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintln(w, "Cam not saved")
				return
			}
		}

		return

	})

	//update camera
	http.HandleFunc("/updatecamera", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		user := req.Form.Get("user")
		loc := req.Form.Get("loc")
		idCam := req.Form.Get("idcam")
		url := req.Form.Get("url")
		active := req.Form.Get("active")
		tokcam := req.Form.Get("tokcam")
		tokcon := req.Form.Get("tokcon")

		act, err := strconv.ParseBool(active)
		//fmt.Println(loc);
		user_id, err := strconv.Atoi(user)
		cam_id, err := strconv.Atoi(idCam)
		if err == nil {
			db.UpdateCam(cam_id, user_id, loc, url, act, tokcam, tokcon)
		} else {
			fmt.Println(err)
		}
		return

	})

	//Delete camera
	http.HandleFunc("/deletecamera", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		idCam := req.Form.Get("id_camera")
		idUser := req.Form.Get("id_user")

		cam_id, err := strconv.Atoi(idCam)
		usr_id, err := strconv.Atoi(idUser)

		if err == nil {
			db.DeleteCam(cam_id, usr_id)
		}

		if err != nil {
			fmt.Println(err)
		}

		return

	})

	//Active/desactive camera
	http.HandleFunc("/activecamera", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		id := req.Form.Get("id_camera")
		active := req.Form.Get("active")
		cam_id, err := strconv.Atoi(id)
		act, err := strconv.ParseBool(active)
		if err == nil {
			db.UpdateActiveCam(act, cam_id)
		} else {
			fmt.Println(err)
		}
		return

	})

	//Login user
	http.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		//email := req.Form.Get("email")
		//loc := req.Form.Get("password")

		//fmt.Printf(email+" "+loc)
		return
	})

	// get cams from user
	http.HandleFunc("/getCameras", func(w http.ResponseWriter, req *http.Request) {

		req.ParseForm() // Parses the request body
		id := req.Form.Get("id")

		cams := make([]models.Camera, 0)
		if err != nil {
			return
		}

		user_id, err := strconv.Atoi(id)
		if err == nil {
			cams, err = db.GetCamsByUser(user_id)
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		data, _ := json.Marshal(cams)

		fmt.Fprintln(w, string(data))

		return
	})

	//get the next id for cams from an specific user
	http.HandleFunc("/getNextCamIdByUser", func(w http.ResponseWriter, req *http.Request) {

		req.ParseForm() // Parses the request body
		id := req.Form.Get("id")

		if err != nil {
			return
		}

		user_id, err := strconv.Atoi(id)

		var nextId int
		if err == nil {
			nextId, err = db.GetNextCamIdByUser(user_id)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintln(w, nextId)

		return
	})

	http.HandleFunc("/getTokenCamFromDB", func(w http.ResponseWriter, req *http.Request) {

		req.ParseForm() // Parses the request body
		id := req.Form.Get("id")
		user := req.Form.Get("user")

		if err != nil {
			return
		}

		user_id, err := strconv.Atoi(user)
		cam_id, err := strconv.Atoi(id)

		var token string
		if err == nil {
			token = db.GetTokenCam(cam_id, user_id)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Fprintln(w, token)

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
	if addr != ":80" {
		ip := strings.Split(addr, ":")[0]
		fmt.Printf(ip + " CONECTING 443")
		go http.ListenAndServe(ip+":80", nil)
		return http.ListenAndServeTLS(ip+":443", PUBLIC_KEY, PRIV_KEY, http.HandlerFunc(redirect443to80))
	}
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
