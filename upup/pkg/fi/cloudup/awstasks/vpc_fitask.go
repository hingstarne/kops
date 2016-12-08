// Code generated by ""fitask" -type=VPC"; DO NOT EDIT

package awstasks

import (
	"encoding/json"

	"k8s.io/kops/upup/pkg/fi"
)

// VPC

// JSON marshalling boilerplate
type realVPC VPC

// UnmarshalJSON implements conversion to JSON, supporitng an alternate specification of the object as a string
func (o *VPC) UnmarshalJSON(data []byte) error {
	var jsonName string
	if err := json.Unmarshal(data, &jsonName); err == nil {
		o.Name = &jsonName
		return nil
	}

	var r realVPC
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	*o = VPC(r)
	return nil
}

var _ fi.HasName = &VPC{}

// GetName returns the Name of the object, implementing fi.HasName
func (o *VPC) GetName() *string {
	return o.Name
}

// SetName sets the Name of the object, implementing fi.SetName
func (o *VPC) SetName(name string) {
	o.Name = &name
}

// String is the stringer function for the task, producing readable output using fi.TaskAsString
func (o *VPC) String() string {
	return fi.TaskAsString(o)
}
