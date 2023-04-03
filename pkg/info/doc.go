// Package info starts an HTTP service to provide the gRPC server information.
//
// The server information can be used by the client to determine:
//   - what is right protocol to use (UDS or TCP)
//   - what is the numaflow sdk version used by the server
//   - what is language used by the server
//
// The gPRC server (UDF, UDSink, etc) must start the InfoServer with correct ServerInfo populated when it starts.
// The client is supposed to call the InfoServer to get the server information, before it starts to communicate with the gRPC server.
package info
