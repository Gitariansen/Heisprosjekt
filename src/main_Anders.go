

/*func main() {
	/* // Test connection functions
	conn := make(chan network.Connection)
	msg_chan := make(chan network.UDPMessage)

	var test_msg network.UDPMessage
	test_msg.Msg = "I'm aliiiiiiveeee"

	go network.Init(conn)
	time.Sleep(100 * time.Millisecond)
	go network.Recive_msg_UDP(msg_chan)
	time.Sleep(100 * time.Millisecond)
	fmt.Println("stuck test 1")

	network.Broadcast_UDP(conn, test_msg, msg_chan)
	

	// Test Peers
	var testPeers []Peer
	peer1 := Test_struct{Address: "test1"}
	peer2 := Test_struct{Address: "test2"}
	peer3 := Test_struct{Address: "test3"}

	testPeers := append(testPeers, peer1)
	testPeers := append(testPeers, peer2)
	testPeers := append(testPeers, peer3)

}*/
