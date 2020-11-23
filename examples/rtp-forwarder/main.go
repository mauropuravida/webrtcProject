package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/pion/rtcp"
	"github.com/pion/webrtc/v3"
	"github.com/pion/webrtc/v3/examples/internal/signal"
	"io/ioutil"
	"flag"
	"strconv"
	"net/http"
	"net/url"
	"strings"
	"os"
)

type udpConn struct {
	conn *net.UDPConn
	port int
}

var(
	host = ""
	)

//Only work in linux
func createSdp(addr string, videoPort string){
	data := []byte("v=0\no=- 0 0 IN IP4 "+addr+"\ns=WebRTC "+addr+":"+videoPort+"\nc=IN IP4 "+addr+"\nt=0 0\nm=video "+videoPort+" RTP/AVP 96\na=rtpmap:96 VP8/90000")
    prefix := addr+"_"+videoPort+"-*.sdp"
    tmpFile, err := ioutil.TempFile(os.TempDir(), prefix)
    if err != nil {
        fmt.Println("Cannot create temporary file", err)
    }
    tmpFile.Write(data)
}

func main() {
	addr := flag.String("address", "127.0.0.200", "Address to host the HTTP server on.")
	portt := flag.Int("port", 4000, "Address to host the HTTP server on.")
	hostt := flag.String("host", "http://localhost", "")
	flag.Parse()

	host = *hostt

	fmt.Printf("Listening address %s \n", *addr)

	if *addr == ""{
		fmt.Printf("Empty address %s \n", *addr)
		return
	}

	createSdp(*addr, strconv.Itoa(*portt))
	// Create context
	ctx, cancel := context.WithCancel(context.Background())

	// Create a MediaEngine object to configure the supported codec
	m := webrtc.MediaEngine{}

	// Setup the codecs you want to use.
	// We'll use a VP8 codec but you can also define your own
	m.RegisterCodec(webrtc.NewRTPOpusCodec(webrtc.DefaultPayloadTypeOpus, 48000))
	m.RegisterCodec(webrtc.NewRTPVP8Codec(webrtc.DefaultPayloadTypeVP8, 90000))

	// Create the API object with the MediaEngine
	api := webrtc.NewAPI(webrtc.WithMediaEngine(m))

	// Everything below is the Pion WebRTC API! Thanks for using it ❤️.

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	peerConnection, err := api.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// Allow us to receive 1 audio track, and 1 video track
	if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeAudio); err != nil {
		panic(err)
	} else if _, err = peerConnection.AddTransceiverFromKind(webrtc.RTPCodecTypeVideo); err != nil {
		panic(err)
	}

	// Create a local addr
	var laddr *net.UDPAddr
	if laddr, err = net.ResolveUDPAddr("udp", *addr+":"); err != nil {
		panic(err)
	}

	// Prepare udp conns
	udpConns := map[string]*udpConn{
		"video": {port: *portt},
	}
	for _, c := range udpConns {
		// Create remote addr
		var raddr *net.UDPAddr
		if raddr, err = net.ResolveUDPAddr("udp", fmt.Sprintf(*addr+":%d", c.port)); err != nil {
			panic(err)
		}

		// Dial udp
		if c.conn, err = net.DialUDP("udp", laddr, raddr); err != nil {
			panic(err)
		}
		defer func(conn net.PacketConn) {
			if closeErr := conn.Close(); closeErr != nil {
				panic(closeErr)
			}
		}(c.conn)
	}

	// Set a handler for when a new remote track starts, this handler will forward data to
	// our UDP listeners.
	// In your application this is where you would handle/process audio/video
	peerConnection.OnTrack(func(track *webrtc.Track, receiver *webrtc.RTPReceiver) {
		// Retrieve udp connection
		c, ok := udpConns[track.Kind().String()]
		if !ok {
			return
		}

		// Send a PLI on an interval so that the publisher is pushing a keyframe every rtcpPLIInterval
		go func() {
			ticker := time.NewTicker(time.Second * 2)
			for range ticker.C {
				if rtcpErr := peerConnection.WriteRTCP([]rtcp.Packet{&rtcp.PictureLossIndication{MediaSSRC: track.SSRC()}}); rtcpErr != nil {
					fmt.Println(rtcpErr)
				}
			}
		}()

		b := make([]byte, 1500)
		for {
			// Read
			n, readErr := track.Read(b)
			if readErr != nil {
				panic(readErr)
			}

			// Write
			if _, err = c.conn.Write(b[:n]); err != nil {
				// For this particular example, third party applications usually timeout after a short
				// amount of time during which the user doesn't have enough time to provide the answer
				// to the browser.
				// That's why, for this particular example, the user first needs to provide the answer
				// to the browser then open the third party application. Therefore we must not kill
				// the forward on "connection refused" errors
				if opError, ok := err.(*net.OpError); ok && opError.Err.Error() == "write: connection refused" {
					continue
				}
				panic(err)
			}
		}
	})

	// Set the handler for ICE connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		fmt.Printf("Connection State has changed %s \n", connectionState.String())

		if connectionState == webrtc.ICEConnectionStateConnected {
			fmt.Println("Ctrl+C the remote client to stop the demo")
		} else if connectionState == webrtc.ICEConnectionStateFailed ||
			connectionState == webrtc.ICEConnectionStateDisconnected {
			fmt.Println("Done forwarding")
			cancel()
		}
	})

	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	signal.Decode(signal.MustReadStdin(), &offer)

	// Set the remote SessionDescription
	if err = peerConnection.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	// Create answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(peerConnection)

	// Sets the LocalDescription, and starts our UDP listeners
	if err = peerConnection.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	// Output the answer in base64 so we can paste it in browser
	fmt.Println(signal.Encode(*peerConnection.LocalDescription()))


    endpoint := host+"/sendtokenconnect"
    data := url.Values{}
    data.Set("token", signal.Encode(*peerConnection.LocalDescription()))

    client := &http.Client{}
    r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload

    r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

    if err != nil {
        panic(err)
    }

    res, err := client.Do(r)
    if err != nil {
        panic(err)
    }
    //fmt.Println(res.Status)
    defer res.Body.Close()
    //body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        panic(err)
    }
    //fmt.Println(string(body))

	// Wait for context to be done
	<-ctx.Done()
}
