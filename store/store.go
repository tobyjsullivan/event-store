package store

import (
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/tobyjsullivan/event-store.v3/events"
    "github.com/aws/aws-sdk-go/aws"
    "encoding/base64"
    "bytes"
    "encoding/json"
)

type Store struct {
    s3svc *s3.S3
    bucket string
}

func NewS3Store(svc *s3.S3, bucket string) *Store {
    return &Store{
        s3svc: svc,
        bucket: bucket,
    }
}

type eventFormat struct {
    Prev string `json:"previous"`
    Type string `json:"type"`
    Data string `json:"data"`
}

func (s *Store) Save(e *events.Event) error {
    id := e.ID()
    key := id.String()

    content := &eventFormat{
        Prev: e.PreviousEvent.String(),
        Type: e.Type,
        Data: base64.StdEncoding.EncodeToString(e.Data),
    }

    var buf bytes.Buffer
    encoder := json.NewEncoder(&buf)
    err := encoder.Encode(content)
    if err != nil {
        return err
    }

    _, err = s.s3svc.PutObject(&s3.PutObjectInput{
        Body: bytes.NewReader(buf.Bytes()),
        Bucket: aws.String(s.bucket),
        Key: aws.String(key),
    })
    if err != nil {
        return err
    }

    return nil
}
