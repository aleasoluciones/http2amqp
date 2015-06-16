// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package http2amqp

import (
	"code.google.com/p/go-uuid/uuid"
)

type IdsRepository interface {
	Next() string
}

func NewIdsRepository() IdsRepository {
	return &idsRepository{}
}

type idsRepository struct {
}

func (repo *idsRepository) Next() string {
	return string(uuid.New())
}
