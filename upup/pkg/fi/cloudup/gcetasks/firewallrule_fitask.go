// Code generated by ""fitask" -type=FirewallRule"; DO NOT EDIT

package gcetasks

import (
	"encoding/json"

	"k8s.io/kops/upup/pkg/fi"
)

// FirewallRule

// JSON marshalling boilerplate
type realFirewallRule FirewallRule

// UnmarshalJSON implements conversion to JSON, supporitng an alternate specification of the object as a string
func (o *FirewallRule) UnmarshalJSON(data []byte) error {
	var jsonName string
	if err := json.Unmarshal(data, &jsonName); err == nil {
		o.Name = &jsonName
		return nil
	}

	var r realFirewallRule
	if err := json.Unmarshal(data, &r); err != nil {
		return err
	}
	*o = FirewallRule(r)
	return nil
}

var _ fi.HasName = &FirewallRule{}

// GetName returns the Name of the object, implementing fi.HasName
func (o *FirewallRule) GetName() *string {
	return o.Name
}

// SetName sets the Name of the object, implementing fi.SetName
func (o *FirewallRule) SetName(name string) {
	o.Name = &name
}

// String is the stringer function for the task, producing readable output using fi.TaskAsString
func (o *FirewallRule) String() string {
	return fi.TaskAsString(o)
}
