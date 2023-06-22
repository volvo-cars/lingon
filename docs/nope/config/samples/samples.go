package samples

import (
	"github.com/volvo-cars/lingon/pkg/kube"
	v1 "github.com/volvo-cars/nope/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewApp() kube.Exporter {
	return &App{
		Account: &v1.Account{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Account",
				APIVersion: "nope.volvocars.com/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "sample",
				Labels: map[string]string{
					"app.kubernetes.io/name":       "account",
					"app.kubernetes.io/instance":   "account-sample",
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
		},
		User: &v1.User{
			TypeMeta: metav1.TypeMeta{
				Kind:       "User",
				APIVersion: "nope.volvocars.com/v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "sample",
				Labels: map[string]string{
					"app.kubernetes.io/name":       "user",
					"app.kubernetes.io/instance":   "user-sample",
					"app.kubernetes.io/part-of":    "nope",
					"app.kubernetes.io/managed-by": "lingon",
					"app.kubernetes.io/created-by": "nope",
				},
				Finalizers: []string{
					"user.nope.volvocars.com/finalizer",
				},
			},
			Spec: v1.UserSpec{
				Account: "sample",
				Name:    "sample",
			},
		},
	}
}

type App struct {
	kube.App

	Account *v1.Account
	User    *v1.User
}
