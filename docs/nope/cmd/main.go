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

package main

import (
	"errors"
	"flag"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/nats-io/nkeys"

	natsv1 "github.com/volvo-cars/nope/api/v1"
	"github.com/volvo-cars/nope/internal/controller"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(natsv1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	// Initialize the logger before doing any validation
	// of environment variables or command line args.
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// Ensure NATS-related environment variables are set.
	var cErr error
	natsURL, ok := os.LookupEnv("NATS_URL")
	if !ok {
		cErr = errors.Join(cErr, errors.New("NATS_URL not set"))
	}
	natsCreds, ok := os.LookupEnv("NATS_CREDS")
	if !ok {
		cErr = errors.Join(cErr, errors.New("NATS_CREDS not set"))
	}
	operatorSeed, ok := os.LookupEnv("NATS_OPERATOR_SEED")
	if !ok {
		cErr = errors.Join(cErr, errors.New("NATS_OPERATOR_SEED not set"))
	}
	if cErr != nil {
		setupLog.Error(cErr, "required environment variables not set")
		os.Exit(1)
	}
	operatorNKey, err := os.ReadFile(operatorSeed)
	if err != nil {
		setupLog.Error(err, "reading NATS_OPERATOR_SEED")
		os.Exit(1)
	}
	// Try to parse the operator seed to ensure it is valid
	if _, err := nkeys.FromSeed(operatorNKey); err != nil {
		setupLog.Error(err, "invalid NATS_OPERATOR_SEED")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "6c0123a9.volvocars.com",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controller.AccountReconciler{
		Client:       mgr.GetClient(),
		Scheme:       mgr.GetScheme(),
		NATSURL:      natsURL,
		NATSCreds:    natsCreds,
		OperatorNKey: operatorNKey,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Account")
		os.Exit(1)
	}
	if err = (&controller.UserReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "User")
		os.Exit(1)
	}
	if err = (&controller.StreamReconciler{
		Client:    mgr.GetClient(),
		Scheme:    mgr.GetScheme(),
		NATSURL:   natsURL,
		NATSCreds: natsCreds,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Stream")
		os.Exit(1)
	}
	// TODO: to re-nable webhooks, uncomment this.
	// Cannot get it work locally, even with `export ENABLE_WEBHOOKS=false`
	// if err = (&natsv1.Stream{}).SetupWebhookWithManager(mgr); err != nil {
	// 	setupLog.Error(err, "unable to create webhook", "webhook", "Stream")
	// 	os.Exit(1)
	// }
	if err = (&controller.ConsumerReconciler{
		Client:    mgr.GetClient(),
		Scheme:    mgr.GetScheme(),
		NATSURL:   natsURL,
		NATSCreds: natsCreds,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Consumer")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
