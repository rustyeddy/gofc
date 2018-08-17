package main

import (
	"flag"
	"fmt"

	"github.com/rustyeddy/gofc"
)

type SampleController struct {
	*gofc.OFController
}

var (
	servport string
)

func init() {
	flag.StringVar(&servport, "server port", ":6633", "Default :6633 to Start Open Flowin")
}

func NewSampleController() *SampleController {
	return &SampleController{gofc.NewOFController()}
}

func main() {
	// regist app
	ofc := NewSampleController()

	fmt.Printf("ofc: %+v\n", ofc)
	gofc.GetAppManager().RegistApplication(ofc)

	// start server
	gofc.ServerLoop(servport)
}
