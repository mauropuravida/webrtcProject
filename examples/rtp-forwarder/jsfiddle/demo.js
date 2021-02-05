
/* eslint-env browser */
var host = "http://localhost:80"
//var host = "http://159.65.97.50"

var pcMap = new Map()
var xhrMap = new Map()
var currentUser=document.getElementById("user-id").innerHTML;
var numberOfCam=1;

function addCamera(url, description, active, camId) {
    
  var table = document.getElementById("cameras");
  var row = document.createElement("tr");
  table.appendChild(row);
  row.id = 'rowcam'+camId;

  var cell1 = row.insertCell(0);
  var element1= document.createElement("input");
  element1.type="text";

  if (url == "")
    url = "http://192.168.0.234:4747/video"

  element1.value= url
  element1.id = "url"+camId
  cell1.appendChild(element1);

  var cell2 = row.insertCell(1);
  var element2= document.createElement("input");
  element2.type="text";
  element2.value = description
  element2.id = "description"+camId
  cell2.appendChild(element2);

  var cell3 = row.insertCell(2)
  cell3.width = 180
  cell3.height = 180
  cell3.style.textAlign = "center"

  var element3 = document.createElement("canvas");
  element3.width = 160
  element3.height = 120
  element3.id = "canvas"+camId
  cell3.appendChild(element3)

  var cell4 = row.insertCell(3);
  var element41 = document.createElement("input");
  element41.type="button";
  element41.value = "Modo Gris";
  element41.id = "greyMode" + camId;

  var element42 = document.createElement("input");
  element42.type="button";
  element42.value = "Rotar";
  element42.id = "rotateStream" + camId;

  var element43 = document.createElement("input");
  element43.type="button";
  element43.value = "Conectar";
  element43.id = "ipcamera" + camId;

  var element44 = document.createElement("input");
  element44.type="button";
  element44.value = "Desconectar";
  element44.id = "stopStream" + camId;

  var element45 = document.createElement("input");
  element45.type="button";
  element45.value = "Eliminar"
  element45.id = "remove" + camId

  var element46 = document.createElement("input");
  element46.type = "button";
  element46.value = "Guardar";
  element46.id = "save" + camId;

  cell4.appendChild(element41)
  cell4.appendChild(element42)
  cell4.appendChild(element43)
  cell4.appendChild(element44)
  cell4.appendChild(element45)
  cell4.appendChild(element46)

  var privateCamId = camId
  document.getElementById("ipcamera"+camId).addEventListener('click', function(){ connectStream(privateCamId)})

  document.getElementById("remove" + camId).addEventListener('click', function () {
    deleteCam(privateCamId);

});

document.getElementById("save" + camId).addEventListener('click', function () {
    let newURL = document.getElementById("url" + camId).value;
    let newDesc = document.getElementById("description" + camId).value;

    //The problem to use update method is that the insert is done when you press the first time "Guardar", to use update is necessary change the logic
    //updateCam(privateCamId,newURL, newDesc, active);
    deleteCam(privateCamId);
    insertCam(newURL, newDesc, active);

});

document.getElementById('stopStream' + camId).addEventListener('click', function () {
    var xhr = new XMLHttpRequest();
    xhr.withCredentials = true;
    xhr.open("POST", host + "/activecamera");
    data = "id_camera=" + camId + "&active=" + 0;
    console.log("data: " + data);
    xhr.setRequestHeader("cache-control", "no-cache");
    xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhr.send(data);


    if (pcMap.has(camId))
        pcMap.get(camId).close()


    if (xhrMap.has(camId)) {
        console.log("abort " + camId)
        xhrMap.get(camId).abort()
        xhrMap.delete(camId)
    }
});

   
}


function insertCam(newURL, newDesc,status) {
    
    var xhr1 = new XMLHttpRequest();
    xhr1.onload = function () {
        console.log(this.readyState);
        if (this.readyState === 4) {
            addCamera(newURL, newDesc, status, this.responseText);
        }
    };
    xhr1.open("POST", host + "/addcamera");
    data = "user=" + currentUser +"&idcam="+ numberOfCam +"&loc=" + newDesc + "&url=" + newURL;
    xhr1.setRequestHeader("cache-control", "no-cache");
    xhr1.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhr1.send(data);
    numberOfCam++;
}

