package native

import (
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/xuperchain/contract-sdk-go/code"
	pbrpc "github.com/xuperchain/contract-sdk-go/pbrpc"
	"google.golang.org/grpc"
)

const (
	xchainPingTimeout = "XCHAIN_PING_TIMEOUT"
	xchainCodePort    = "XCHAIN_CODE_PORT"
	xchainChainAddr   = "XCHAIN_CHAIN_ADDR"
	// XCHAIN_CODE_ADDR is standard networking address
	// see documentation for net.Listen in standard linrary for more detials
	xchainCodeAddr = "XCHAIN_CODE_ADDR"
)

type driver struct {
}

// New returns a native driver
func New() code.Driver {
	return new(driver)
}

func (d *driver) Serve(contract code.Contract) {
	chainAddr := os.Getenv(xchainChainAddr)
	codePort := os.Getenv(xchainCodePort)
	codeAddr := os.Getenv(xchainCodeAddr)

	if chainAddr == "" {
		panic("empty XCHAIN_CHAIN_ADDR env")
	}
	if codeAddr == "" && codePort == "" {
		panic("empty XCHAIN_CODE_PORT env")
	}

	listenAddress := "127.0.0.1:" + codePort
	if codeAddr != "" {
		uri, err := url.Parse(codeAddr)
		if err != nil {
			panic(err)
		}
		listenAddress = uri.Host + uri.Path
	}

	nativeCodeService := newNativeCodeService(chainAddr, contract)
	rpcServer := grpc.NewServer()
	pbrpc.RegisterNativeCodeServer(rpcServer, nativeCodeService)

	var listener net.Listener
	listener, err := net.Listen("tcp", listenAddress)
	if err != nil {
		panic(err)
	}

	go rpcServer.Serve(listener)

	sigch := make(chan os.Signal, 2)
	signal.Notify(sigch, os.Interrupt, syscall.SIGTERM, syscall.SIGPIPE)

	timer := time.NewTicker(1 * time.Second)
	running := true
	pingTimeout := getPingTimeout()
	for running {
		select {
		case sig := <-sigch:
			running = false
			log.Print("receive signal ", sig)
		case <-timer.C:
			lastping := nativeCodeService.LastpingTime()
			if time.Since(lastping) > pingTimeout {
				log.Print("ping timeout")
				running = false
			}
		}
	}
	rpcServer.GracefulStop()
	nativeCodeService.Close()
	log.Print("native code ended")
}

func getPingTimeout() time.Duration {
	envtimeout := os.Getenv(xchainPingTimeout)
	if envtimeout == "" {
		return 3 * time.Second
	}
	timeout, err := strconv.Atoi(envtimeout)
	if err != nil {
		return 3 * time.Second
	}
	return time.Duration(timeout) * time.Second
}
