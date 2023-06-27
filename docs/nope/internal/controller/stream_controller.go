/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	natsv1 "github.com/volvo-cars/nope/api/v1"
	v1 "github.com/volvo-cars/nope/api/v1"
	"github.com/volvo-cars/nope/internal/bla"
)

// StreamReconciler reconciles a Stream object
type StreamReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	// NATS related configs
	NATSURL      string
	NATSCreds    string
	OperatorNKey []byte
}

//+kubebuilder:rbac:groups=nope.volvocars.com,resources=streams,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nope.volvocars.com,resources=streams/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nope.volvocars.com,resources=streams/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Stream object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *StreamReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Get the Stream resource.
	var stream v1.Stream
	if err := r.Get(ctx, req.NamespacedName, &stream); err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "unable to fetch Stream")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Get the Account resource referenced by this stream.
	targetAccount := types.NamespacedName{
		Namespace: stream.Namespace,
		Name:      stream.Spec.Account,
	}
	var account v1.Account
	if err := r.Get(ctx, targetAccount, &account); err != nil {
		if client.IgnoreNotFound(err) != nil {
			logger.Error(err, "unable to fetch Account for Stream")
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if account.Status.ServiceUser == nil {
		logger.Info("Account does not have a service user yet")
		return ctrl.Result{
			RequeueAfter: time.Second * 2,
		}, nil
	}

	// Get account server user and create a connection to the NATS server.
	userJWT := account.Status.ServiceUser.JWT
	userNKey := account.Status.ServiceUser.NKeySeed
	userKeyPair, err := nkeys.FromSeed(userNKey)
	if err != nil {
		logger.Error(err, "unable to create nkey pair from seed")
		return ctrl.Result{}, fmt.Errorf("getting key pair from seed: %w", err)
	}

	nc, err := nats.Connect(
		r.NATSURL,
		bla.UserJWTOption(userJWT, userKeyPair),
	)
	if err != nil {
		logger.Error(err, "unable to connect to NATS server")
		return ctrl.Result{}, fmt.Errorf("connecting to NATS server: %w", err)
	}

	var managedStream *bla.Stream
	if stream.Status.Name != "" {
		managedStream = &bla.Stream{
			Name: stream.Status.Name,
		}
	}

	syncdStream, err := bla.SyncStream(nc, managedStream, bla.StreamRequest{
		Name:     stream.Spec.Name,
		Subjects: stream.Spec.Subjects,
	})
	if err != nil {
		logger.Error(err, "unable to sync stream")
		return ctrl.Result{}, fmt.Errorf("syncing stream: %w", err)
	}

	stream.Status.Name = syncdStream.Name
	if err := r.Status().Update(ctx, &stream); err != nil {
		logger.Error(err, "unable to update Stream status")
		return ctrl.Result{}, fmt.Errorf("updating Stream status: %w", err)
	}
	logger.Info("strean successfully reconciled", "name", syncdStream.Name)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StreamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&natsv1.Stream{}).
		Complete(r)
}
