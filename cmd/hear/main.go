package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/wzshiming/hear"
)

var (
	modelPath = "./models/ggml-base.bin"
)

func init() {
	flag.StringVar(&modelPath, "m", modelPath, "modelPath")

	flag.Parse()
}

func main() {

	h := hear.NewHear(hear.WithModelPath(modelPath))
	if err := h.Check(); err != nil {
		slog.Error("failed to check", "err", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)

	slog.Info("Start to hear from you, press Ctrl+C to finish.", "modelPath", modelPath)

	out, err := h.Run(ctx)
	if err != nil {
		slog.Error("failed to hear", "err", err)
		os.Exit(1)
	}

	fmt.Println(out)
}
