package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/denismakogon/s3-pollster/api"
	"github.com/google/uuid"
	"os"
	"testing"
)

func TestFunc(t *testing.T) {

	t.Run("sox-exec-csv-test", func(t *testing.T) {
		statUUID := uuid.New().String()
		ctx := context.Background()
		wavFile, err := os.Open(fmt.Sprintf("/tmp/sox.%v", Ext))
		if err != nil {
			t.Fatal(err.Error())
		}

		err = runExec(ctx, fmt.Sprintf(soxExec, statUUID), wavFile, nil, statUUID)
		if err != nil {
			t.Fatal(err.Error())
		}

	})

	t.Run("s3-sox-stream-dft-test", func(t *testing.T) {
		wavFile, err := os.Open(fmt.Sprintf("/tmp/sox.%v", Ext))
		if err != nil {
			t.Fatal(err.Error())
		}
		ctx := context.Background()

		store, err := api.NewFromEndpoint(s3URL)
		if err != nil {
			t.Fatal(err.Error())
		}
		statUUID := uuid.New().String()
		_, err = store.Uploader.UploadWithContext(ctx, &s3manager.UploadInput{
			Bucket: aws.String(store.Config.Bucket),
			Key:    aws.String(fmt.Sprintf("%v.%v", statUUID, Ext)),
			Body:   wavFile,
		})
		if err != nil {
			t.Fatal(err.Error())
		}

		var in, out bytes.Buffer

		json.NewEncoder(&in).Encode(Response{StatUUID: statUUID})
		err = Handle(ctx, &in, &out)
		if err != nil {
			t.Fatal(err.Error())
		}

		store.Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(store.Config.Bucket),
			Key:    aws.String(fmt.Sprintf("%v.%v", statUUID, Ext)),
		})
		store.Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(store.Config.Bucket),
			Key:    aws.String(fmt.Sprintf("%v.dat", statUUID)),
		})
	})
}
