package guacamole

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestConnectionAttributesSchemaElem(t *testing.T) {
	elem := connectionAttributesSchemaElem()

	expectedKeys := []string{
		"guacd_hostname",
		"guacd_port",
		"guacd_encryption",
		"failover_only",
		"weight",
		"max_connections",
		"max_connections_per_user",
	}
	for _, key := range expectedKeys {
		if _, ok := elem.Schema[key]; !ok {
			t.Errorf("connectionAttributesSchemaElem: missing schema key %q", key)
		}
	}
	if len(elem.Schema) != len(expectedKeys) {
		t.Errorf("connectionAttributesSchemaElem: expected %d schema keys, got %d", len(expectedKeys), len(elem.Schema))
	}
}

func TestConnectionAttributesSchemaElemTypes(t *testing.T) {
	elem := connectionAttributesSchemaElem()

	stringFields := []string{"guacd_hostname", "guacd_port", "guacd_encryption", "weight", "max_connections", "max_connections_per_user"}
	for _, key := range stringFields {
		s, ok := elem.Schema[key]
		if !ok {
			t.Errorf("missing key %q", key)
			continue
		}
		if s.Type != schema.TypeString {
			t.Errorf("key %q: expected TypeString, got %v", key, s.Type)
		}
	}

	boolFields := []string{"failover_only"}
	for _, key := range boolFields {
		s, ok := elem.Schema[key]
		if !ok {
			t.Errorf("missing key %q", key)
			continue
		}
		if s.Type != schema.TypeBool {
			t.Errorf("key %q: expected TypeBool, got %v", key, s.Type)
		}
	}
}

func TestReadConnectionAttributes(t *testing.T) {
	attrs := map[string]string{
		"guacd-hostname":           "proxy.example.com",
		"guacd-port":               "4822",
		"guacd-encryption":         "ssl",
		"failover-only":            "true",
		"weight":                   "10",
		"max-connections":          "4",
		"max-connections-per-user": "2",
	}

	result := readConnectionAttributes(attrs)

	if len(result) != 1 {
		t.Fatalf("expected 1 element, got %d", len(result))
	}
	m, ok := result[0].(map[string]interface{})
	if !ok {
		t.Fatal("result[0] is not map[string]interface{}")
	}

	cases := []struct {
		key      string
		expected interface{}
	}{
		{"guacd_hostname", "proxy.example.com"},
		{"guacd_port", "4822"},
		{"guacd_encryption", "ssl"},
		{"failover_only", true},
		{"weight", "10"},
		{"max_connections", "4"},
		{"max_connections_per_user", "2"},
	}
	for _, tc := range cases {
		if m[tc.key] != tc.expected {
			t.Errorf("readConnectionAttributes[%q]: expected %v, got %v", tc.key, tc.expected, m[tc.key])
		}
	}
}

func TestReadConnectionAttributesEmpty(t *testing.T) {
	result := readConnectionAttributes(map[string]string{})

	if len(result) != 1 {
		t.Fatalf("expected 1 element, got %d", len(result))
	}
	m, ok := result[0].(map[string]interface{})
	if !ok {
		t.Fatal("result[0] is not map[string]interface{}")
	}

	if m["guacd_hostname"] != "" {
		t.Errorf("guacd_hostname: expected empty string, got %q", m["guacd_hostname"])
	}
	if m["failover_only"] != false {
		t.Errorf("failover_only: expected false, got %v", m["failover_only"])
	}
}

func TestReadConnectionAttributesFalseBoolean(t *testing.T) {
	attrs := map[string]string{
		"failover-only": "",
	}
	result := readConnectionAttributes(attrs)
	m := result[0].(map[string]interface{})
	if m["failover_only"] != false {
		t.Errorf("failover_only with empty string: expected false, got %v", m["failover_only"])
	}
}

func TestReadConnectionAttributesTrueBooleanVariants(t *testing.T) {
	cases := []struct {
		apiValue string
		expected bool
	}{
		{"true", true},
		{"", false},
	}
	for _, tc := range cases {
		attrs := map[string]string{"failover-only": tc.apiValue}
		result := readConnectionAttributes(attrs)
		m := result[0].(map[string]interface{})
		got, ok := m["failover_only"].(bool)
		if !ok {
			t.Errorf("failover_only with %q: result is not bool (%T)", tc.apiValue, m["failover_only"])
			continue
		}
		if got != tc.expected {
			t.Errorf("failover_only with %q: expected %v, got %v", tc.apiValue, tc.expected, got)
		}
	}
}

func TestBuildConnectionAttributes(t *testing.T) {
	attrList := []interface{}{
		map[string]interface{}{
			"guacd_hostname":           "proxy.example.com",
			"guacd_port":               "4822",
			"guacd_encryption":         "ssl",
			"failover_only":            true,
			"weight":                   "10",
			"max_connections":          "4",
			"max_connections_per_user": "2",
		},
	}

	result := buildConnectionAttributes(attrList)

	cases := []struct {
		key      string
		expected string
	}{
		{"guacd-hostname", "proxy.example.com"},
		{"guacd-port", "4822"},
		{"guacd-encryption", "ssl"},
		{"failover-only", "true"},
		{"weight", "10"},
		{"max-connections", "4"},
		{"max-connections-per-user", "2"},
	}
	for _, tc := range cases {
		if result[tc.key] != tc.expected {
			t.Errorf("buildConnectionAttributes[%q]: expected %q, got %q", tc.key, tc.expected, result[tc.key])
		}
	}
}

func TestBuildConnectionAttributesFalseBoolean(t *testing.T) {
	attrList := []interface{}{
		map[string]interface{}{
			"guacd_hostname":           "",
			"guacd_port":               "",
			"guacd_encryption":         "",
			"failover_only":            false,
			"weight":                   "",
			"max_connections":          "",
			"max_connections_per_user": "",
		},
	}

	result := buildConnectionAttributes(attrList)
	if result["failover-only"] != "" {
		t.Errorf("failover-only with false: expected empty string, got %q", result["failover-only"])
	}
}

func TestBuildConnectionAttributesEmpty(t *testing.T) {
	result := buildConnectionAttributes([]interface{}{})
	if result != nil {
		t.Errorf("expected nil for empty attrList, got %v", result)
	}
}

func TestBuildConnectionAttributesRoundtrip(t *testing.T) {
	original := map[string]string{
		"guacd-hostname":           "proxy.example.com",
		"guacd-port":               "4822",
		"guacd-encryption":         "ssl",
		"failover-only":            "true",
		"weight":                   "5",
		"max-connections":          "10",
		"max-connections-per-user": "3",
	}

	// Read (API → TF)
	tfList := readConnectionAttributes(original)

	// Build (TF → API)
	rebuilt := buildConnectionAttributes(tfList)

	for k, v := range original {
		if rebuilt[k] != v {
			t.Errorf("roundtrip[%q]: expected %q, got %q", k, v, rebuilt[k])
		}
	}
}
