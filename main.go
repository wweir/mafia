package main

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/wweir/mafia/drivers"
	"github.com/wweir/mafia/drivers/sftp"
	"goftp.io/server/v2"
)

type auth struct {
}

func (*auth) CheckPasswd(ctx *server.Context, user string, passwd string) (bool, error) {
	log.Info().
		Str("user", user).
		Str("password", passwd).
		Msg("auth info")
	return true, nil
}

type logger struct {
	server.DiscardLogger
	zerolog.Logger
}

func (l *logger) Print(sessionID string, message interface{}) {
	l.Info().Str("sess", sessionID).Interface("msg", message).Send()
}
func (l *logger) Printf(sessionID string, format string, v ...interface{}) {
	l.Info().Str("sess", sessionID).Msgf(format, v...)
}

// func (l *logger) PrintCommand(sessionID string, command string, params string) {
// 	l.Info().Str("sess", sessionID).Str("cmd", command).Str("params", params).Send()
// }
// func (l *logger) PrintResponse(sessionID string, code int, message string) {
// 	l.Info().Str("sess", sessionID).Int("code", code).Msg(message)
// }

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out: os.Stderr,
		FormatCaller: func(i interface{}) string {
			caller := i.(string)
			if idx := strings.Index(caller, "/pkg/mod/"); idx > 0 {
				return caller[idx+9:]
			}
			if idx := strings.LastIndexByte(caller, '/'); idx > 0 {
				return caller[idx+1:]
			}
			return caller
		},
	}).With().Timestamp().Stack().Caller().Logger()
}

func main() {
	// ftp, err := tar.NewFTP("drivers/tar/a.tar")
	dav, err := sftp.NewWebdav(&sftp.SSHConfig{
		Addr: "127.0.0.1",
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	ftpServer, err := server.NewServer(&server.Options{
		Driver:    drivers.NewFTPDriver(nil, dav),
		Name:      "Mafia FTP Server",
		Auth:      &auth{},
		Perm:      server.NewSimplePerm("wweir", "wweir"),
		Port:      3000,
		RateLimit: 1 << 20,
		PublicIP:  "139.196.34.166",
		Logger: &logger{
			Logger: log.Logger.With().CallerWithSkipFrameCount(3).Logger(),
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("creating server")
	}

	err = ftpServer.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msg("starting server")
	}
}
