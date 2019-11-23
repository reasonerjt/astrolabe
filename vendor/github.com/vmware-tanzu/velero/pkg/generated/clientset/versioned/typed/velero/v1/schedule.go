/*
Copyright the Velero contributors.

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

package v1

import (
	"time"

	v1 "github.com/vmware-tanzu/velero/pkg/apis/velero/v1"
	scheme "github.com/vmware-tanzu/velero/pkg/generated/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SchedulesGetter has a method to return a ScheduleInterface.
// A group's client should implement this interface.
type SchedulesGetter interface {
	Schedules(namespace string) ScheduleInterface
}

// ScheduleInterface has methods to work with Schedule resources.
type ScheduleInterface interface {
	Create(*v1.Schedule) (*v1.Schedule, error)
	Update(*v1.Schedule) (*v1.Schedule, error)
	UpdateStatus(*v1.Schedule) (*v1.Schedule, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.Schedule, error)
	List(opts metav1.ListOptions) (*v1.ScheduleList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Schedule, err error)
	ScheduleExpansion
}

// schedules implements ScheduleInterface
type schedules struct {
	client rest.Interface
	ns     string
}

// newSchedules returns a Schedules
func newSchedules(c *VeleroV1Client, namespace string) *schedules {
	return &schedules{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the schedule, and returns the corresponding schedule object, and an error if there is any.
func (c *schedules) Get(name string, options metav1.GetOptions) (result *v1.Schedule, err error) {
	result = &v1.Schedule{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("schedules").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Schedules that match those selectors.
func (c *schedules) List(opts metav1.ListOptions) (result *v1.ScheduleList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ScheduleList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("schedules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested schedules.
func (c *schedules) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("schedules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a schedule and creates it.  Returns the server's representation of the schedule, and an error, if there is any.
func (c *schedules) Create(schedule *v1.Schedule) (result *v1.Schedule, err error) {
	result = &v1.Schedule{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("schedules").
		Body(schedule).
		Do().
		Into(result)
	return
}

// Update takes the representation of a schedule and updates it. Returns the server's representation of the schedule, and an error, if there is any.
func (c *schedules) Update(schedule *v1.Schedule) (result *v1.Schedule, err error) {
	result = &v1.Schedule{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("schedules").
		Name(schedule.Name).
		Body(schedule).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *schedules) UpdateStatus(schedule *v1.Schedule) (result *v1.Schedule, err error) {
	result = &v1.Schedule{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("schedules").
		Name(schedule.Name).
		SubResource("status").
		Body(schedule).
		Do().
		Into(result)
	return
}

// Delete takes name of the schedule and deletes it. Returns an error if one occurs.
func (c *schedules) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("schedules").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *schedules) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("schedules").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched schedule.
func (c *schedules) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Schedule, err error) {
	result = &v1.Schedule{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("schedules").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
