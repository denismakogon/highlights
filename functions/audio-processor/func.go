package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/denismakogon/s3-pollster/api"
	"github.com/denismakogon/s3-pollster/common"
	"github.com/fnproject/fdk-go"
)

func main() {
	fdk.Handle(fdk.HandlerFunc(withError))
}

var (
	Ext     = "wav"
	soxExec = `sox -t wav - /tmp/%v.dat`
	s3URL   = common.WithDefault("S3_URL",
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

/* Battle plan:
1. Get WAV stream from an S3.
2. Run SoX to get amplification peaks.
3. Upload stat file to an S3.
*/
func Handle(ctx context.Context, in io.Reader, out io.Writer) error {
	log.Println("in Handle")
	var rq Response
	err := json.NewDecoder(in).Decode(&rq)
	if err != nil {
		return err
	}
	log.Println("body decoded")
	defer timeTrack(time.Now(), fmt.Sprintf("audio-processor-func-%v", rq.StatUUID))

	store, err := api.NewFromEndpoint(s3URL)
	if err != nil {
		return err
	}

	wavFile, err := func() (io.Reader, error) {
		defer timeTrack(time.Now(), fmt.Sprintf("s3-get-audio-stream-%v", rq.StatUUID))
		wavFile, err := store.Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
			Bucket: aws.String(store.Config.Bucket),
			Key:    aws.String(fmt.Sprintf("%v.%v", rq.StatUUID, Ext)),
		})
		if err != nil {
			return nil, err
		}

		return wavFile.Body, nil
	}()
	if err != nil {
		return err
	}

	err = runExec(ctx, fmt.Sprintf(soxExec, rq.StatUUID), wavFile, nil, rq.StatUUID)
	if err != nil {
		return err
	}

	statFileName := fmt.Sprintf("%v.dat", rq.StatUUID)
	statFile, err := os.Open("/tmp/" + statFileName)
	if err != nil {
		return err
	}

	err = func() error {
		defer timeTrack(time.Now(), fmt.Sprintf("s3-upload-stat-file-%v", rq.StatUUID))
		_, err := store.Uploader.UploadWithContext(ctx, &s3manager.UploadInput{
			Bucket: aws.String(store.Config.Bucket),
			Key:    aws.String(statFileName),
			Body:   statFile,
		})
		return err
	}()
	if err != nil {
		return err
	}

	return json.NewEncoder(out).Encode(rq)
}
