package p2p

import (
	"fmt"
	"net/http"

	"github.com/JhyeonLee/BlockChain/blockchain"
	"github.com/JhyeonLee/BlockChain/utils"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func Upgrade(rw http.ResponseWriter, r *http.Request) {
	// Port :3000 will upgrade the request from :4000
	openPort := r.URL.Query().Get("openPort")
	ip := utils.Splitter(r.RemoteAddr, ":", 0)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return openPort != "" && ip != ""
	}
	fmt.Printf("%s wants an upgrade\n", openPort)
	conn, err := upgrader.Upgrade(rw, r, nil)
	utils.HandleErr(err)
	initPeer(conn, ip, openPort)
	// time.Sleep(20 * time.Second)
	/// conn.WriteMessage(websocket.TextMessage, []byte("Hello from Port 3000!"))
	// peer.inbox <- []byte("Hello from Port 3000!")
}

// being called exclusively by rest api
func AddPeer(address, port, openPort string, broadcast bool) {
	// Port :4000 is requesting an upgrade from the port :3000
	fmt.Printf("%s wants to connect to port %s\n", openPort, port)
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%s/ws?openPort=%s", address, port, openPort), nil)
	utils.HandleErr(err)
	p := initPeer(conn, address, port) // peer
	// time.Sleep(10 * time.Second)
	// conn.WriteMessage(websocket.TextMessage, []byte("Hello from Port 4000!"))
	// peer.inbox <- []byte("Hello from Port 4000!")
	if broadcast {
		BroadcastNewPeer(p)
		return
	}
	sendNewestblock(p) // request and response
}

func BroadcastNewBlock(b *blockchain.Block) {
	for _, p := range Peers.v {
		notifyNewBlock(b, p)
	}
}

func BroadcastNewTx(tx *blockchain.Tx) {
	for _, p := range Peers.v {
		notifyNewTx(tx, p)
	}
}

func BroadcastNewPeer(newPeer *peer) {
	for key, p := range Peers.v {
		if key != newPeer.key {
			payload := fmt.Sprintf("%s:%s", newPeer.key, p.port) // newPeer.key : address and port of new peer (to) // p.port : our open port that is port of old peer (from)
			notifyNewPeer(payload, p)
		}
	}
}
