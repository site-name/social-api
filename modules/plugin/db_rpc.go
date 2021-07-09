package plugin

import (
	"net/rpc"
)

// dbRPCClient contains the client-side logic to handle the RPC communication
// with the server. It's API is hand-written because we do not expect
// new methods to be added very frequently.
type dbRPCClient struct {
	client *rpc.Client
}

// dbRPCServer is the server-side component which is responsible for calling
// the driver methods and properly encoding the responses back to the RPC client.
type dbRPCServer struct {
	dbImpl Driver
}
