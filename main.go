package main

import (
	"log"

	"goftp.io/server/v2"
	"goftp.io/server/v2/driver/file"
)

type auth struct {
}

func (*auth) CheckPasswd(ctx *server.Context, user string, passwd string) (bool, error) {
	log.Println(user, passwd)
	return true, nil
}

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	d, err := file.NewDriver("./")
	if err != nil {
		log.Fatal(err)
	}

	ftpServer, err := server.NewServer(&server.Options{
		// Driver: server.NewMultiDriver(driver.Drivers),
		Driver:       d,
		Name:         "Mafia FTP Server",
		Auth:         &auth{},
		Perm:         server.NewSimplePerm("wweir", "wweir"),
		Port:         3000,
		RateLimit:    1 << 20,
		PublicIP:     "139.196.34.166",
		PassivePorts: "50000-60000",
	})
	if err != nil {
		log.Fatal("Error creating server:", err)
	}

	err = ftpServer.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
