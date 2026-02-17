package guacamole

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// boolToString converts a boolean to the string Guacamole expects ("true" or "").
func boolToString(b bool) string {
	if b {
		return "true"
	}
	return ""
}

// stringToBool converts Guacamole's boolean strings ("true"/"") to a Go bool.
func stringToBool(v string) bool {
	if v == "" {
		return false
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false
	}
	return b
}

// sliceDiff returns elements in slice1 that are not in slice2.
// When bidirectional is true it also returns elements in slice2 not in slice1.
func sliceDiff(slice1 []string, slice2 []string, bidirectional bool) []string {
	var diff []string

	loopCount := 1
	if bidirectional {
		loopCount = 2
	}

	for i := 0; i < loopCount; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			if !found {
				diff = append(diff, s1)
			}
		}
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

// stringInSlice returns an error diagnostic if any test value is not in valid.
func stringInSlice(valid []string, test []string) diag.Diagnostics {
	var diags diag.Diagnostics

	for _, t := range test {
		found := false
		for _, v := range valid {
			if v == t {
				found = true
				break
			}
		}
		if !found {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Invalid value entered",
				Detail:   fmt.Sprintf("%s is not one of supported values: %s", t, strings.Join(valid, ", ")),
			})
		}
	}
	return diags
}

// checkForDuplicates returns an error diagnostic if any string appears more than once.
func checkForDuplicates(slice1 []string) diag.Diagnostics {
	var diags diag.Diagnostics
	seen := make(map[string]bool)
	var duplicates []string

	for _, v := range slice1 {
		if seen[v] {
			duplicates = append(duplicates, v)
		}
		seen[v] = true
	}
	if len(duplicates) > 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Duplicate entries found in array",
			Detail:   fmt.Sprintf("Found the duplicate entries: %s", strings.Join(duplicates, ", ")),
		})
	}
	return diags
}

// toHclString converts a value to an HCL-safe string representation.
// Used by acceptance tests to generate configuration snippets.
func toHclString(value interface{}, isNested bool) string {
	if slice, isSlice := tryToConvertToGenericSlice(value); isSlice {
		return sliceToHclString(slice)
	} else if m, isMap := tryToConvertToGenericMap(value); isMap {
		return mapToHclString(m)
	}
	return primitiveToHclString(value, isNested)
}

func tryToConvertToGenericSlice(value interface{}) ([]interface{}, bool) {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice {
		return nil, false
	}
	out := make([]interface{}, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		out[i] = rv.Index(i).Interface()
	}
	return out, true
}

func tryToConvertToGenericMap(value interface{}) (map[string]interface{}, bool) {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Map {
		return nil, false
	}
	if reflect.TypeOf(value).Key().Kind() != reflect.String {
		return nil, false
	}
	out := make(map[string]interface{}, rv.Len())
	for _, k := range rv.MapKeys() {
		out[k.String()] = rv.MapIndex(k).Interface()
	}
	return out, true
}

func sliceToHclString(slice []interface{}) string {
	parts := make([]string, len(slice))
	for i, v := range slice {
		parts[i] = toHclString(v, true)
	}
	return fmt.Sprintf("[%s]", strings.Join(parts, ", "))
}

func mapToHclString(m map[string]interface{}) string {
	var pairs []string
	for k, v := range m {
		if _, isMap := tryToConvertToGenericMap(v); isMap {
			pairs = append(pairs, fmt.Sprintf(`%s %s`, k, toHclString(v, true)))
		} else {
			pairs = append(pairs, fmt.Sprintf(`%s = %s`, k, toHclString(v, true)))
		}
	}
	return fmt.Sprintf("{\n%s\n}", strings.Join(pairs, "\n"))
}

func primitiveToHclString(value interface{}, isNested bool) string {
	if value == nil {
		return "null"
	}
	switch v := value.(type) {
	case bool:
		return strconv.FormatBool(v)
	case string:
		if isNested {
			return fmt.Sprintf("%q", v)
		}
		return fmt.Sprintf("%v", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// validateTimestring validates that a date string matches YYYY-MM-DD format.
func validateTimestring(timeString string, name string) diag.Diagnostics {
	var diags diag.Diagnostics
	matched, err := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, timeString)
	if err != nil || !matched {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  fmt.Sprintf("Invalid timestring format for: %s", name),
			Detail:   "Date string must be in the form of YYYY-MM-DD",
		})
	}
	return diags
}

// testAccCheckTestSliceVals is a test helper that verifies a TypeSet in
// Terraform state contains exactly the expected values.
func testAccCheckTestSliceVals(resourceName string, key string, expected []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		v, ok := rs.Primary.Attributes[fmt.Sprintf("%s.#", key)]
		if !ok {
			return fmt.Errorf("%s: attribute '%s.#' not found", resourceName, key)
		}
		testCount, _ := strconv.Atoi(v)
		if testCount == 0 {
			return fmt.Errorf("no entries found in state for key: %s", key)
		}

		var sv []string
		for i := 0; i < testCount; i++ {
			sv = append(sv, rs.Primary.Attributes[fmt.Sprintf("%s.%d", key, i)])
		}

		diff := sliceDiff(expected, sv, true)
		if len(diff) > 0 {
			return fmt.Errorf("set values %v do not match expected %v", sv, expected)
		}
		return nil
	}
}

// validSystemPermissions returns the list of valid Guacamole system permission names.
func validSystemPermissions() []string {
	return []string{
		"CREATE_USER",
		"CREATE_USER_GROUP",
		"CREATE_CONNECTION",
		"CREATE_CONNECTION_GROUP",
		"CREATE_SHARING_PROFILE",
		"ADMINISTER",
	}
}
