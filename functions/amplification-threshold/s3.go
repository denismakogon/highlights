package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/denismakogon/s3-pollster/api"
	"io"
	"time"
)

func getReader(ctx context.Context, rq *Request) (io.Reader, error) {
	defer timeTrack(time.Now(), "s3-get-reader")
	store, err := api.NewFromEndpoint(s3URL)
	if err != nil {
		return nil, err
	}

	obj, err := store.Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(store.Config.Bucket),
		Key:    aws.String(fmt.Sprintf("%v.dat", rq.StatUUID)),
	})

	return obj.Body, err
}
