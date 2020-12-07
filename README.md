# webrtcProject

## Requirements
**Install go**
- Install go 1.15 https://tecadmin.net/install-go-on-linuxmint/

**Import env vars**
    ```bash
        export GOROOT=/usr/local/go ; export GOPATH=$HOME/Apps/app1 ; export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
    ```

**Note:** This project use port 80 and 443


## Run project
 Go to root project folder (../webrtcProject/examples) and execute this command:

**Localhost**
    ```bash
        go run  examples.go &
    ```

**For remote server**
    ```bash
        go run  examples.go --address <ADDRESS> &
    ```
    Edit ../webrtcProject/examples/rtp-forwarder/jsfiddle/demo.js for specify your host.


## How to use:

<HOST>: Reference your host
<ADDRES>: Local listen address for stream. If not specify address, you listen stream on 127.0.0.200 by default.
<PORT>: Local port for listen stream. If not specify port, use port 4000 by default.


1. Go to url http://<HOST>/example/js/rtp-forwarder/
2. Add and connect your IP-CAM
3. For consume stream use this command:
    ```bash
        -curl -X GET http://<HOST>/checkstream | .<PATH_TO>/rtp-forwarder --address <ADDRESS> --port <PORT> --host <HOST>
    ```
4. Go to the tmp folder and choose the generated .sdp file. Open it with a media player.


## Useful commands for develop

//List conections
```bash

    sudo lsof -i -P -n | grep LISTEN
```

//Kill process
```bash

    sudo kill -9 PROCESSID
```

