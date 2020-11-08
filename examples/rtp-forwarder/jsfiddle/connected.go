package main

import (
	"fmt"
	"flag"
	"bufio"
	"os/exec"
)

func main() {
	token := flag.String("token","","")
	flag.Parse()

	fmt.Println(token)
	//cd /src/github.com/pion/webrtc/examples/rtp-forwarder
	//curl -X GET https://testingwebrtc.ml/checkstream | ./rtp-forwarder | tail -n +3

	//req.ParseForm()		// Parses the request body
	//token := "sdasfasf"

	cmd := exec.Command("/src/github.com/pion/webrtc/examples/rtp-forwarder/jsfiddle/script.sh", "string(token)")
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
			//fmt.Fprintln(w,"e"+string(line))
			return
		}
	}
	cmd.Wait()
	return
}