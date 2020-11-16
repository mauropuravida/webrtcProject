
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
function connectStream(){

  myImg = new Image()
  myImg.src = document.getElementById('url').value
  myImg.crossOrigin = "Anonymous"

  var canvas = document.getElementById("myCanvas");
  var context = canvas.getContext("2d");
  

  (function loop() {
    if (!this.paused && !this.ended) {
      context.drawImage(myImg, 0, 0, canvas.width, canvas.height);
      setTimeout(loop, 1000 / 30); // drawing at 30fps
    }
  })();

  var stream = canvas.captureStream(30)
  pc.addStream(stream)
  pc.createOffer().then(d => pc.setLocalDescription(d)).catch(log)
  
}

var  localSesion = 'Test'
pc.oniceconnectionstatechange = e => log(pc.iceConnectionState)
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