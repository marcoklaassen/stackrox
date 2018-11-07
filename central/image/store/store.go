package store

import (
	"github.com/boltdb/bolt"
	"github.com/stackrox/rox/generated/api/v1"
	"github.com/stackrox/rox/pkg/bolthelper"
)

const (
	orchShaToRegShaBucket = "orchShaToRegShaBucket"
	imageBucket           = "imageBucket"
	listImageBucket       = "images_list"
)

// Store provides storage functionality for alerts.
//go:generate mockgen-wrapper Store
type Store interface {
	ListImage(sha string) (*v1.ListImage, bool, error)
	ListImages() ([]*v1.ListImage, error)

	GetImages() ([]*v1.Image, error)
	CountImages() (int, error)
	GetImage(sha string) (*v1.Image, bool, error)
	GetImagesBatch(shas []string) ([]*v1.Image, error)

	UpsertImage(image *v1.Image) error
	DeleteImage(sha string) error
}

// New returns a new Store instance using the provided bolt DB instance.
func New(db *bolt.DB) Store {
	bolthelper.RegisterBucketOrPanic(db, imageBucket)
	bolthelper.RegisterBucketOrPanic(db, listImageBucket)
	return &storeImpl{
		db: db,
	}
}
