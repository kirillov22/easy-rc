package main

import (
	"flag"
	"fmt"
	"mouse-server/websocket"
	"net"
	//"github.com/go-vgo/robotgo"
	"log"
	"net/http"
	//"strconv"
	//"strings"
)

//func home(w http.ResponseWriter, r *http.Request) {
//	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
//}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", websocket.Server)
	//http.HandleFunc("/", home)
	var serverAddress, port = generateAddress()
	var outboundIp = getOutboundIP()

	log.Printf("Starting websocket server at: %s. Outbound address to connect to: %s:%d\n", serverAddress, outboundIp, port)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

func generateAddress() (string, int) {
	var port, err = getFreePort()
	if err != nil {
		log.Fatal("Failed to get a free port", err)
	}
	return fmt.Sprintf("0.0.0.0:%d", port), port
}

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

//var homeTemplate = template.Must(template.New("").Parse(`
//<!DOCTYPE html>
//<html>
//<head>
//<meta charset="utf-8">
//<script>
//window.addEventListener("load", function(evt) {
//
//    var output = document.getElementById("output");
//    var input = document.getElementById("input");
//    var ws;
//
//    var print = function(message) {
//        var d = document.createElement("div");
//        d.textContent = message;
//        output.appendChild(d);
//        output.scroll(0, output.scrollHeight);
//    };
//
//    document.getElementById("open").onclick = function(evt) {
//        if (ws) {
//            return false;
//        }
//        ws = new WebSocket("{{.}}");
//        ws.onopen = function(evt) {
//            print("OPEN");
//        }
//        ws.onclose = function(evt) {
//            print("CLOSE");
//            ws = null;
//        }
//        ws.onmessage = function(evt) {
//            print("RESPONSE: " + evt.data);
//        }
//        ws.onerror = function(evt) {
//            print("ERROR: " + evt.data);
//        }
//        return false;
//    };
//
//    document.getElementById("send").onclick = function(evt) {
//        if (!ws) {
//            return false;
//        }
//        print("SEND: " + input.value);
//        ws.send(input.value);
//        return false;
//    };
//
//    document.getElementById("close").onclick = function(evt) {
//        if (!ws) {
//            return false;
//        }
//        ws.close();
//        return false;
//    };
//
//});
//</script>
//</head>
//<body>
//<table>
//<tr><td valign="top" width="50%">
//<p>Click "Open" to create a connection to the server,
//"Send" to send a message to the server and "Close" to close the connection.
//You can change the message and send multiple times.
//<p>
//<form>
//<button id="open">Open</button>
//<button id="close">Close</button>
//<p><input id="input" type="text" value="Hello world!">
//<button id="send">Send</button>
//</form>
//</td><td valign="top" width="50%">
//<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
//</td></tr></table>
//</body>
//</html>
//`))
