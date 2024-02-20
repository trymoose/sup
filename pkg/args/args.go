package args

import (
	"github.com/jessevdk/go-flags"
	"os"
	"time"
)

type Args struct {
	Data             flags.Filename `short:"d" long:"data" description:"directory files are stored" default:"/data" env:"DATA_DIR"`
	UploadAddress    string         `short:"u" long:"upload-address" description:"address of upload server" default:"0.0.0.0:9090" env:"UPLOAD_ADDRESS"`
	UIAddress        string         `short:"c" long:"ui-address" description:"address of web client" default:"0.0.0.0:8080" env:"UI_ADDRESS"`
	ShutdownTimeout  time.Duration  `short:"t" long:"shutdown-timeout" description:"how long to wait before force stopping the server" default:"5s" env:"SHUTDOWN_TIMEOUT"`
	UploadTimeout    time.Duration  `short:"o" long:"upload-timeout" description:"how long client have to upload data" default:"15s" env:"UPLOAD_TIMEOUT"`
	LimitUpload      int64          `short:"l" long:"limit-upload" description:"size in bytes to limit upload size, 0 to disable" default:"0" env:"LIMIT_UPLOAD"`
	DestinationsFile flags.Filename `short:"f" long:"destinations" description:"file containing destination configs" default:"/destinations.yaml" env:"DESTINATIONS_FILE"`
}

func NewArgs() (*Args, error) {
	var args Args
	_, err := flags.Parse(&args)
	if err == nil {
		_, err = os.Stat(string(args.Data))
		return &args, err
	}
	return &args, err
}