function updateCam(numberOfCam, newURL, newDesc,status) {
    
    var xhr1 = new XMLHttpRequest();
    xhr1.onload = function () {
        console.log(this.readyState);
        if (this.readyState === 4) {
            addCamera(newURL, newDesc, status, numberOfCam);
        }
    };
    xhr1.open("POST", host + "/updatecamera");
    data = "user=" + currentUser +"&idcam="+ numberOfCam +"&loc=" + newDesc + "&url=" + newURL+"&active=" + status;
    xhr1.setRequestHeader("cache-control", "no-cache");
    xhr1.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhr1.send(data);
    numberOfCam++;
}

function deleteCam(camId) {
    var row = document.getElementById('rowcam' + camId);
    row.parentNode.removeChild(row);
    var xhr = new XMLHttpRequest();
    xhr.withCredentials = true;
    xhr.open("POST", host + "/deletecamera");
    data = "id_camera=" + camId+ "&id_user="+ currentUser;
    xhr.setRequestHeader("cache-control", "no-cache");
    xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhr.send(data);
 
}


function connectStream(id) {
    console.log("CONECTANDO " + id)

    if (document.getElementById('url' + id).value == "") {
        console.log("URL no valida")
        return
    }

    var myImg = new Image()
    myImg.src = document.getElementById('url' + id).value
    myImg.crossOrigin = "Anonymous"

    var canvasClient = document.getElementById('canvas' + id);

    var canvasStream = document.createElement("CANVAS");
    canvasStream.width = 480
    canvasStream.height = 360

    var xhr = new XMLHttpRequest();
    xhr.withCredentials = true;
    xhr.open("POST", host + "/activecamera");
    data = "id_camera=" + id + "&active=" + 1;
    console.log("data: " + data);
    xhr.setRequestHeader("cache-control", "no-cache");
    xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhr.send(data);


    let rotateAngle = 0;
    let xrC = canvasClient.width;
    let yrC = canvasClient.height;
    let xrS = canvasStream.width;
    let yrS = canvasStream.height;
    function rotate() {
        rotateAngle = rotateAngle + 90;
        angle = parseInt(rotateAngle) * Math.PI / 180;

        switch (rotateAngle) {
            case 90: case 270:
                yrC = -yrC
                yrS = -yrS
                canvasClient.width = 120
                canvasClient.height = 160
                break;
            case 180: case 360:
                xrC = -xrC
                xrS = -xrS
                canvasClient.width = 160
                canvasClient.height = 120

                if (rotateAngle == 360)
                    rotateAngle = 0
                break;
            default:
            // code block
        }
    }

    var rotateStream = document.getElementById('rotateStream' + id)
    rotateStream.addEventListener('click', rotate, false)

    var transform = false;
    function changeColorStream() {
        if (transform)
            transform = false
        else
            transform = true
    }

    var greyMode = document.getElementById('greyMode' + id)
    greyMode.addEventListener('click', changeColorStream, false)

    function handlerImg(element, xr, yr) {
        var context = element.getContext("2d");

        context.clearRect(0, 0, element.width, element.height);
        context.rotate(angle);
        context.drawImage(myImg, 0, 0, xr, yr);

        //transform grey color
        if (transform) {
            var imgData = context.getImageData(0, 0, element.width, element.height);
            var px = imgData.data;
            for (i = 0; i < imgData.data.length; i += 4) {
                var grey = px[i] * .3 + px[i + 1] * .59 + px[i + 2] * .11;
                px[i] = px[i + 1] = px[i + 2] = grey;
            }

            context.putImageData(imgData, 0, 0);
        }
    }

    var angle = 0;



    (function loop() {
        if (!this.paused && !this.ended && (!pcMap.has(id) || pcMap.get(id).connectionState != 'closed')) {
            //draw in canvas client
            handlerImg(canvasClient, xrC, yrC);
            //draw in canvas for stream best resolution
            handlerImg(canvasStream, xrS, yrS);
        } else {
            var context = canvasClient.getContext("2d");
            context.clearRect(0, 0, canvasClient.width, canvasClient.height)
            pcMap.delete(id)
            myImg.src = "#"
            return
        }

        setTimeout(loop, 1000 / 60); // drawing at 60fps
    })();

    var stream = canvasStream.captureStream(60)

    pc = new RTCPeerConnection({
        iceServers: [
            {
                urls: 'stun:stun.l.google.com:19302'
            }
        ]
    })

    pcMap.set(id, pc);
    listenPc(id)

    pc.addStream(stream)
    pc.createOffer().then(function (d) {
        pc.setLocalDescription(d)
    }).catch("error")

    
 }

