package main

import (
	Http "net/http"
	OS "os"
	Strings "strings"

	Bcrypt "golang.org/x/crypto/bcrypt"

	Echo "github.com/labstack/echo/v4"
	Middleware "github.com/labstack/echo/v4/middleware"
	Config "github.com/simsibimsiwimsi/hetzner-dyndns/config"
	Hetzner "github.com/simsibimsiwimsi/hetzner-dyndns/hetzner"
)

func main() {

	config := Config.Initialize()
	e := Echo.New()

	e.Use(Middleware.Logger())
	e.Use(Middleware.Recover())

	e.Use(Middleware.BasicAuth(func(username, password string, c Echo.Context) (bool, error) {
		hostname := c.QueryParam("hostname")
		dnsRecordName := Strings.Split(hostname, ".")[0]
		userAndPassword := config.Users[dnsRecordName]
		if userAndPassword == nil {
			return false, nil
		}

		if username == userAndPassword.User && Bcrypt.CompareHashAndPassword([]byte(userAndPassword.Password), []byte(password)) == nil {
			return true, nil
		}
		createPasswordHash(userAndPassword.Password)
		return false, nil
	}))

	// Route level middleware
	track := func(next Echo.HandlerFunc) Echo.HandlerFunc {
		return func(c Echo.Context) error {
			println("request to /")
			return next(c)
		}
	}

	e.GET("/", func(c Echo.Context) error {

		hostname := c.QueryParam("hostname")
		dnsRecordName := Strings.Split(hostname, ".")[0]
		ipv4 := c.QueryParam("ipv4")
		ipv6 := c.QueryParam("ipv6")

		dnsZone := Hetzner.NewDnsZone(config.Hetzner.Dns["zone-id"], config.Hetzner.Dns["auth-api-token"])

		ipv4record, ipv6record, err := dnsZone.CreateOrUpdateIpV4andV6Records(dnsRecordName, ipv4, ipv6)
		if err != nil {
			return c.JSON(Http.StatusInternalServerError, struct{ Error string }{Error: err.Error()})
		}

		return c.JSON(Http.StatusOK, struct{ Message string }{"Upserted dyndns records on hetzner cloud for " + dnsRecordName + " pointing to " + ipv4record.Value + " / " + ipv6record.Value})
	}, track)

	e.GET("/health", func(c Echo.Context) error {
		return c.JSON(Http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	httpPort := OS.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8053"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

func createPasswordHash(password string) {
	passwordBytes, _ := Bcrypt.GenerateFromPassword([]byte(password), Bcrypt.MinCost)
	println(string(passwordBytes))
}
