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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/inspirit941/eventing-prometheus/pkg/apis/sources/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePrometheusSources implements PrometheusSourceInterface
type FakePrometheusSources struct {
	Fake *FakeSourcesV1alpha1
	ns   string
}

var prometheussourcesResource = schema.GroupVersionResource{Group: "sources.knative.dev", Version: "v1alpha1", Resource: "prometheussources"}

var prometheussourcesKind = schema.GroupVersionKind{Group: "sources.knative.dev", Version: "v1alpha1", Kind: "PrometheusSource"}

// Get takes name of the prometheusSource, and returns the corresponding prometheusSource object, and an error if there is any.
func (c *FakePrometheusSources) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.PrometheusSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(prometheussourcesResource, c.ns, name), &v1alpha1.PrometheusSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PrometheusSource), err
}

// List takes label and field selectors, and returns the list of PrometheusSources that match those selectors.
func (c *FakePrometheusSources) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.PrometheusSourceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(prometheussourcesResource, prometheussourcesKind, c.ns, opts), &v1alpha1.PrometheusSourceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.PrometheusSourceList{ListMeta: obj.(*v1alpha1.PrometheusSourceList).ListMeta}
	for _, item := range obj.(*v1alpha1.PrometheusSourceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested prometheusSources.
func (c *FakePrometheusSources) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(prometheussourcesResource, c.ns, opts))

}

// Create takes the representation of a prometheusSource and creates it.  Returns the server's representation of the prometheusSource, and an error, if there is any.
func (c *FakePrometheusSources) Create(ctx context.Context, prometheusSource *v1alpha1.PrometheusSource, opts v1.CreateOptions) (result *v1alpha1.PrometheusSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(prometheussourcesResource, c.ns, prometheusSource), &v1alpha1.PrometheusSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PrometheusSource), err
}

// Update takes the representation of a prometheusSource and updates it. Returns the server's representation of the prometheusSource, and an error, if there is any.
func (c *FakePrometheusSources) Update(ctx context.Context, prometheusSource *v1alpha1.PrometheusSource, opts v1.UpdateOptions) (result *v1alpha1.PrometheusSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(prometheussourcesResource, c.ns, prometheusSource), &v1alpha1.PrometheusSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PrometheusSource), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakePrometheusSources) UpdateStatus(ctx context.Context, prometheusSource *v1alpha1.PrometheusSource, opts v1.UpdateOptions) (*v1alpha1.PrometheusSource, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(prometheussourcesResource, "status", c.ns, prometheusSource), &v1alpha1.PrometheusSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PrometheusSource), err
}

// Delete takes name of the prometheusSource and deletes it. Returns an error if one occurs.
func (c *FakePrometheusSources) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(prometheussourcesResource, c.ns, name), &v1alpha1.PrometheusSource{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePrometheusSources) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(prometheussourcesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.PrometheusSourceList{})
	return err
}

// Patch applies the patch and returns the patched prometheusSource.
func (c *FakePrometheusSources) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.PrometheusSource, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(prometheussourcesResource, c.ns, name, pt, data, subresources...), &v1alpha1.PrometheusSource{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.PrometheusSource), err
}
