package client

import (
	"github.com/numaproj/numaflow-go/pkg/function"
	"google.golang.org/grpc/resolver"
	"log"
	"runtime"
	"strconv"
	"strings"
)

const (
	exampleScheme      = "example"
	exampleServiceName = "lb.example.grpc.io"
	connAddr           = "0.0.0.0"
)

var addrsList []string

type exampleResolverBuilder struct{}

func (*exampleResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &exampleResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			exampleServiceName: addrsList,
		},
	}
	r.start()
	return r, nil
}
func (*exampleResolverBuilder) Scheme() string { return exampleScheme }

type exampleResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *exampleResolver) start() {
	addrStrs := r.addrsStore[r.target.Endpoint()]
	addrs := make([]resolver.Address, len(addrStrs))
	//attrs := make(*attributes.Attributes, len(addrStrs))
	for i, s := range addrStrs {
		adr := strings.Split(s, ",")
		addr := adr[0]
		serv := adr[1]
		addrs[i] = resolver.Address{Addr: addr, ServerName: serv}
	}
	err := r.cc.UpdateState(resolver.State{Addresses: addrs})
	if err != nil {
		return
	}
}
func (*exampleResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*exampleResolver) Close()                                  {}
func (*exampleResolver) Resolve(target resolver.Target) {
	log.Println(target)
}

func buildConnAddrs(numCpu int) {
	var conn = make([]string, numCpu)
	for i := 0; i < numCpu; i++ {
		conn[i] = connAddr + function.Addr + "," + strconv.Itoa(i+1)
	}
	addrsList = conn
}
func init() {
	if function.MAP_MULTIPROC_SERV == true {
		resolver.Register(&exampleResolverBuilder{})
		numCpu := runtime.NumCPU()
		buildConnAddrs(numCpu)
		log.Println("TCP client list:", addrsList)
	}
}
