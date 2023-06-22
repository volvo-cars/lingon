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
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	natsv1 "github.com/volvo-cars/nope/api/v1"
	v1 "github.com/volvo-cars/nope/api/v1"
	"github.com/volvo-cars/nope/internal/bla"
)

// UserReconciler reconciles a User object
type UserReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=nope.volvocars.com,resources=users,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nope.volvocars.com,resources=users/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nope.volvocars.com,resources=users/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the User object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *UserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	finalizer := "user.nope.volvocars.com/finalizer"

	var user v1.User
	if err := r.Get(ctx, req.NamespacedName, &user); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	if user.Spec.Account == "" {
		return ctrl.Result{}, fmt.Errorf("user %s/%s has no account", user.Namespace, user.Name)
	}

	refAccount := types.NamespacedName{
		Namespace: user.Namespace,
		Name:      user.Spec.Account,
	}

	var account v1.Account
	if err := r.Get(ctx, refAccount, &account); err != nil {
		return ctrl.Result{}, fmt.Errorf("could not find account %s/%s for user %s/%s", refAccount.Namespace, refAccount.Name, user.Namespace, user.Name)
	}

	// Check that account is ready
	if account.Status.NKeySeed == nil {
		return ctrl.Result{
			RequeueAfter: time.Second * 5,
		}, fmt.Errorf("account %s/%s is not ready for user %s/%s", refAccount.Namespace, refAccount.Name, user.Namespace, user.Name)
	}

	// Check if the user is being deleted
	if !user.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("user is being deleted")
		// Account is being deleted
		if controllerutil.ContainsFinalizer(&user, finalizer) {
			// TODO: to revoke a user we need to add the user's JWT subject
			// (i.e. public key) to the list of revoked users in the account
			// JWT. This is not implemented yet.
			logger.Info("TODO: revoke user")

			controllerutil.RemoveFinalizer(&user, finalizer)
			if err := r.Update(ctx, &user); err != nil {
				return ctrl.Result{}, fmt.Errorf("updating user: %w", err)
			}
		}
		logger.Info("finish reconciliation as the user is being deleted")
		// Finish reconciliation as the object is being deleted
		return ctrl.Result{}, nil
	}

	var managedUser *bla.User
	if user.Status.ID != "" {
		managedUser = &bla.User{
			ID:   user.Status.ID,
			NKey: user.Status.NKeySeed,
			JWT:  user.Status.JWT,
		}
	}

	userReq := bla.UserRequest{
		Name: "user",
	}
	syncdUser, err := bla.SyncUser(account.Status.NKeySeed, managedUser, userReq)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("syncing NATS user %s/%s: %w", user.Namespace, user.Name, err)
	}

	user.Status = v1.UserStatus{
		ID:       syncdUser.ID,
		NKeySeed: syncdUser.NKey,
		JWT:      syncdUser.JWT,
	}
	if err := r.Status().Update(ctx, &user); err != nil {
		return ctrl.Result{}, fmt.Errorf("updating user status %s/%s: %w", user.Namespace, user.Name, err)
	}

	logger.Info("user successfully reconciled", "user_id", user.Status.ID)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&natsv1.User{}).
		Complete(r)
}
