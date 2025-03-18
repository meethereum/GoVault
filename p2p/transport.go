package p2p

//represents  the remote node over a TCP established connection
type Peer interface{
	
}

//it is anything that handles the communication between nodes in nw

type  Transport interface{
	ListenAndAccept() error //listen and accept udp.tcp,grpc, or literrally any thing
}