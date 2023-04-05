package client

import (
	"github.com/numaproj/numaflow-go/pkg/function"
	"google.golang.org/grpc/resolver"
	"log"
	"strconv"
	"strings"
)

const (
	custScheme      = "example"
	custServiceName = "lb.example.grpc.io"
	connAddr        = "0.0.0.0"
)

var addrsList []string

type multiProcResolverBuilder struct{}

func (*multiProcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &multiProcResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			custServiceName: addrsList,
		},
	}
	r.start()
	return r, nil
}
func (*multiProcResolverBuilder) Scheme() string { return custScheme }

type multiProcResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func (r *multiProcResolver) start() {
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
func (*multiProcResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*multiProcResolver) Close()                                  {}
func (*multiProcResolver) Resolve(target resolver.Target) {
	log.Println(target)
}

func buildConnAddrs(numCpu int) {
	var conn = make([]string, numCpu)
	for i := 0; i < numCpu; i++ {
		conn[i] = connAddr + function.Addr + "," + strconv.Itoa(i+1)
	}
	addrsList = conn
}

//func init() {
//	if function.MAP_MULTIPROC_SERV == true {
//		resolver.Register(&multiProcResolverBuilder{})
//		numCpu := runtime.NumCPU()
//		buildConnAddrs(numCpu)
//		log.Println("TCP client list:", addrsList)
//	}
//}
