// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package queries_service

type Id int

type IdsRepository interface {
	Next() Id
}

func NewIdsRepository() IdsRepository {
	return &idsRepository{}
}

type idsRepository struct {
}

func (repo *idsRepository) Next() Id {
	return Id(1)
}
