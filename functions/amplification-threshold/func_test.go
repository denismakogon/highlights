package main

import (
	"context"
	"github.com/google/uuid"
	"os"
	"testing"
)

func TestAmplificationThreshold(t *testing.T) {
	t.Run("test-get-amp-threshold", func(t *testing.T) {
		soxStat, err := os.Open("/tmp/sox.dat")
		if err != nil {
			t.Fatal(err.Error())
		}
		thresholdPtr, err := ObtainThreshold(
			context.Background(),
			soxStat,
			&Request{StatUUID: uuid.New().String()},
		)
		if err != nil {
			t.Fatal(err.Error())
		}
		t.Log("threshold: ", *thresholdPtr)
		if *thresholdPtr <= 0 {
			t.Fatal("invalid threshold value, suppose be greater than 0.0")
		}
	})
}
