// Copyright 2015 The http2amqp Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package queries_service_test

import (
	. "github.com/aleasoluciones/http2amqp/queries_service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ids repository", func() {
	Describe("Next", func() {
		It("returns a string id", func() {
			idsRepository := NewIdsRepository()

			id := idsRepository.Next()

			Expect(id).To(BeAssignableToTypeOf(Id("")))
		})

		It("does not repeat id", func() {
			idsRepository := NewIdsRepository()

			id1 := idsRepository.Next()
			id2 := idsRepository.Next()

			Expect(id1).NotTo(Equal(id2))
		})

		It("does not repeat id in two different instances", func() {
			idsRepository1 := NewIdsRepository()
			idsRepository2 := NewIdsRepository()

			id1 := idsRepository1.Next()
			id2 := idsRepository2.Next()

			Expect(id1).NotTo(Equal(id2))
		})
	})
})
