package hear

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

type Hear struct {
	ModelPath      string
	FfmpegPath     string
	WhisperCppPath string
}

func WithModelPath(modelPath string) func(*Hear) {
	return func(h *Hear) {
		h.ModelPath = modelPath
	}
}

func WithFfmpegPath(ffmpegPath string) func(*Hear) {
	return func(h *Hear) {
		h.FfmpegPath = ffmpegPath
	}
}

func WithWhisperCppPath(whisperCppPath string) func(*Hear) {
	return func(h *Hear) {
		h.WhisperCppPath = whisperCppPath
	}
}

func NewHear(opts ...func(*Hear)) *Hear {
	h := &Hear{
		ModelPath:      "./models/ggml-base.bin",
		FfmpegPath:     "ffmpeg",
		WhisperCppPath: "whisper-cpp",
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}

func (h *Hear) Check() error {
	ffmpegPath, err := exec.LookPath(h.FfmpegPath)
	if err != nil {
		return err
	}

	whisperCppPath, err := exec.LookPath(h.WhisperCppPath)
	if err != nil {
		return err
	}

	h.FfmpegPath = ffmpegPath
	h.WhisperCppPath = whisperCppPath
	return nil
}

func (h *Hear) Run(ctx context.Context) (string, error) {
	buf := bytes.NewBuffer(nil)
	errBuf := bytes.NewBuffer(nil)
	ffmpegCmd := exec.CommandContext(ctx, h.FfmpegPath, "-f", "avfoundation", "-i", ":0", "-ar", "16000", "-f", "wav", "-")
	ffmpegCmd.Stdout = buf
	ffmpegCmd.Stderr = errBuf
	if err := ffmpegCmd.Run(); err != nil {
		if buf.Len() == 0 {
			return "", fmt.Errorf("failed to ffmpeg: %w: %s", err, errBuf.String())
		}
	}
	if buf.Len() == 0 {
		return "", fmt.Errorf("nothing to hear: %s", errBuf.String())
	}

	out := bytes.NewBuffer(nil)

	errBuf.Reset()

	whisperCppCmd := exec.CommandContext(context.Background(), h.WhisperCppPath, "-m", h.ModelPath, "-nt", "-f", "-", "-of", "-")
	whisperCppCmd.Stdin = buf
	whisperCppCmd.Stdout = out
	whisperCppCmd.Stderr = errBuf
	if err := whisperCppCmd.Run(); err != nil {
		return "", fmt.Errorf("failed to whisper-cpp: %w: %s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}
