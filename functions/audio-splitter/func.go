package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/denismakogon/s3-pollster/api"
	"github.com/denismakogon/s3-pollster/common"
	"github.com/fnproject/fdk-go"
	"github.com/google/uuid"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(withError))
}

var (
	Ext          = "wav"
	FFMPEG       = "ffmpeg"
	ArgVerbose   = "-loglevel panic"
	Args         = "-y -i %v -vn -f wav -"
	DebugMode, _ = strconv.Atoi(common.WithDefault("DEBUG_MODE", "1"))
	s3URL        = common.WithDefault("S3_URL",
		"s3://admin:password@docker.for.mac.localhost:9000/us-east-1/default-bucket")
)

func withError(ctx context.Context, in io.Reader, out io.Writer) {
	err := Handle(ctx, in, out)
	if err != nil {
		fdk.WriteStatus(out, http.StatusInternalServerError)
		out.Write([]byte(err.Error()))
		return
	}
	fdk.WriteStatus(out, http.StatusOK)
}

func Handle(ctx context.Context, in io.Reader, out io.Writer) error {
	defer timeTrack(time.Now(), "ffmpeg-audio-splitter")
	var rp RequestPayload
	if err := json.NewDecoder(in).Decode(&rp); err != nil {
		return err
	}

	ffmpegArgs := strings.Split(fmt.Sprintf(Args, rp.PreSignedURLs.GetURL), " ")
	if DebugMode != 0 {
		ffmpegArgs = append(strings.Split(ArgVerbose, " "), ffmpegArgs...)
	}
	log.Println("ffmpeg ", ffmpegArgs)
	statUUID := uuid.New().String()
	var wavOut bytes.Buffer
	err := runFFMPEG(ctx, ffmpegArgs, &wavOut, statUUID)
	if err != nil {
		return err
	}

	store, err := api.NewFromEndpoint(s3URL)
	if err != nil {
		return err
	}

	err = func() error {
		defer timeTrack(time.Now(), fmt.Sprintf("s3-upload-stat-file-%v", statUUID))
		_, err := store.Uploader.UploadWithContext(ctx, &s3manager.UploadInput{
			Bucket: aws.String(store.Config.Bucket),
			Key:    aws.String(fmt.Sprintf("%v.%v", statUUID, Ext)),
			Body:   &wavOut,
		})
		return err
	}()
	if err != nil {
		return err
	}
	return json.NewEncoder(out).Encode(Response{StatUUID: statUUID})
}
