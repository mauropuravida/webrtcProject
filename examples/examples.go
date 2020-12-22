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
	"project/webrtcProject/examples/db"
	"strconv"
	"project/webrtcProject/examples/models"


)

const (
    PUBLIC_KEY = "/etc/letsencrypt/live/testingwebrtc.ddns.net/fullchain.pem"
    PRIV_KEY = "/etc/letsencrypt/live/testingwebrtc.ddns.net/privkey.pem"
)

// Examples represents the examples loaded from examples.json.
type Examples []*Example

var (
	errListExamples  = errors.New("failed to list examples (please run in the examples folder)")
	errParseExamples = errors.New("failed to parse examples")
	token_stream = ""
	token_connect = ""
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

		//Reset vars connection for test
		token_stream = ""
		token_connect = ""

		//user := req.Form.Get("user")
		//id_camera := req.Form.Get("id_camera")
		token := req.Form.Get("token")

		if token_stream == ""{
			token_stream = token
		}

		// await for token_connect
		endTime := time.Now().Add(time.Second * 61)
		for time.Now().Before(endTime) {
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
		token_connect = req.Form.Get("token")

		token := token_stream
		if token != ""{
			fmt.Fprintln(w,token_connect)
			return
		}

		fmt.Printf(token_connect)
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

	//Add camera
	http.HandleFunc("/addcamera", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		user := req.Form.Get("user")
		loc := req.Form.Get("loc")
		var idcam int64
		
		user_id, err:= strconv.Atoi(user)
		if err == nil {
			idcam,err=db.InsertCam(user_id,loc)
		}

		
		
		fmt.Fprintln(w, idcam)
		return
		
	})

	//update camera
	http.HandleFunc("/updatecamera", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		user := req.Form.Get("user")
		loc := req.Form.Get("loc")
		idCam := req.Form.Get("id_camera")

		fmt.Println("AUPDATE");
		fmt.Println(loc);
		user_id, err:= strconv.Atoi(user)
		cam_id, err:= strconv.Atoi(idCam)
		if err == nil {
		fmt.Println("abase")
			db.UpdateCam(cam_id,user_id,loc)
		}
		return
		
	})

	//Delete camera
	http.HandleFunc("/deletecamera", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		idCam := req.Form.Get("id_camera")
		fmt.Println("idCam")
		fmt.Println(idCam)

		cam_id, err:= strconv.Atoi(idCam)
		if err == nil {
		fmt.Println("ESTA BORRANDO EN BASE")
			db.DeleteCam(cam_id)
		}

		if err != nil {
		fmt.Println("Error: ")
		fmt.Println(err)
			
		}
			
		fmt.Printf("idcam: "+idCam)

		return
		
	})

	//Active/desactive camera
	http.HandleFunc("/activecamera", func(w http.ResponseWriter, req *http.Request) {
		// Parses the request body
		req.ParseForm()
		id := req.Form.Get("id_camera")
		cam_id, err:= strconv.Atoi(id)
		if err == nil {
			db.UpdateActiveCam(true,cam_id)
		}
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



	// get cams from user 
	http.HandleFunc("/getCameras", func(w http.ResponseWriter, req *http.Request) {
		
		req.ParseForm()		// Parses the request body
		id := req.Form.Get("id")
		
		cams := make([]models.Camera, 0)
		if err!=nil {
			return
		}

		user_id, err:= strconv.Atoi(id)
		if err == nil {
			cams, err= db.GetCamsByUser(user_id)
		}
		if err!=nil {
			fmt.Println("ERROR")
			fmt.Println(err)
			return
		}
		

		data, _ := json.Marshal(cams)
		

		fmt.Fprintln(w, string(data))
		
		return
	})



	//Generate session token
	http.HandleFunc("/gentoken", func(w http.ResponseWriter, req *http.Request) {

		req.ParseForm()		// Parses the request body
		token := req.Form.Get("token") // x will be "" if parameter is not set
		
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
	if  addr != ":80" {
		ip := strings.Split(addr, ":")[0]
		fmt.Printf(ip+" CONECTING 443")
		go http.ListenAndServe( ip+":80", nil)
		return http.ListenAndServeTLS( ip+":443", PUBLIC_KEY, PRIV_KEY, http.HandlerFunc(redirect443to80))
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


