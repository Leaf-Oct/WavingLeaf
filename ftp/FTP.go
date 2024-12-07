package ftp

import (
	"log"

	"goftp.io/server/v2"
	"goftp.io/server/v2/driver/file"
)

func FTPTest() {
	driver, err := file.NewDriver("/home/leaf")
	drivers_map := make(map[string]server.Driver)
	drivers_map["leaf"] = driver
	if err != nil {
		log.Fatal(err)
	}

	s, err := server.NewServer(&server.Options{
		Driver:    NewMultiUserDriver(drivers_map),
		Auth:      &UserAuth{},
		Perm:      server.NewSimplePerm("root", "root"),
		RateLimit: 0,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
