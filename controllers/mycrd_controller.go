/*
Copyright 2022.

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

package controllers

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	sofiworkerv1alpha1 "crd/api/v1alpha1"
)

// MyCrdReconciler reconciles a MyCrd object
type MyCrdReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=sofiworker.sofiworker.me,resources=mycrds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=sofiworker.sofiworker.me,resources=mycrds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=sofiworker.sofiworker.me,resources=mycrds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyCrd object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *MyCrdReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	klog.Info("-----------hello-------------")
	mycrd := &sofiworkerv1alpha1.MyCrd{}
	if err := r.Get(ctx, req.NamespacedName, mycrd); err != nil {
		if errors.IsNotFound(err) {
			klog.Info("------mycrd resource is not found------")
			return ctrl.Result{}, nil
		} else {
			klog.Info("you should update mycrd, name: " + mycrd.Name)
			return ctrl.Result{}, err
		}
	}
	klog.Info("--------------", mycrd)
	klog.Info("---------------------------------------------------------")

	for i := int64(0); i < mycrd.Spec.Replicas; i++ {
		pod := &v1.Pod{}
		podName := fmt.Sprintf("%s-%s-%d", mycrd.Name, mycrd.Spec.Image, i)
		if err := r.Get(ctx, types.NamespacedName{Namespace: mycrd.Namespace, Name: podName}, pod); err != nil && errors.IsNotFound(err) {
			pod = &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      podName,
					Namespace: mycrd.Namespace,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  fmt.Sprintf("%s-%d", mycrd.Spec.Image, i),
							Image: mycrd.Spec.Image,
						},
					},
				},
			}
			if err := r.Create(ctx, pod); err != nil {
				klog.Info("-----", err)
				return ctrl.Result{}, err
			}
		}
	}
	klog.Info("create pod success")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyCrdReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&sofiworkerv1alpha1.MyCrd{}).
		Complete(r)
}
