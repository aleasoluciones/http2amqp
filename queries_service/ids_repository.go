// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package queries_service

import (
	"code.google.com/p/go-uuid/uuid"
)

type Id string

type IdsRepository interface {
	Next() Id
}

func NewIdsRepository() IdsRepository {
	return &idsRepository{}
}

type idsRepository struct {
}

func (repo *idsRepository) Next() Id {
	return Id(uuid.New())
}
