// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package http2amqp

import (
	"code.google.com/p/go-uuid/uuid"
)

// IdsGenerator generate a diferent ID each time Next is called
type IdsGenerator interface {
	Next() string
}

// NewUUIDIdsGenerator return a uuid
func NewUUIDIdsGenerator() IdsGenerator {
	return &uuidGenerator{}
}

type uuidGenerator struct {
}

func (g *uuidGenerator) Next() string {
	return string(uuid.New())
}
