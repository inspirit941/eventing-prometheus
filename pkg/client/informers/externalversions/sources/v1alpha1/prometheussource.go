/*
Copyright 2020 The Knative Authors

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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	sourcesv1alpha1 "github.com/inspirit941/eventing-prometheus/pkg/apis/sources/v1alpha1"
	versioned "github.com/inspirit941/eventing-prometheus/pkg/client/clientset/versioned"
	internalinterfaces "github.com/inspirit941/eventing-prometheus/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/inspirit941/eventing-prometheus/pkg/client/listers/sources/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// PrometheusSourceInformer provides access to a shared informer and lister for
// PrometheusSources.
type PrometheusSourceInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.PrometheusSourceLister
}

type prometheusSourceInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewPrometheusSourceInformer constructs a new informer for PrometheusSource type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewPrometheusSourceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredPrometheusSourceInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredPrometheusSourceInformer constructs a new informer for PrometheusSource type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredPrometheusSourceInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SourcesV1alpha1().PrometheusSources(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.SourcesV1alpha1().PrometheusSources(namespace).Watch(context.TODO(), options)
			},
		},
		&sourcesv1alpha1.PrometheusSource{},
		resyncPeriod,
		indexers,
	)
}

func (f *prometheusSourceInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredPrometheusSourceInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *prometheusSourceInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&sourcesv1alpha1.PrometheusSource{}, f.defaultInformer)
}

func (f *prometheusSourceInformer) Lister() v1alpha1.PrometheusSourceLister {
	return v1alpha1.NewPrometheusSourceLister(f.Informer().GetIndexer())
}
