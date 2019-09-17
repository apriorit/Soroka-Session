package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"

	c "github.com/Soroka-EDMS/svc/sessions/pkgs/config"
	e "github.com/Soroka-EDMS/svc/sessions/pkgs/endpoints"
	h "github.com/Soroka-EDMS/svc/sessions/pkgs/handlers"
	s "github.com/Soroka-EDMS/svc/sessions/pkgs/service"
)

func main() {
	//Command parameters

	var (
		httpAddr   = flag.String("address", ":443", "")
		consulAddr = flag.String("consul.address", "localhost:8500", "Consul agent address")
		certKey    = flag.String("consul.tls.pubkey", "tls/pubKey", "tls certificate")
		privateKey = flag.String("consul.tls.privkey", "tls/privKey", "tls private key")
		secret     = flag.String("consul.service.secret", "service/secret", "Secret to sign JWT")
	)

	//Parse CLI parameters
	flag.Parse()

	//Get global logger
	logger := c.GetLogger().Logger

	//Log CLI parameters
	logger.Log(
		"address", *httpAddr,
		"consul.address", *consulAddr,
		"consul.tls.pubkey", *certKey,
		"consul.tls.privkey", *privateKey,
		"consul.service.secret", *secret,
	)

	//Obtain consul k/v storage
	consulStorage, err := GetConsulClient(*consulAddr)
	c.LogAndTerminateOnError(err, "create consul client")

	//Obtain tls pair (raw)
	certKeyData, err := ConsulGetKey(consulStorage, *certKey)
	c.LogAndTerminateOnError(err, "obtain certificate key")
	privateKeyData, err := ConsulGetKey(consulStorage, *privateKey)
	c.LogAndTerminateOnError(err, "obtain private key")
	//Save public key
	c.SetPublicKey(certKeyData)

	signSecret, err := ConsulGetKey(consulStorage, *secret)
	c.LogAndTerminateOnError(err, "obtain secret")

	//Get kpair from raw data
	logger.Log("pub", string(certKeyData), "priv", string(privateKeyData))
	cert, err := tls.X509KeyPair(certKeyData, privateKeyData)
	c.LogAndTerminateOnError(err, "create cert from raw key pair data")

	//Build service layers
	var handler http.Handler
	{
		logger.Log("Loading", "Creating Session service...")
		logger.Log("Sign secret", string(signSecret))
		svc := s.Build(logger, string(signSecret))
		endp := e.MakeServerEndpoints(svc)
		handler = h.MakeHTTPHandler(endp, logger)
	}

	logger.Log("Loading", "Starting Session service...")
	var g run.Group
	{
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
		httpsListener, err := tls.Listen("tcp", *httpAddr, tlsConfig)
		c.LogAndTerminateOnError(err, "create https listener")

		g.Add(func() error {
			logger.Log("transport", "https", "addr", *httpAddr)
			return http.Serve(httpsListener, handler)
		}, func(err error) {
			logger.Log("transport", "https", "err", err)
			httpsListener.Close()
		})
	}
	{
		var (
			cancelInterrupt = make(chan struct{})
			cancel          = make(chan os.Signal, 2)
		)

		defer close(cancel)

		g.Add(func() error {
			signal.Notify(cancel, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-cancel:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	logger.Log("exit", g.Run())
}
