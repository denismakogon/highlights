package main

type PreSignedURLs struct {
	GetURL string `json:"get_url"`
	PutURL string `json:"put_url"`
}

type RequestPayload struct {
	S3Endpoint    string        `json:"s3_endpoint,omitempty"`
	Bucket        string        `json:"bucket,omitempty"`
	Object        string        `json:"object,omitempty"`
	PreSignedURLs PreSignedURLs `json:"presigned_urls"`
}

type Response struct {
	StatUUID string `json:"stat_uuid"`
}
