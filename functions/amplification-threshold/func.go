package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/denismakogon/s3-pollster/common"
	"github.com/fnproject/fdk-go"
)

var (
	theoreticalThreshold = common.WithDefault("THEORETICAL_THRESHOLD", "0.6")
	upperLimitThresholdQuota = common.WithDefault("THRESHOLD_LIMIT_QUOTA", "0.75")
	awkMax = `awk '$2 > %v { print $2 }' | sort -g | awk 'END { print $1*%v }'`
	s3URL  = common.WithDefault("S3_URL",
		"s3://admin:password@docker.for.mac.localhost:9000/us-east-1/default-bucket")
)



func main() {
	fdk.Handle(fdk.HandlerFunc(withError))
}

/* Battle plan:
1. Get reader from an S3
2. Pipe reader from an S3 to AWK
3. Get threshold (75% of highest peak value).
*/
func myHandler(ctx context.Context, in io.Reader, out io.Writer) error {
	var rq Request
	err := json.NewDecoder(in).Decode(&rq)
	if err != nil {
		return err
	}
	defer timeTrack(time.Now(), fmt.Sprintf("handler-%v", rq.StatUUID))

	reader, err := getReader(ctx, &rq)
	if err != nil {
		return err
	}

	threshold, err := ObtainThreshold(ctx, reader, &rq)

	return json.NewEncoder(out).Encode(
		Threshold{
			ThresholdValue: *threshold,
			StatUUID:       rq.StatUUID,
		},
	)
}

func withError(ctx context.Context, in io.Reader, out io.Writer) {
	err := myHandler(ctx, in, out)
	if err != nil {
		fdk.WriteStatus(out, http.StatusInternalServerError)
		out.Write([]byte(err.Error()))
		return
	}
	fdk.WriteStatus(out, http.StatusOK)
}
