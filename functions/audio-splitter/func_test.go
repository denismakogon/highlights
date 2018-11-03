package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/denismakogon/s3-pollster/api"
	"os"
	"testing"
)

func TestFunc(t *testing.T) {
	t.Run("test-ffmpeg-exec", func(t *testing.T) {
		ctx := context.Background()
		var out bytes.Buffer
		f, err := os.Open("payload.json")
		if err != nil {
			t.Fatal(err.Error())
		}
		err = Handle(ctx, f, &out)
		if err != nil {
			t.Fatal(err.Error())
		}
		var id Response
		if err := json.NewDecoder(&out).Decode(&id); err != nil {
			t.Fatal(err.Error())
		}
		store, err := api.NewFromEndpoint(s3URL)
		if err != nil {
			t.Fatal(err.Error())
		}

		store.Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(store.Config.Bucket),
			Key:    aws.String(fmt.Sprintf("%v.%v", id.StatUUID, Ext)),
		})
	})
}
