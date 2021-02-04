# webrtcProject

## Requirements
**Install go**
- Install go 1.15 https://tecadmin.net/install-go-on-linuxmint/
- Install docker
- Pull this image:
    ```bash
        docker pull mysql
    ```
- Create and run container: 
    ```bash
        docker run -d -p 3306:3306 --name dbcam -e MYSQL_ROOT_PASSWORD=secret mysql
    ```
	If you cant use this port, you can change it and change in /examples/db/Connection.go the information about connection.

- Enter the container. Use password "secret"
    ```bash
        docker exec -it dbcam mysql -p
    ```

- Create the schema using the script in examples/db/dbCam.sql. Copy this file content inside the container.



**Import env vars**
    ```bash
        export GOROOT=/usr/local/go ; export GOPATH=$HOME/Apps/app1 ; export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
    ```

**Import new packages**
    Copy ../webrtcProject/examples/db and ../webrtcProject/examples/models in /usr/local/go/src folder.

**Note:** This project use port 8080 and 443


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
<IDCAMERA>: is the number in id_camera field in "Cameras" table. 
<USER>: is the number in users_id field in "Cameras" table.

Use <IDCAMERA> and <USER> to identify an specific camera.


1. Go to url http://<HOST>/example/js/rtp-forwarder/
2. Add and connect your IP-CAM
3. For consume stream use this command:
    ```bash
        curl -X POST -d "cam=<IDCAMERA>&user=<USER>" http://<HOST>/checkstream | .<PATH_TO>/rtp-forwarder --address <ADDRESS> --port <PORT> --host <HOST> --idCam <IDCAMERA> --idUser <USER>
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

