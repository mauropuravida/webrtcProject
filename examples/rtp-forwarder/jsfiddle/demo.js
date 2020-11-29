
/* eslint-env browser */
var host = "http://localhost"
//var host = "https://testingwebrtc.ddns.net"

var pcMap = new Map()
var xhrMap = new Map()

let camId = 5

function addCamera(url, description){
  var table = document.getElementById("cameras");
  var row = table.insertRow();
  row.id = 'rowcam'+camId

  var cell1 = row.insertCell(0);
  var element1= document.createElement("input");
  element1.type="text";
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
  element41.value = "Modo Gris"
  element41.id = "greyMode"+camId

  var element42 = document.createElement("input");
  element42.type="button";
  element42.value = "Rotar"
  element42.id = "rotateStream"+camId

  var element43 = document.createElement("input");
  element43.type="button";
  element43.value = "Conectar"
  element43.id = "ipcamera"+camId

  var element44 = document.createElement("input");
  element44.type="button";
  element44.value = "Desconectar"
  element44.id = "stopStream"+camId

  var element45 = document.createElement("input");
  element45.type="button";
  element45.value = "Eliminar"
  element45.id = "remove"+camId

  cell4.appendChild(element41)
  cell4.appendChild(element42)
  cell4.appendChild(element43)
  cell4.appendChild(element44)
  cell4.appendChild(element45)

  var privateCamId = camId
  document.getElementById("ipcamera"+camId).addEventListener('click', function(){ connectStream(privateCamId)})
  document.getElementById("remove"+camId).addEventListener('click', function(){
    var row = document.getElementById('rowcam'+privateCamId);
    row.parentNode.removeChild(row);
  })

  camId = camId + 1
}

function connectStream(id){
  console.log("CONECTANDO "+id)

  if (document.getElementById('url'+id).value == ""){
    console.log("URL no valida")
    return
  }

  var myImg = new Image()
  myImg.src = document.getElementById('url'+id).value
  myImg.crossOrigin = "Anonymous"

  var canvasClient = document.getElementById('canvas'+id);

  var canvasStream = document.createElement("CANVAS");
  canvasStream.width = 480
  canvasStream.height = 360

  let rotateAngle = 0;
  let xrC = canvasClient.width;
  let yrC = canvasClient.height;
  let xrS = canvasStream.width;
  let yrS = canvasStream.height;
  function rotate(){
    rotateAngle = rotateAngle + 90;
    angle = parseInt(rotateAngle) * Math.PI / 180;

    switch(rotateAngle) {
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

  var rotateStream = document.getElementById('rotateStream'+id)
  rotateStream.addEventListener('click', rotate, false)

  var transform = false;
  function changeColorStream(){
    if(transform)
      transform = false
    else
      transform = true
  }

  var greyMode = document.getElementById('greyMode'+id)
  greyMode.addEventListener('click', changeColorStream, false)

  function handlerImg(element, xr, yr){
    var context = element.getContext("2d");

    context.clearRect(0,0,element.width, element.height);  
    context.rotate(angle);
    context.drawImage(myImg, 0, 0, xr, yr);

    //transform grey color
    if (transform){
      var imgData = context.getImageData(0, 0, element.width, element.height);
      var px = imgData.data;
      for (i = 0; i < imgData.data.length; i += 4) {
        var grey = px[i] * .3 + px[i+1] * .59 + px[i+2] * .11;
        px[i] = px[i+1] = px[i+2] = grey;
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
    }else{
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

  pcMap.set(id, pc)
  listenPc(id)

  pc.addStream(stream)
  pc.createOffer().then(d => {
    pc.setLocalDescription(d)
  }).catch("error")

  //obtain element
  var stopStream = document.getElementById('stopStream'+id)

  //stop stream
  function stop(){
    if (pcMap.has(id))
      pcMap.get(id).close()


    if(xhrMap.has(id)){
      console.log("abort "+id)
      xhrMap.get(id).abort()
      xhrMap.delete(id)
    }
  }

  stopStream.addEventListener('click', stop, false)

  var remove = document.getElementById('remove'+id)

  remove.addEventListener('click', stop, false)

}

function connectStreamNavigator(id){
  navigator.mediaDevices.enumerateDevices()
  .then(devices => {
    var camera = devices.find(device => device.kind == "videoinput" && device.label == 'HP Wide Vision FHD Camera (0bda:58e6)');
    if (camera) {
      var constraints = { deviceId: { exact: camera.deviceId } };
      return navigator.mediaDevices.getUserMedia({ video: constraints });
    }
  })
  .then(stream => {
      pc = new RTCPeerConnection({
      iceServers: [
        {
          urls: 'stun:stun.l.google.com:19302'
        }
        ]
      })
      pcMap.set(id, pc)
      listenPc(id)

      pc.addStream(document.getElementById('localcamera').srcObject = stream)
      pc.createOffer().then(d => pc.setLocalDescription(d)).catch(log)
    })
  .catch(e => console.error(e));
}

function stopStream(id){
  pcMap.get(id).close()

  if (id == 'localcamera'){
    document.getElementById('localcamera').srcObject = null
    pcMap.delete(id)
  }

  if(xhrMap.has(id)){
    console.log("abort "+id)
    xhrMap.get(id).abort()
    xhrMap.delete(id)
  }
}

function listenPc(id){
  var pc = pcMap.get(id)
  pc.oniceconnectionstatechange = e => console.log(pc.iceConnectionState)
  pc.onicecandidate = event => {
    if (event.candidate === null) {
      var localSesion = btoa(JSON.stringify(pc.localDescription));

      var data = "user=mauro&id_camera=1&token="+localSesion;
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
  xhr.timeout = 60000;
  xhr.open("POST", host+"/sendtokenstreamer");
  xhr.setRequestHeader("cache-control", "no-cache");
  xhr.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
  xhr.setRequestHeader('token', 'eltoken');
  xhr.send(data);
}

window.startSession = (id, remoteSesion) => {
  try {
    pcMap.get(id).setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(remoteSesion))))
  } catch (e) {
    alert(e)
  }
}