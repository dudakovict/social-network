package main

import (
	"errors"
	"expvar"
	"fmt"
	"net"
	"net/smtp"
	"os"

	"github.com/ardanlabs/conf"
	es "github.com/dudakovict/social-network/business/core/email"
	"github.com/dudakovict/social-network/business/data/email"
	"github.com/dudakovict/social-network/foundation/logger"
	_ "go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var build = "develop"

func main() {

	// Construct the application logger.
	log, err := logger.New("email-api")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	// Perform the startup and shutdown sequence.
	if err = run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {

	// =========================================================================
	// GOMAXPROCS

	// Want to see what maxprocs reports.

	//opt := maxprocs.Logger(log.Infof)

	// Set the correct number of threads for the service
	// based on what is available either by the machine or quotas.
	/*
		if _, err := maxprocs.Set(opt); err != nil {
			return fmt.Errorf("maxprocs: %w", err)
		}
		log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))
	*/

	// =========================================================================
	// Configuration
	cfg := struct {
		conf.Version
		Auth struct {
			KeysFolder string `conf:"default:zarf/keys/"`
			ActiveKID  string `conf:"default:54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"`
		}
		GRPC struct {
			Network string `conf:"default:tcp"`
			Address string `conf:"default:0.0.0.0:50084"`
		}
		SMTP struct {
			Username string `conf:"default:example@gmail.com,env:SMTP_USERNAME"`
			Password string `conf:"default:jusesendsmtpauth,env:SMTP_PASSWORD"`
			Host     string `conf:"default:smtp.gmail.com"`
			Port     string `conf:"default:587"`
			Address  string `conf:"default:smtp.gmail.com:587"`
		}
	}{
		Version: conf.Version{
			SVN:  build,
			Desc: "copyright information here",
		},
	}

	const prefix = "EMAIL"
	help, err := conf.ParseOSArgs(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

	// =========================================================================
	// App Starting

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	expvar.NewString("build").Set(build)

	// =========================================================================
	// Start gRPC Service

	log.Infow("startup", "status", "gRPC service started", "host", cfg.GRPC.Address)

	conn, err := net.Listen(cfg.GRPC.Network, cfg.GRPC.Address)

	if err != nil {
		return fmt.Errorf("announcing on the local network: %w", err)
	}

	grpc := grpc.NewServer()

	/*
		cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
		cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")
	*/

	auth := smtp.PlainAuth("", cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Host)
	es := es.NewEmailServer(log, cfg.SMTP.Address, cfg.SMTP.Username, auth)

	email.RegisterEmailServer(grpc, &es)

	if err := grpc.Serve(conn); err != nil {
		log.Errorw("shutdown", "status", "gRPC v1 server closed", "host", cfg.GRPC.Address, "ERROR", err)
	}

	return nil
}
