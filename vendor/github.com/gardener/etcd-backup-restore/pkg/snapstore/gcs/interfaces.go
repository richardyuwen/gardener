// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gcs

import (
	"context"
	"io"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// Client is storage clint interface for google cloud storage SDK
type Client interface {
	// Bucket ...
	Bucket(name string) BucketHandle
	// Buckets ...
	Buckets(ctx context.Context, projectID string) BucketIterator
	// Close ...
	Close() error

	embedToIncludeNewMethods()
}

// ObjectHandle ...
type ObjectHandle interface {
	// ACL ...
	ACL() ACLHandle
	// Generation ...
	Generation(int64) ObjectHandle
	// If ...
	If(storage.Conditions) ObjectHandle
	// Key ...
	Key([]byte) ObjectHandle
	// ReadCompressed ...
	ReadCompressed(bool) ObjectHandle
	// Attrs ...
	Attrs(context.Context) (*storage.ObjectAttrs, error)
	// Update ...
	Update(context.Context, storage.ObjectAttrsToUpdate) (*storage.ObjectAttrs, error)
	// NewReader ...
	NewReader(context.Context) (Reader, error)
	// NewRangeReader ...
	NewRangeReader(context.Context, int64, int64) (Reader, error)
	// NewWriter ...
	NewWriter(context.Context) Writer
	// Delete ...
	Delete(context.Context) error
	// CopierFrom ...
	CopierFrom(ObjectHandle) Copier
	// ComoserFrom ...
	ComposerFrom(...ObjectHandle) Composer

	embedToIncludeNewMethods()
}

// BucketHandle ...
type BucketHandle interface {
	// Create ...
	Create(context.Context, string, *storage.BucketAttrs) error
	// Delete ...
	Delete(context.Context) error
	// DefaultObjectACL ...
	DefaultObjectACL() ACLHandle
	// Object ...
	Object(string) ObjectHandle
	// Attrs ...
	Attrs(context.Context) (*storage.BucketAttrs, error)
	// Update ...
	Update(context.Context, storage.BucketAttrsToUpdate) (*storage.BucketAttrs, error)
	// If ...
	If(storage.BucketConditions) BucketHandle
	// Objects ...
	Objects(context.Context, *storage.Query) ObjectIterator
	// ACL ...
	ACL() ACLHandle
	// IAM ...
	IAM() *iam.Handle
	// UserProject ...
	UserProject(projectID string) BucketHandle
	// Notifications ...
	Notifications(context.Context) (map[string]*storage.Notification, error)
	// AddNotification ...
	AddNotification(context.Context, *storage.Notification) (*storage.Notification, error)
	// DeleteNotification ...
	DeleteNotification(context.Context, string) error
	// LockRetentionPolicy ...
	LockRetentionPolicy(context.Context) error

	embedToIncludeNewMethods()
}

// ObjectIterator ...
type ObjectIterator interface {
	// Next ...
	Next() (*storage.ObjectAttrs, error)
	// PageInfo ...
	PageInfo() *iterator.PageInfo

	embedToIncludeNewMethods()
}

// BucketIterator ...
type BucketIterator interface {
	// SetPrefix ...
	SetPrefix(string)
	// Next ...
	Next() (*storage.BucketAttrs, error)
	// PageInfo ...
	PageInfo() *iterator.PageInfo

	embedToIncludeNewMethods()
}

// ACLHandle ...
type ACLHandle interface {
	// Delete ...
	Delete(context.Context, storage.ACLEntity) error
	// Set ...
	Set(context.Context, storage.ACLEntity, storage.ACLRole) error
	// List ...
	List(context.Context) ([]storage.ACLRule, error)

	embedToIncludeNewMethods()
}

// Reader ...
type Reader interface {
	io.ReadCloser
	// Size ...
	Size() int64
	// Remain ...
	Remain() int64
	// ContentType ...
	ContentType() string
	// ContentEncoding ...
	ContentEncoding() string
	// CacheControl ...
	CacheControl() string

	embedToIncludeNewMethods()
}

// Writer ...
type Writer interface {
	io.WriteCloser
	// ObjectAttrs ...
	ObjectAttrs() *storage.ObjectAttrs
	// SetChunkSize ...
	SetChunkSize(int)
	// SetProgressFunc ...
	SetProgressFunc(func(int64))
	// SetCRC32C ...
	SetCRC32C(uint32) // Sets both CRC32C and SendCRC32C.
	// CloseWithError ...
	CloseWithError(err error) error
	// Attrs ...
	Attrs() *storage.ObjectAttrs

	embedToIncludeNewMethods()
}

// Copier ...
type Copier interface {
	// ObjectAttrs ...
	ObjectAttrs() *storage.ObjectAttrs
	// SetRewriteToken ...
	SetRewriteToken(string)
	// SetProgressFunc ...
	SetProgressFunc(func(uint64, uint64))
	// SetDestinationKMSKeyName ...
	SetDestinationKMSKeyName(string)
	// Run ...
	Run(context.Context) (*storage.ObjectAttrs, error)

	embedToIncludeNewMethods()
}

// Composer ...
type Composer interface {
	// ObjectAttrs ...
	ObjectAttrs() *storage.ObjectAttrs
	// Run ...
	Run(context.Context) (*storage.ObjectAttrs, error)

	embedToIncludeNewMethods()
}
