package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"github.com/google/uuid"
	"os"
	"testing"
)

func TestCSVConverter(t *testing.T) {
	t.Run("test-cvs-converter-pipeline", func(t *testing.T) {
		soxStat, err := os.Open("/tmp/sox.dat")
		if err != nil {
			t.Fatal(err.Error())
		}
		var out bytes.Buffer
		err = runExec(context.Background(), toCSVFormat, soxStat, &out, uuid.New().String())
		if err != nil {
			t.Fatal(err.Error())
		}
		_, err = csv.NewReader(&out).Read()
		if err != nil {
			t.Log("unable to validate CSV content from AWK")
			t.Fatal(err.Error())
		}
	})
}
