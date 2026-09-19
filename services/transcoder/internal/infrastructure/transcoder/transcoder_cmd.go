package transcoder

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

type TranscoderCmd struct {
	bin string
}

func NewTranscoderCmd(bin string) *TranscoderCmd {
	return &TranscoderCmd{
		bin: bin,
	}
}

func (transcoder *TranscoderCmd) Transcode(ctx context.Context, inputFilePath, outputFilePath string) error {
	cmd := exec.CommandContext(ctx, transcoder.bin,
		"-nostdin", "-y",
		"-i", inputFilePath,
		"-vf", "scale=-2:720",
		"-c:v", "libx264", "-preset", "veryfast", "-crf", "23",
		"-c:a", "aac", "-b:a", "128k",
		"-movflags", "+faststart",
		outputFilePath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg: %w: %s", err, stderr.String())
	}
	return nil
}
