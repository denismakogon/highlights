# S3/Minio pollster - alternative to S3/Minio object notification

# webhook payload

```json
{
  "bucket": "<bucket>",
  "object": "<object>",
  "presigned_urls": {
    "get_url": "https://...",
    "put_url": "https://..."
  }
}
```

# Env vars

```bash
S3_URL="s3://admin:password@s3:9000/us-east-1/default-bucket"
WEBHOOK_ENDPOINT="https://<your-endpoint>"
POLLSTER_BACKOFF=5
```
