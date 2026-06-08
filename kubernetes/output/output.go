package output

import (
	"encoding/json"

	"sigs.k8s.io/yaml"
)

// Format returns formatted output for a K8s object or slice.
// format can be "json", "yaml", or "" (returns empty string, caller handles summary).
func Format(format string, raw interface{}) (string, error) {
	switch format {
	case "json":
		b, err := json.MarshalIndent(raw, "", "  ")
		if err != nil {
			return "", err
		}
		return string(b), nil
	case "yaml":
		b, err := yaml.Marshal(raw)
		if err != nil {
			return "", err
		}
		return string(b), nil
	default:
		return "", nil
	}
}
