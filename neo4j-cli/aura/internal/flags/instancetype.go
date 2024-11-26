package flags

import "errors"

type InstanceType string

// String is used both by fmt.Print and by Cobra in help text
func (e *InstanceType) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *InstanceType) Set(v string) error {
	switch v {
	case "free-db", "professional-db", "business-critical", "enterprise-db", "professional-ds", "enterprise-ds":
		*e = InstanceType(v)
		return nil
	default:
		return errors.New(`must be one of "free-db", "professional-db", "business-critical", "enterprise-db", "professional-ds", or "enterprise-ds"`)
	}
}

// Type is only used in help text
func (e *InstanceType) Type() string {
	return "type"
}
