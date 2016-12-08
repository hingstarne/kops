// Code generated by ""fitask" -type=LoadBalancer"; DO NOT EDIT

package awstasks

import (
	"encoding/json"

	"k8s.io/kops/upup/pkg/fi"
)

// LoadBalancer

// JSON marshalling boilerplate
type realLoadBalancer LoadBalancer

// UnmarshalJSON implements conversion to JSON, supporitng an alternate specification of the object as a string
func (o *LoadBalancer) UnmarshalJSON(data []byte) error {
	var jsonName string
	if err := json.Unmarshal(data, &jsonName); err == nil {
		o.Name = &jsonName
		return nil
	}

	var r realLoadBalancer
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	*o = LoadBalancer(r)
	return nil
}

var _ fi.HasName = &LoadBalancer{}

// GetName returns the Name of the object, implementing fi.HasName
func (o *LoadBalancer) GetName() *string {
	return o.Name
}

// SetName sets the Name of the object, implementing fi.SetName
func (o *LoadBalancer) SetName(name string) {
	o.Name = &name
}

// String is the stringer function for the task, producing readable output using fi.TaskAsString
func (o *LoadBalancer) String() string {
	return fi.TaskAsString(o)
}
