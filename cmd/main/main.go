package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"strings"

	mcrlockup "github.com/ntnn/mcr-lockup"
)

var (
	fKubeconfigs = flag.String("kubeconfigs", "", "Path to kubeconfig files, comma separated")
)

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	if err := mcrlockup.Run(ctx, strings.Split(*fKubeconfigs, ",")); err != nil {
		panic(err)
	}
}
