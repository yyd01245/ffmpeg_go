package main

import (
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/giorgisio/goav/avcodec"
	"github.com/giorgisio/goav/avdevice"
	"github.com/giorgisio/goav/avfilter"
	"github.com/giorgisio/goav/avformat"
	"github.com/giorgisio/goav/avutil"
	"github.com/giorgisio/goav/swresample"
	"github.com/giorgisio/goav/swscale"
	"github.com/jessevdk/go-flags"
)

var opts struct {
	InFile  string `short:"i" long:"infile" description:"Input file path"`
	OutFile string `short:"o" long:"outfile" default:"./output.264" description:"output file path"`
	Filter  string `short:"f" long:"filter" description:"filter data"`
	Log     string `short:"l" long:"log" description:"the log file to tail -f"`
}

var log *logrus.Logger

func init() {
	log = logrus.New()
	log.Level = logrus.InfoLevel
	f := new(logrus.TextFormatter)
	f.TimestampFormat = "2006-01-02 15:04:05"
	f.FullTimestamp = true
	log.Formatter = f
}
func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		if !strings.Contains(err.Error(), "Usage") {
			log.Fatalf("error: %v", err)
		} else {
			return
		}
	}
	out, err := os.Create(opts.Log)
	if err != nil {
		log.Errorf("creat %s file error", opts.Log)
		return
	}
	defer out.Close()
	log.Out = out
	filename := opts.OutFile

	// Register all formats and codecs
	avformat.AvRegisterAll()

	avcodec.AvcodecRegisterAll()

	log.Printf("AvFilter Version:\t%v", avfilter.AvfilterVersion())
	log.Printf("AvDevice Version:\t%v", avdevice.AvdeviceVersion())
	log.Printf("SWScale Version:\t%v", swscale.SwscaleVersion())
	log.Printf("AvUtil Version:\t%v", avutil.AvutilVersion())
	log.Printf("AvCodec Version:\t%v", avcodec.AvcodecVersion())
	log.Printf("Resample Version:\t%v", swresample.SwresampleLicense())

	log.Println("register all success")
	// Open video file
	var (
		ctxtFormat    *avformat.Context
		packet        *avcodec.Packet
		videoCodec    *avcodec.Codec
		videoFrame    *avutil.Frame
		url           string
		videoStream   int
		audioStream   int
		frameFinished int
		numBytes      int
		frameSize     int
	)
	if avformat.AvformatOpenInput(&ctxtFormat, filename, nil, nil) != 0 {
		log.Println("Error: Couldn't open file.")
		return
	}

	// Retrieve stream information
	if ctxtFormat.AvformatFindStreamInfo(nil) < 0 {
		log.Println("Error: Couldn't find stream information.")
		return
	}

	ctxtFormat.AvDumpFormat(0, url, 0)
	videoStream = -1
	n := ctxtFormat.NbStreams()
	s := ctxtFormat.Streams()
	log.Print("Number of Streams:", n)
	for i := 0; i < int(n); i++ {
		log.Println("stream Number : ", i)
		if (*avformat.CodecContext)(s.Codec()) != nil {
			videoStream = i
			break
		}
	}
	if videoStream == -1 {
		log.Println("couldn't find video stream")
		return
	}
}
