
/* eslint-env browser */
var host = "http://localhost"
//var host = "https://testingwebrtc.ml"

let pc = new RTCPeerConnection({
  iceServers: [
  {
    urls: 'stun:stun.l.google.com:19302'
  }
  ]
})
var log = msg => {
  document.getElementById('logs').innerHTML += msg + '<br>'
}

navigator.mediaDevices.getUserMedia({ video: true, audio: true })
.then(stream => {
  pc.addStream(document.getElementById('video1').srcObject = stream)
  pc.createOffer().then(d => pc.setLocalDescription(d)).catch(log)
}).catch(log)


var  localSesion = 'Test'
pc.oniceconnectionstatechange = e => log(pc.iceConnectionState)
pc.onicecandidate = event => {
  if (event.candidate === null) {
    localSesion = btoa(JSON.stringify(pc.localDescription))
    document.getElementById('localSessionDescription').value = localSesion

    getToken("user=mauro&id_camera=1&token="+localSesion);
  }
}

function getToken(data){
  var xhr = new XMLHttpRequest();

  xhr.onload = function () {
  console.log(this.readyState);
    if (this.readyState === 4) {
      document.getElementById('remoteSessionDescription').value = this.responseText
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

  xhr.send(data);
}

window.startSession = () => {
  let sd = document.getElementById('remoteSessionDescription').value
  if (sd === '') {
    return alert('Session Description must not be empty NOW')
  }

  try {
    pc.setRemoteDescription(new RTCSessionDescription(JSON.parse(atob(sd))))
  } catch (e) {
    alert(e)
  }
}
