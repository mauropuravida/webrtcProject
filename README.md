# webrtcProject


**Install go**
- Install go 1.15 https://tecadmin.net/install-go-on-linuxmint/

**Import env vars**
    ```bash
        export GOROOT=/usr/local/go ; export GOPATH=$HOME/Apps/app1 ; export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
    ```

**Run project**
 Go to root project folder and execute this command:
```bash

    go run  examples.go &
```



## Examples:

//Get current token ssesion streamer and connect.
-curl -X GET http://localhost/checkstream | ./rtp-forwarder --address 127.0.0.10 --port 4040


//Reset vars for testing
-curl -X GET http://localhost/reset

//Send token for connect, this action is done inside rtp-forwarder.
-curl -X POST https://localhost/sendtokenconnect

**Note:** content type x-wwww-form-urlencoded, in body declare "token" and your value.


 ## Useful commands

//List conections
```bash

    sudo lsof -i -P -n | grep LISTEN
```

//Kill process
```bash

    sudo kill -9 PROCESSID
```



**Note:** For deploy in local host you need change:

- in "webrtcProject/examples/rtp-forwarder/jsfiddle/demo.js" comment line 22 and uncomment line 21
- in "webrtcProject/examples/example.go" comment line 4 and uncomment line 3

