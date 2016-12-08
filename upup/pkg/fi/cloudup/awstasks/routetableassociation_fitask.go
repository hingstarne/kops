// Code generated by ""fitask" -type=RouteTableAssociation"; DO NOT EDIT

package awstasks

import (
	"encoding/json"

	"k8s.io/kops/upup/pkg/fi"
)

// RouteTableAssociation

// JSON marshalling boilerplate
type realRouteTableAssociation RouteTableAssociation

// UnmarshalJSON implements conversion to JSON, supporitng an alternate specification of the object as a string
func (o *RouteTableAssociation) UnmarshalJSON(data []byte) error {
	var jsonName string
	if err := json.Unmarshal(data, &jsonName); err == nil {
		o.Name = &jsonName
		return nil
	}

	var r realRouteTableAssociation
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	*o = RouteTableAssociation(r)
	return nil
}

var _ fi.HasName = &RouteTableAssociation{}

// GetName returns the Name of the object, implementing fi.HasName
func (o *RouteTableAssociation) GetName() *string {
	return o.Name
}

// SetName sets the Name of the object, implementing fi.SetName
func (o *RouteTableAssociation) SetName(name string) {
	o.Name = &name
}

// String is the stringer function for the task, producing readable output using fi.TaskAsString
func (o *RouteTableAssociation) String() string {
	return fi.TaskAsString(o)
}