function listenPc(id){
  var pc = pcMap.get(id)
  pc.oniceconnectionstatechange = function(e) { console.log(pc.iceConnectionState) }
  pc.onicecandidate = function(event) {
    if (event.candidate === null) {    
        //var localSesion = btoa(JSON.stringify(pc.localDescription));
        var localSesion;  
        getTokenCamFromDB(id, pc, localSesion);

      var data = "user="+currentUser+"&id_camera="+id+"&token="+localSesion;
      console.log(data);
      getToken(data, id);
    }
  }
}

function getToken(data, id){
  var xhr = new XMLHttpRequest();

  xhrMap.set(id, xhr)
  xhr.withCredentials = true;
  xhr.onload = function () {
  console.log(this.readyState);
    if (this.readyState === 4) {
      //remoteSesion = this.responseText
      setTimeout(window.startSession(id, this.responseText), 3000);
      console.log(this.responseText);
    }
  };
  xhr.ontimeout = function(e){
    console.log("time lost");
    getToken("");
  }
  // await one minute for response   
    //primero
    //mandar id user
  xhr.timeout = 60000;
  xhr.open("POST", host+"/sendtokenstreamer");
  xhr.setRequestHeader("cache-control", "no-cache");
  xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
  xhr.setRequestHeader('token', 'eltoken');
  xhr.send(data);
}

function login() {
  var xhr = new XMLHttpRequest();

  currentUser = document.getElementById("user-id").value;
  
    xhr.onload = function () {
        console.log(this.readyState);
        if (this.readyState === 4) {
            console.log(this.responseText);
            var data = JSON.parse(this.responseText);
            document.getElementById("cameras").innerHTML= " " ;
            for (var key in data) {
                addCamera(data[key].Url, data[key].Loc, data[key].Active, data[key].Id_cam);
            }
            getNextCamIdByUser(currentUser);

        }
    };


    data = 'id='+ currentUser;
    xhr.withCredentials = true;
    xhr.open("POST", host + "/getCameras");
    xhr.setRequestHeader("cache-control", "no-cache");
    xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhr.setRequestHeader('token', 'eltoken');
    xhr.send(data);
}

function getNextCamIdByUser(currentUser){
    var xhr = new XMLHttpRequest();

    xhr.onload = function () {
        console.log(this.readyState);
        if (this.readyState === 4) {
            numberOfCam=parseInt(this.responseText);
        }
    };
    data = 'id='+ currentUser;
    xhr.withCredentials = true;
    xhr.open("POST", host + "/getNextCamIdByUser");
    xhr.setRequestHeader("cache-control", "no-cache");
    xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhr.setRequestHeader('token', 'eltoken');
    xhr.send(data);

    
}

function getTokenCamFromDB(id, pc, localSesion){
    var xhr = new XMLHttpRequest();

    xhr.onload = function () {
        console.log(this.responseText);
        if (this.readyState === 4) {
            if(this.responseText.length<100){
                localSesion = btoa(JSON.stringify(pc.localDescription));
                saveTokenCam(id,localSesion);
            }
            else{
                localSesion = this.responseText;

            } 
        }
    };
    data = 'id='+ id + '&user='+currentUser;
    xhr.withCredentials = true;
    xhr.open("POST", host + "/getTokenCamFromDB");
    xhr.setRequestHeader("cache-control", "no-cache");
    xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhr.setRequestHeader('token', 'eltoken');
    xhr.send(data);
}

function saveTokenCam(id, localSesion){
    var xhr = new XMLHttpRequest();
    data = 'id='+ id + '&user='+currentUser+"&token="+localSesion;
    xhr.withCredentials = true;
    xhr.open("POST", host + "/saveTokenCam");
    xhr.setRequestHeader("cache-control", "no-cache");
    xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
    xhr.setRequestHeader('token', 'eltoken');
    xhr.send(data);
}


window.startSession = function(id, remoteSesion) {
  try {
    pcMap.get(id).setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(remoteSesion))))
  } catch (e) {
    alert(e)
  }
}