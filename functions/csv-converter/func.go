package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type Request struct {
	ThresholdValue float64 `json:"threshold_value"`
	StatUUID       string  `json:"stat_uuid"`
}

var (
	toCSVFormat = `awk '{ if (NR>2) print $1","$2}'`
	s3URL       = common.WithDefault("S3_URL",
		"s3://admin:password@docker.for.mac.localhost:9000/us-east-1/default-bucket")
)

func Handler(ctx context.Context, in io.Reader, out io.Writer) error {
	var rq Request
	err := json.NewDecoder(in).Decode(&rq)
	if err != nil {
		return err
	}
	defer timeTrack(time.Now(), fmt.Sprintf("handler-%v", rq.StatUUID))

	store, err := api.NewFromEndpoint(s3URL)
	if err != nil {
		return err
	}

	reader, err := func() (io.Reader, error) {
		defer timeTrack(time.Now(), fmt.Sprintf("s3-download-stat-file-%v", rq.StatUUID))
		obj, err := store.Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
			Bucket: aws.String(store.Config.Bucket),
			Key:    aws.String(fmt.Sprintf("%v.dat", rq.StatUUID)),
		})
		if err != nil {
			return nil, err
		}
		return obj.Body, nil
	}()
	if err != nil {
		return err
	}

	var csvContent bytes.Buffer
	err = runExec(ctx, toCSVFormat, reader, &csvContent, rq.StatUUID)
	if err != nil {
		return err
	}

	err = func() error {
		defer timeTrack(time.Now(), fmt.Sprintf("s3-upload-csv-stat-file-%v", rq.StatUUID))
		_, err := store.Uploader.UploadWithContext(ctx, &s3manager.UploadInput{
			Bucket: aws.String(store.Config.Bucket),
			Key:    aws.String(fmt.Sprintf("%v.csv", rq.StatUUID)),
			Body:   &csvContent,
		})
		return err
	}()
	if err != nil {
		return err
	}

	return json.NewEncoder(out).Encode(rq)
}

func withError(ctx context.Context, in io.Reader, out io.Writer) {
	err := Handler(ctx, in, out)
	if err != nil {
		fdk.WriteStatus(out, http.StatusInternalServerError)
		out.Write([]byte(err.Error()))
		return
	}
	fdk.WriteStatus(out, http.StatusOK)
}
