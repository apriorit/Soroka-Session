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

	"github.com/Soroka-EDMS/svc/sessions/pkgs/config"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/db"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/endpoints"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/handlers"
	"github.com/Soroka-EDMS/svc/sessions/pkgs/service"
)

func main() {
	//Command parameters

	var (
		httpAddr   = flag.String("address", ":443", "")
		consulAddr = flag.String("consul.address", "localhost:8500", "Consul agent address")
		conn       = flag.String("consul.sessionsdb", "sessionsdb", "database connection string")
		certKey    = flag.String("consul.tls.pubkey", "tls/pubKey", "tls certificate")
		privateKey = flag.String("consul.tls.privkey", "tls/privKey", "tls private key")
		signingKey = flag.String("consul.service.signingKey", "service/signingKey", "Secret key to sign JWT")
	)

	//Parse CLI parameters
	flag.Parse()

	//Get global logger
	logger := config.GetLogger().Logger

	//Log CLI parameters
	logger.Log(
		"address", *httpAddr,
		"consul.address", *consulAddr,
		"consul.sessionsdb", *conn,
		"consul.tls.pubkey", *certKey,
		"consul.tls.privkey", *privateKey,
		"consul.service.secret", *signingKey,
	)

	//Obtain consul k/v storage
	consulStorage, err := GetConsulClient(*consulAddr)
	config.LogAndTerminateOnError(err, "create consul client")

	//Obtain tls pair (raw)
	certKeyData, err := ConsulGetKey(consulStorage, *certKey)
	config.LogAndTerminateOnError(err, "obtain certificate key")
	privateKeyData, err := ConsulGetKey(consulStorage, *privateKey)
	config.LogAndTerminateOnError(err, "obtain private key")

	signSecret, err := ConsulGetKey(consulStorage, *signingKey)
	config.LogAndTerminateOnError(err, "obtain secret")

	//Get kpair from raw data
	logger.Log("pub", string(certKeyData), "priv", string(privateKeyData))
	cert, err := tls.X509KeyPair(certKeyData, privateKeyData)
	config.LogAndTerminateOnError(err, "create cert from raw key pair data")

	//Create sessions database
	rawConnectionStr, err := ConsulGetKey(consulStorage, *conn)
	config.LogAndTerminateOnError(err, "obtain database connection string")

	dbs, err := db.Connection(logger, string(rawConnectionStr))

	//Build service layers
	var handler http.Handler
	{
		logger.Log("Loading", "Creating Session service...")
		svc := service.Build(logger, dbs, signSecret, certKeyData)
		endp := endpoints.MakeServerEndpoints(svc)
		handler = handlers.MakeHTTPHandler(endp, logger)
	}

	logger.Log("Loading", "Starting Session service...")
	var g run.Group
	{
		tlsConfig := &tls.Config{Certificates: []tls.Certificate{cert}}
		httpsListener, err := tls.Listen("tcp", *httpAddr, tlsConfig)
		config.LogAndTerminateOnError(err, "create https listener")

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
