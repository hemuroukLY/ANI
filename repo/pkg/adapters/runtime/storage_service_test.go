package runtime

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/kubercloud/ani/pkg/ports"
)

func TestLocalStorageServiceVolumeDevProfile(t *testing.T) {
	service := NewLocalStorageService()
	volume, err := service.CreateVolume(context.Background(), ports.StorageVolumeCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "storage-volume-a",
		Name:           "data-a",
		SizeGiB:        100,
		StorageClass:   "fast",
	})
	if err != nil {
		t.Fatalf("CreateVolume() error = %v", err)
	}
	if volume.VolumeID == "" || volume.State != ports.StorageResourceAvailable || volume.StorageClass != "fast" {
		t.Fatalf("volume = %#v, want available fast volume", volume)
	}
	replay, err := service.CreateVolume(context.Background(), ports.StorageVolumeCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "storage-volume-a",
		Name:           "data-a-retry",
		SizeGiB:        200,
		StorageClass:   "slow",
	})
	if err != nil {
		t.Fatalf("CreateVolume replay error = %v", err)
	}
	if replay.VolumeID != volume.VolumeID || replay.SizeGiB != volume.SizeGiB {
		t.Fatalf("replay volume = %#v, want original %#v", replay, volume)
	}
	if _, err := service.GetVolume(context.Background(), ports.StorageResourceGetRequest{TenantID: "tenant-b", ResourceID: volume.VolumeID}); err == nil {
		t.Fatalf("GetVolume from another tenant succeeded, want isolation error")
	}
	deleted, err := service.DeleteVolume(context.Background(), ports.StorageResourceGetRequest{TenantID: "tenant-a", ResourceID: volume.VolumeID})
	if err != nil {
		t.Fatalf("DeleteVolume() error = %v", err)
	}
	if deleted.State != ports.StorageResourceDeleted {
		t.Fatalf("deleted state = %q, want deleted", deleted.State)
	}
}

func TestLocalStorageServiceFilesystemAndObjectDevProfile(t *testing.T) {
	service := NewLocalStorageService()
	filesystem, err := service.CreateFilesystem(context.Background(), ports.StorageFilesystemCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "storage-fs-a",
		Name:           "shared",
		Protocol:       "cephfs",
		SizeGiB:        500,
	})
	if err != nil {
		t.Fatalf("CreateFilesystem() error = %v", err)
	}
	if filesystem.FilesystemID == "" || filesystem.Protocol != "cephfs" || filesystem.Endpoint == "" {
		t.Fatalf("filesystem = %#v, want cephfs endpoint", filesystem)
	}
	object, err := service.CreateObject(context.Background(), ports.StorageObjectCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "storage-object-a",
		Bucket:         "models",
		Key:            "llm/model.bin",
		SizeBytes:      1024,
		ContentType:    "application/octet-stream",
	})
	if err != nil {
		t.Fatalf("CreateObject() error = %v", err)
	}
	if object.ObjectID == "" || object.State != ports.StorageResourceAvailable || object.Bucket != "models" {
		t.Fatalf("object = %#v, want available object metadata", object)
	}
	objects, err := service.ListObjects(context.Background(), ports.StorageResourceListRequest{TenantID: "tenant-a"})
	if err != nil {
		t.Fatalf("ListObjects() error = %v", err)
	}
	if len(objects) != 1 {
		t.Fatalf("objects = %d, want 1", len(objects))
	}
}

func TestLocalStorageServiceSnapshotsAndMountTargets(t *testing.T) {
	service := NewLocalStorageService()
	volume, err := service.CreateVolume(context.Background(), ports.StorageVolumeCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "snapshot-volume-a",
		Name:           "db-data",
		SizeGiB:        8,
	})
	if err != nil {
		t.Fatalf("CreateVolume error = %v", err)
	}
	snapshot, err := service.CreateVolumeSnapshot(context.Background(), ports.VolumeSnapshotCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "snapshot-a",
		VolumeID:       volume.VolumeID,
		Name:           "db-data-snap",
		Description:    "daily backup",
	})
	if err != nil {
		t.Fatalf("CreateVolumeSnapshot error = %v", err)
	}
	retry, err := service.CreateVolumeSnapshot(context.Background(), ports.VolumeSnapshotCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "snapshot-a",
		VolumeID:       volume.VolumeID,
		Name:           "changed-name",
	})
	if err != nil {
		t.Fatalf("CreateVolumeSnapshot retry error = %v", err)
	}
	if retry.SnapshotID != snapshot.SnapshotID || retry.Name != snapshot.Name {
		t.Fatalf("idempotent snapshot = %+v, want original %+v", retry, snapshot)
	}
	snapshots, err := service.ListVolumeSnapshots(context.Background(), ports.VolumeSnapshotListRequest{TenantID: "tenant-a", VolumeID: volume.VolumeID})
	if err != nil {
		t.Fatalf("ListVolumeSnapshots error = %v", err)
	}
	if len(snapshots) != 1 || snapshots[0].Status != ports.VolumeSnapshotAvailable {
		t.Fatalf("snapshots = %+v, want one available snapshot", snapshots)
	}

	filesystem, err := service.CreateFilesystem(context.Background(), ports.StorageFilesystemCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "mount-fs-a",
		Name:           "shared",
		SizeGiB:        32,
	})
	if err != nil {
		t.Fatalf("CreateFilesystem error = %v", err)
	}
	targets, err := service.ListFilesystemMountTargets(context.Background(), ports.FilesystemMountTargetListRequest{
		TenantID:     "tenant-a",
		FilesystemID: filesystem.FilesystemID,
	})
	if err != nil {
		t.Fatalf("ListFilesystemMountTargets error = %v", err)
	}
	if len(targets) != 1 || targets[0].FilesystemID != filesystem.FilesystemID || targets[0].Status != ports.MountTargetAvailable {
		t.Fatalf("mount targets = %+v, want generated available target", targets)
	}
}

