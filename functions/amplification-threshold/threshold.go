package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func ObtainThreshold(ctx context.Context, input io.Reader, rq *Request) (*Threshold, error) {
	defer timeTrack(time.Now(), fmt.Sprintf("threshold-obtain-%v", rq.StatUUID))
	var buf bytes.Buffer
	err := runExec(ctx, fmt.Sprintf(awkMax, theoreticalThreshold, upperLimitThresholdQuota),
		input, &buf, rq.StatUUID)
	if err != nil {
		return nil, err
	}
	threshold, err := strconv.ParseFloat(
		strings.TrimRight(buf.String(), "\n"), 64)
	if err != nil {
		return nil, err
	}

	return &Threshold{threshold, rq.StatUUID}, nil
}
