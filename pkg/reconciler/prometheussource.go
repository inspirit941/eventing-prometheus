/*
Copyright 2019 The Knative Authors

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

package reconciler

import (
	"context"
	"fmt"

	"knative.dev/eventing/pkg/utils"

	"k8s.io/client-go/kubernetes"
	"knative.dev/eventing/pkg/reconciler/source"
	"knative.dev/pkg/controller"

	"github.com/inspirit941/eventing-prometheus/pkg/reconciler/resources"
	"github.com/kelseyhightower/envconfig"
	"github.com/robfig/cron"
	"go.uber.org/zap"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1listers "k8s.io/client-go/listers/apps/v1"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
	"knative.dev/pkg/resolver"

	"github.com/inspirit941/eventing-prometheus/pkg/apis/sources/v1alpha1"
	promreconciler "github.com/inspirit941/eventing-prometheus/pkg/client/injection/reconciler/sources/v1alpha1/prometheussource"
)

const (
	// Name of the corev1.Events emitted from the reconciliation process
	prometheussourceDeploymentCreated = "PrometheusSourceDeploymentCreated"
	prometheussourceDeploymentUpdated = "PrometheusSourceDeploymentUpdated"
	prometheussourceDeploymentDeleted = "PrometheusSourceDeploymentDeleted"
)

type envConfig struct {
	Image string `envconfig:"PROMETHEUS_RA_IMAGE" required:"true"`
}

// Reconciler reconciles a PrometheusSource object
type Reconciler struct {
	kubeClientSet kubernetes.Interface

	// listers index properties about resources
	deploymentLister appsv1listers.DeploymentLister

	sinkResolver *resolver.URIResolver
	configs      *source.ConfigWatcher
}

var _ promreconciler.Interface = (*Reconciler)(nil)

func (r *Reconciler) ReconcileKind(ctx context.Context, source *v1alpha1.PrometheusSource) pkgreconciler.Event {
	source.Status.InitializeConditions()

	if source.Spec.Sink == nil {
		source.Status.MarkNoSink("SinkMissing", "")
		return fmt.Errorf("spec.sink missing")
	}

	dest := source.Spec.Sink.DeepCopy()
	if dest.Ref != nil {
		// To call URIFromDestination(), dest.Ref must have a Namespace. If there is
		// no Namespace defined in dest.Ref, we will use the Namespace of the source
		// as the Namespace of dest.Ref.
		if dest.Ref.Namespace == "" {
			dest.Ref.Namespace = source.GetNamespace()
		}
	}

	sinkURI, err := r.sinkResolver.URIFromDestinationV1(ctx, *dest, source)
	if err != nil {
		source.Status.MarkNoSink("NotFound", "")
		return err
	}
	source.Status.MarkSink(sinkURI)

	_, err = cron.ParseStandard(source.Spec.Schedule)
	if err != nil {
		source.Status.MarkInvalidSchedule("Invalid", "Reason: "+err.Error())
		return fmt.Errorf("invalid schedule: %v", err)
	}
	source.Status.MarkValidSchedule()

	ra, err := r.createReceiveAdapter(ctx, source, sinkURI)
	if err != nil {
		logging.FromContext(ctx).Errorw("Unable to create the receive adapter", zap.Error(err))
		return err
	}
	// Update source status// Update source status
	source.Status.PropagateDeploymentAvailability(ra)

	source.Status.CloudEventAttributes = []duckv1.CloudEventAttributes{{
		Type:   v1alpha1.PromQLPrometheusSourceEventType,
		Source: r.makeEventSource(source),
	}}

	return nil
}

func (r *Reconciler) createReceiveAdapter(ctx context.Context, src *v1alpha1.PrometheusSource, sinkURI *apis.URL) (*appsv1.Deployment, error) {
	eventSource := r.makeEventSource(src)
	logging.FromContext(ctx).Debug("event source", zap.Any("source", eventSource))

	env := &envConfig{}
	if err := envconfig.Process("", env); err != nil {
		logging.FromContext(ctx).Panicf("required environment variable is not defined: %v", err)
	}

	adapterArgs := resources.ReceiveAdapterArgs{
		EventSource:    eventSource,
		Image:          env.Image,
		Source:         src,
		Labels:         resources.Labels(src.Name),
		SinkURI:        sinkURI.String(),
		AdditionalEnvs: r.configs.ToEnvVars(),
	}
	expected := resources.MakeReceiveAdapter(&adapterArgs)

	ra, err := r.kubeClientSet.AppsV1().Deployments(src.Namespace).Get(ctx, expected.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		// Issue eventing#2842: Adater deployment name uses kmeta.ChildName. If a deployment by the previous name pattern is found, it should
		// be deleted. This might cause temporary downtime.
		if deprecatedName := utils.GenerateFixedName(adapterArgs.Source, fmt.Sprintf("prometheussource-%s", adapterArgs.Source.Name)); deprecatedName != expected.Name {
			if err := r.kubeClientSet.AppsV1().Deployments(src.Namespace).Delete(ctx, deprecatedName, metav1.DeleteOptions{}); err != nil && !apierrors.IsNotFound(err) {
				return nil, fmt.Errorf("error deleting deprecated named deployment: %v", err)
			}
			controller.GetEventRecorder(ctx).Eventf(src, corev1.EventTypeNormal, prometheussourceDeploymentDeleted, "Deprecated deployment removed: \"%s/%s\"", src.Namespace, deprecatedName)
		}
		ra, err = r.kubeClientSet.AppsV1().Deployments(src.Namespace).Create(ctx, expected, metav1.CreateOptions{})
		controller.GetEventRecorder(ctx).Eventf(src, corev1.EventTypeNormal, prometheussourceDeploymentCreated, "Deployment created, error: %v", err)
		return ra, err
	} else if err != nil {
		return nil, fmt.Errorf("error getting receive adapter: %v", err)
	} else if !metav1.IsControlledBy(ra, src) {
		return nil, fmt.Errorf("deployment %q is not owned by PrometheusSource %q", ra.Name, src.Name)
	} else if r.podSpecChanged(ra.Spec.Template.Spec, expected.Spec.Template.Spec) {
		ra.Spec.Template.Spec = expected.Spec.Template.Spec
		if ra, err = r.kubeClientSet.AppsV1().Deployments(src.Namespace).Update(ctx, ra, metav1.UpdateOptions{}); err != nil {
			return ra, err
		}
		controller.GetEventRecorder(ctx).Eventf(src, corev1.EventTypeNormal, prometheussourceDeploymentUpdated, "Deployment updated")
		return ra, nil
	} else {
		logging.FromContext(ctx).Debug("Reusing existing receive adapter", zap.Any("receiveAdapter", ra))
	}
	return ra, nil
}

func (r *Reconciler) podSpecChanged(oldPodSpec corev1.PodSpec, newPodSpec corev1.PodSpec) bool {
	if !equality.Semantic.DeepDerivative(newPodSpec, oldPodSpec) {
		return true
	}
	if len(oldPodSpec.Containers) != len(newPodSpec.Containers) {
		return true
	}
	for i := range newPodSpec.Containers {
		if !equality.Semantic.DeepEqual(newPodSpec.Containers[i].Env, oldPodSpec.Containers[i].Env) {
			return true
		}
	}
	return false
}

// makeEventSource computes the Cloud Event source attribute for the given source
func (r *Reconciler) makeEventSource(src *v1alpha1.PrometheusSource) string {
	return src.Namespace + "/" + src.Name
}
