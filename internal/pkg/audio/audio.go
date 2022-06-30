package audio

import (
	"context"
	"errors"
	"github.com/namejlt/geektime-downloader/internal/pkg/logger"
	"github.com/namejlt/geektime-downloader/pconst"
	"os"
	"path/filepath"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/namejlt/geektime-downloader/internal/pkg/filenamify"
)

const (
	// MP3Extension ...
	MP3Extension = ".mp3"
)

// DownloadAudio ...
func DownloadAudio(ctx context.Context, downloadAudioURL, dir, title string) error {
	filenamifyTitle := filenamify.Filenamify(title)
	c := resty.New()
	c.SetOutputDirectory(dir).
		SetRetryCount(1).
		SetTimeout(time.Minute).
		SetHeader(pconst.UserAgentHeaderName, pconst.UserAgentHeaderValue).
		SetHeader(pconst.OriginHeaderName, pconst.GeekBang).
		SetLogger(logger.DiscardLogger{})

	_, err := c.R().
		SetContext(ctx).
		SetOutput(filenamifyTitle + MP3Extension).
		Get(downloadAudioURL)

	if errors.Is(err, context.Canceled) {
		fullName := filepath.Join(dir, filenamifyTitle+MP3Extension)
		_ = os.Remove(fullName)
	}

	return err
}
