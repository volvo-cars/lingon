// Copyright (c) 2023 Volvo Car Corporation
// SPDX-License-Identifier: Apache-2.0

package samples

import (
	"github.com/volvo-cars/lingon/pkg/kube"
	v1 "github.com/volvo-cars/nope/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewApp() kube.Exporter {
	account := &v1.Account{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Account",
			APIVersion: "nope.volvocars.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "sample",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "sample",
				"app.kubernetes.io/instance":   "sample",
				"app.kubernetes.io/part-of":    "nope",
				"app.kubernetes.io/managed-by": "lingon",
				"app.kubernetes.io/created-by": "nope",
			},
			Finalizers: []string{
				"account.nope.volvocars.com/finalizer",
			},
		},
		Spec: v1.AccountSpec{
			Name: "sample",
		},
	}

	user := &v1.User{
		TypeMeta: metav1.TypeMeta{
			Kind:       "User",
			APIVersion: "nope.volvocars.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "sample",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "sample",
				"app.kubernetes.io/instance":   "sample",
				"app.kubernetes.io/part-of":    "nope",
				"app.kubernetes.io/managed-by": "lingon",
				"app.kubernetes.io/created-by": "nope",
			},
			Finalizers: []string{
				"user.nope.volvocars.com/finalizer",
			},
		},
		Spec: v1.UserSpec{
			Account: account.Name,
			Name:    "sample",
		},
	}

	stream := &v1.Stream{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Stream",
			APIVersion: "nope.volvocars.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "sample",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "sample",
				"app.kubernetes.io/instance":   "sample",
				"app.kubernetes.io/part-of":    "nope",
				"app.kubernetes.io/managed-by": "lingon",
				"app.kubernetes.io/created-by": "nope",
			},
			Finalizers: []string{
				"stream.nope.volvocars.com/finalizer",
			},
		},
		Spec: v1.StreamSpec{
			Account:  account.Name,
			Name:     "sample",
			Subjects: []string{"sample"},
		},
	}

	consumer := &v1.Consumer{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Consumer",
			APIVersion: "nope.volvocars.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "sample",
			Labels: map[string]string{
				"app.kubernetes.io/name":       "sample",
				"app.kubernetes.io/instance":   "sample",
				"app.kubernetes.io/part-of":    "nope",
				"app.kubernetes.io/managed-by": "lingon",
				"app.kubernetes.io/created-by": "nope",
			},
			Finalizers: []string{
				"consumer.nope.volvocars.com/finalizer",
			},
		},
		Spec: v1.ConsumerSpec{
			Stream: stream.Name,
			Name:   "sample",
		},
	}

	return &App{
		Account:  account,
		User:     user,
		Stream:   stream,
		Consumer: consumer,
	}
}

type App struct {
	kube.App

	Account  *v1.Account
	User     *v1.User
	Stream   *v1.Stream
	Consumer *v1.Consumer
}
