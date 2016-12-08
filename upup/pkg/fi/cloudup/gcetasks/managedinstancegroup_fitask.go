// Code generated by ""fitask" -type=ManagedInstanceGroup"; DO NOT EDIT

package gcetasks

import (
	"encoding/json"

	"k8s.io/kops/upup/pkg/fi"
)

// ManagedInstanceGroup

// JSON marshalling boilerplate
type realManagedInstanceGroup ManagedInstanceGroup

// UnmarshalJSON implements conversion to JSON, supporitng an alternate specification of the object as a string
func (o *ManagedInstanceGroup) UnmarshalJSON(data []byte) error {
	var jsonName string
	if err := json.Unmarshal(data, &jsonName); err == nil {
		o.Name = &jsonName
		return nil
	}

	var r realManagedInstanceGroup
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	*o = ManagedInstanceGroup(r)
	return nil
}

var _ fi.HasName = &ManagedInstanceGroup{}

// GetName returns the Name of the object, implementing fi.HasName
func (o *ManagedInstanceGroup) GetName() *string {
	return o.Name
}

// SetName sets the Name of the object, implementing fi.SetName
func (o *ManagedInstanceGroup) SetName(name string) {
	o.Name = &name
}

// String is the stringer function for the task, producing readable output using fi.TaskAsString
func (o *ManagedInstanceGroup) String() string {
	return fi.TaskAsString(o)
}
