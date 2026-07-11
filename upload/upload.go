package upload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/objectstorage"
)

type Storage struct {
	namespace string
	bucket    string
	region    string
	client    objectstorage.ObjectStorageClient
}

func NewStorage(bucket string, config common.ConfigurationProvider) (*Storage, error) {
	region, err := common.DefaultConfigProvider().Region()

	if err != nil {
		return nil, err
	}

	client, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(config)

	if err != nil {
		return nil, err
	}

	ctx, timeout := context.WithTimeout(context.Background(), 10*time.Second)

	defer timeout()

	nsReq := objectstorage.GetNamespaceRequest{}
	nsRes, err := client.GetNamespace(ctx, nsReq)

	if err != nil {
		return nil, err
	}

	return &Storage{
		client:    client,
		region:    region,
		bucket:    bucket,
		namespace: *nsRes.Value,
	}, nil
}

func (s *Storage) GetTempUrl(objectName string, hours int) (string, error) {
	ctx, timeout := context.WithTimeout(context.Background(), 20*time.Second)

	defer timeout()

	expiry := time.Now().Add(time.Duration(hours) * time.Hour)

	req := objectstorage.CreatePreauthenticatedRequestRequest{
		NamespaceName: &s.namespace,
		BucketName:    &s.bucket,
		CreatePreauthenticatedRequestDetails: objectstorage.CreatePreauthenticatedRequestDetails{
			Name:        common.String(fmt.Sprintf("temp-%d", time.Now().Unix())),
			ObjectName:  &objectName,
			AccessType:  objectstorage.CreatePreauthenticatedRequestDetailsAccessTypeObjectread,
			TimeExpires: &common.SDKTime{Time: expiry},
		},
	}

	res, err := s.client.CreatePreauthenticatedRequest(ctx, req)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://objectstorage.%s.oraclecloud.com%s", s.region, *res.AccessUri), nil
}

func (s *Storage) UploadFromBytes(objectName string, data []byte) (string, error) {
	ctx, timeout := context.WithTimeout(context.Background(), 20*time.Second)

	defer timeout()

	req := objectstorage.PutObjectRequest{
		NamespaceName: &s.namespace,
		BucketName:    &s.bucket,
		ObjectName:    &objectName,
		ContentLength: common.Int64(int64(len(data))),
		PutObjectBody: io.NopCloser(bytes.NewReader(data)),
	}

	_, err := s.client.PutObject(ctx, req)

	if err != nil {
		return "", err
	}

	return s.GetTempUrl(objectName, 24)
}
