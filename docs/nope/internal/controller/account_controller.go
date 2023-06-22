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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/nats-io/nats.go"
	natsv1 "github.com/volvo-cars/nope/api/v1"
	v1 "github.com/volvo-cars/nope/api/v1"
	"github.com/volvo-cars/nope/internal/bla"
)

// AccountReconciler reconciles a Account object
type AccountReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	// NATS related configs
	NATSURL      string
	NATSCreds    string
	OperatorNKey []byte
}

//+kubebuilder:rbac:groups=nope.volvocars.com,resources=accounts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=nope.volvocars.com,resources=accounts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=nope.volvocars.com,resources=accounts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Account object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *AccountReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	finalizer := "account.nope.volvocars.com/finalizer"

	var account v1.Account
	if err := r.Get(ctx, req.NamespacedName, &account); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	nc, err := nats.Connect(r.NATSURL, nats.UserCredentials(r.NATSCreds))
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("connecting to NATS: %w", err)
	}

	// Check if the account is being deleted
	if !account.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info("account is being deleted")
		// Account is being deleted
		if controllerutil.ContainsFinalizer(&account, finalizer) {
			if err := bla.DeleteAccount(nc, r.OperatorNKey, account.Status.ID); err != nil {
				return ctrl.Result{}, fmt.Errorf("deleting NATS account: %w", err)
			}
			controllerutil.RemoveFinalizer(&account, finalizer)
			if err := r.Update(ctx, &account); err != nil {
				return ctrl.Result{}, fmt.Errorf("updating account: %w", err)
			}
			logger.Info("account deleted and finalizer removed")
		}
		logger.Info("finish reconciliation as the account is being deleted")
		// Finish reconciliation as the object is being deleted
		return ctrl.Result{}, nil
	}

	logger.Info("checking NATS account")
	// Else proceed with reconciliation
	var managedAccount *bla.Account
	if account.Status.ID != "" {
		managedAccount = &bla.Account{
			ID:   account.Status.ID,
			NKey: account.Status.NKeySeed,
			JWT:  account.Status.JWT,
		}
	}

	// Create the account request an sync the account
	accountReq := bla.AccountRequest{
		Name: account.Name,
	}
	syncdAccount, err := bla.SyncAccount(nc, r.OperatorNKey, managedAccount, accountReq)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("syncing NATS account %s/%s: %w", account.Namespace, account.Name, err)
	}
	logger.Info("checking NATS account service user")
	// We also need to sync the account user, which we use to manage streams and consumers
	// in the account
	var serviceUser *bla.User
	if account.Status.ServiceUser != nil {
		serviceUser = &bla.User{
			ID:   account.Status.ServiceUser.ID,
			NKey: account.Status.ServiceUser.NKeySeed,
			JWT:  account.Status.ServiceUser.JWT,
		}
	}
	userReq := bla.UserRequest{
		Name: "nope-service-user",
	}
	syncdUser, err := bla.SyncUser(syncdAccount.NKey, serviceUser, userReq)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("syncing NATS service user for account %s/%s: %w", account.Namespace, account.Name, err)
	}
	account.Status = v1.AccountStatus{
		ID:       syncdAccount.ID,
		NKeySeed: syncdAccount.NKey,
		JWT:      syncdAccount.JWT,
		ServiceUser: &v1.AccountServiceUser{
			ID:       syncdUser.ID,
			NKeySeed: syncdUser.NKey,
			JWT:      syncdUser.JWT,
		},
	}
	if err := r.Status().Update(ctx, &account); err != nil {
		return ctrl.Result{}, fmt.Errorf("updating account status %s/%s: %w", account.Namespace, account.Name, err)
	}

	logger.Info("account successfully reconciled", "account_id", syncdAccount.ID)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AccountReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&natsv1.Account{}).
		Complete(r)
}
