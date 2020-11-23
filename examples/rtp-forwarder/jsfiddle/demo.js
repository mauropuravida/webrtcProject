
/* eslint-env browser */
var host = "http://localhost"
//var host = "https://testingwebrtc.ddns.net"

let pc = new RTCPeerConnection({
  iceServers: [
  {
    urls: 'stun:stun.l.google.com:19302'
  }
  ]
})

/*navigator.mediaDevices.getUserMedia({ video: true, audio: true })
.then(stream => {
  pc.addStream(document.getElementById('video1').srcObject = stream)
  pc.createOffer().then(d => pc.setLocalDescription(d)).catch(log)
}).catch(log)*/
//var sobject = document.getElementById('streamtest').srcObject
//var stream = new MediaStream()

let rotateAngle = 0;
let xr;
let yr;
function rotateStream(){
  rotateAngle = rotateAngle + 90;

  switch(rotateAngle) {
  case 90:
    yr = yr * -1
    break;
  case 180:
    xr = xr * -1
    break;
  case 270:
    yr = yr * -1
    break;
  case 360:
    xr = xr * -1
    rotateAngle = 0
    break;
  default:
    // code block
  }
}

var transform = false;
function changeColorStream(){
  if(transform)
    transform = false
  else
    transform = true
}

function handlerImage(element, myImg){
    var context = element.getContext("2d");

    context.setTransform(1, 0, 0, 1, 0, 0);
    context.clearRect(0,0,element.width, element.height);
    context.setTransform(1, 0, 0, 1, element.xr, element.yr);      
    context.rotate(parseInt(rotateAngle) * Math.PI / 180);
    context.drawImage(myImg, 0, 0, element.width, element.height);

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


async function connectStream(){
  var myImg = new Image()
  myImg.src = document.getElementById('url').value
  myImg.crossOrigin = "Anonymous"

  var canvasClient = document.getElementById("myCanvas");

  var canvasStream = document.createElement("CANVAS");
  canvasStream.width = 480
  canvasStream.height = 360

  xr = canvasClient.width;
  yr = canvasClient.height;

  (function loop() {
    if (!this.paused && !this.ended) {
        //draw in canvas client
        handlerImage(canvasClient, myImg);
        //draw in canvas for stream best resolution
        handlerImage(canvasStream, myImg);
      }

      setTimeout(loop, 1000 / 60); // drawing at 60fps
  })();

  var stream = canvasStream.captureStream(60)
  pc.addStream(stream)
  pc.createOffer().then(d => pc.setLocalDescription(d)).catch("error")
  
}

var  localSesion = 'Test'
pc.oniceconnectionstatechange = e => console.log(pc.iceConnectionState)
pc.onicecandidate = event => {
  if (event.candidate === null) {
    localSesion = btoa(JSON.stringify(pc.localDescription));

    var data = "user=mauro&id_camera=1&token="+localSesion;
    console.log(data);
    getToken(data);
  }
}

var remoteSesion = ''
function getToken(data){
  var xhr = new XMLHttpRequest();
  xhr.withCredentials = true;
  xhr.onload = function () {
  console.log(this.readyState);
    if (this.readyState === 4) {
      remoteSesion = this.responseText
      setTimeout(window.startSession(), 3000);
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

window.startSession = () => {
  try {
    pc.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(remoteSesion))))
  } catch (e) {
    alert(e)
  }
}