func TestLocalStorageServiceBucketsAndSignedObjectURLsUseObjectStorePort(t *testing.T) {
	objectStore := &fakeObjectStore{
		uploadURL:   "https://objects.local/upload/model.bin",
		downloadURL: "https://objects.local/download/model.bin",
		expiresAt:   time.Date(2026, 6, 19, 10, 30, 0, 0, time.UTC),
	}
	service := NewLocalStorageService(WithStorageObjectStore(objectStore))

	bucket, err := service.CreateStorageBucket(context.Background(), ports.StorageBucketCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "bucket-a",
		Name:           "models-a",
		Region:         "local",
		AccessMode:     "private",
	})
	if err != nil {
		t.Fatalf("CreateStorageBucket() error = %v", err)
	}
	if bucket.BucketID == "" || bucket.Name != "models-a" || bucket.AccessMode != "private" {
		t.Fatalf("bucket = %#v, want private models-a bucket", bucket)
	}
	if objectStore.ensureBucket != ports.BucketClass("models-a") {
		t.Fatalf("EnsureBucket class = %q, want models-a", objectStore.ensureBucket)
	}

	replay, err := service.CreateStorageBucket(context.Background(), ports.StorageBucketCreateRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "bucket-a",
		Name:           "changed-name",
	})
	if err != nil {
		t.Fatalf("CreateStorageBucket replay error = %v", err)
	}
	if replay.BucketID != bucket.BucketID || replay.Name != bucket.Name {
		t.Fatalf("replay bucket = %#v, want original %#v", replay, bucket)
	}

	upload, err := service.CreateStorageObjectUpload(context.Background(), ports.StorageObjectUploadRequest{
		TenantID:       "tenant-a",
		IdempotencyKey: "upload-a",
		BucketID:       bucket.BucketID,
		Key:            "llm/model.bin",
		ContentType:    "application/octet-stream",
	})
	if err != nil {
		t.Fatalf("CreateStorageObjectUpload() error = %v", err)
	}
	if upload.UploadURL != objectStore.uploadURL || upload.ObjectID == "" || upload.ExpiresAt != objectStore.expiresAt {
		t.Fatalf("upload = %#v, want signed upload response", upload)
	}
	if objectStore.uploadRef.BucketClass != ports.BucketClass("models-a") || objectStore.uploadRef.ObjectKey != "llm/model.bin" {
		t.Fatalf("upload ref = %#v, want bucket models-a key llm/model.bin", objectStore.uploadRef)
	}

	download, err := service.GetStorageObjectDownload(context.Background(), ports.StorageObjectDownloadRequest{
		TenantID:       "tenant-a",
		ObjectID:       upload.ObjectID,
		ExpiresSeconds: 600,
	})
	if err != nil {
		t.Fatalf("GetStorageObjectDownload() error = %v", err)
	}
	if download.DownloadURL != objectStore.downloadURL || download.ContentType != "application/octet-stream" {
		t.Fatalf("download = %#v, want signed download response", download)
	}
	if objectStore.downloadRef.BucketClass != ports.BucketClass("models-a") || objectStore.downloadRef.ObjectKey != "llm/model.bin" {
		t.Fatalf("download ref = %#v, want bucket models-a key llm/model.bin", objectStore.downloadRef)
	}

	buckets, err := service.ListStorageBuckets(context.Background(), ports.StorageResourceListRequest{TenantID: "tenant-a"})
	if err != nil {
		t.Fatalf("ListStorageBuckets() error = %v", err)
	}
	if len(buckets) != 1 || buckets[0].ObjectCount != 1 {
		t.Fatalf("buckets = %#v, want one bucket with one object", buckets)
	}
}

type fakeObjectStore struct {
	ensureBucket ports.BucketClass
	uploadRef    ports.ObjectRef
	downloadRef  ports.ObjectRef
	uploadURL    string
	downloadURL  string
	expiresAt    time.Time
}

func (s *fakeObjectStore) EnsureBucket(_ context.Context, class ports.BucketClass) error {
	s.ensureBucket = class
	return nil
}

func (s *fakeObjectStore) PutObject(context.Context, ports.PutObjectInput) (ports.ObjectMetadata, error) {
	return ports.ObjectMetadata{}, ports.ErrUnsupported
}

func (s *fakeObjectStore) GetObject(context.Context, ports.ObjectRef) (io.ReadCloser, ports.ObjectMetadata, error) {
	return nil, ports.ObjectMetadata{}, ports.ErrUnsupported
}

func (s *fakeObjectStore) DeleteObject(context.Context, ports.ObjectRef) error {
	return ports.ErrUnsupported
}

func (s *fakeObjectStore) StatObject(context.Context, ports.ObjectRef) (ports.ObjectMetadata, error) {
	return ports.ObjectMetadata{}, ports.ErrUnsupported
}

func (s *fakeObjectStore) SignedUploadURL(_ context.Context, ref ports.ObjectRef, _ time.Duration) (ports.SignedURL, error) {
	s.uploadRef = ref
	return ports.SignedURL{URL: s.uploadURL, ExpiresAt: s.expiresAt}, nil
}

func (s *fakeObjectStore) SignedDownloadURL(_ context.Context, ref ports.ObjectRef, _ time.Duration) (ports.SignedURL, error) {
	s.downloadRef = ref
	return ports.SignedURL{URL: s.downloadURL, ExpiresAt: s.expiresAt}, nil
}
