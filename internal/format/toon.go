package format

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// TOONEncoder encodes Go values to TOON format
type TOONEncoder struct {
	indent       int
	indentString string
}

// NewTOONEncoder creates a new TOON encoder
func NewTOONEncoder() *TOONEncoder {
	return &TOONEncoder{
		indent:       0,
		indentString: "  ", // 2 spaces per level
	}
}

// Encode encodes a Go value to TOON format string
func (e *TOONEncoder) Encode(v interface{}) (string, error) {
	var sb strings.Builder
	if err := e.encodeValue(&sb, reflect.ValueOf(v), ""); err != nil {
		return "", err
	}
	return strings.TrimRight(sb.String(), "\n"), nil
}

// encodeValue encodes a reflect.Value to the string builder
func (e *TOONEncoder) encodeValue(sb *strings.Builder, v reflect.Value, key string) error {
	// Handle interface and pointer indirection
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		if v.IsNil() {
			e.writeIndent(sb)
			if key != "" {
				sb.WriteString(key)
				sb.WriteString(": ")
			}
			sb.WriteString("null\n")
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		return e.encodeString(sb, v.String(), key)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return e.encodeNumber(sb, fmt.Sprintf("%d", v.Int()), key)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return e.encodeNumber(sb, fmt.Sprintf("%d", v.Uint()), key)
	case reflect.Float32, reflect.Float64:
		return e.encodeNumber(sb, fmt.Sprintf("%g", v.Float()), key)
	case reflect.Bool:
		return e.encodeBool(sb, v.Bool(), key)
	case reflect.Slice, reflect.Array:
		return e.encodeArray(sb, v, key)
	case reflect.Map:
		return e.encodeMap(sb, v, key)
	case reflect.Struct:
		return e.encodeStruct(sb, v, key)
	default:
		return fmt.Errorf("unsupported type: %v", v.Kind())
	}
}

// encodeString encodes a string value
func (e *TOONEncoder) encodeString(sb *strings.Builder, s string, key string) error {
	e.writeIndent(sb)
	if key != "" {
		e.writeKey(sb, key)
		sb.WriteString(": ")
	}

	// Check if multiline (contains newlines)
	if strings.Contains(s, "\n") {
		sb.WriteString("|\n")
		e.indent++
		lines := strings.Split(s, "\n")
		for _, line := range lines {
			if line == "" {
				// Empty lines should not have indentation
				sb.WriteString("\n")
			} else {
				e.writeIndent(sb)
				sb.WriteString(line)
				sb.WriteString("\n")
			}
		}
		e.indent--
		return nil
	}

	// Quote if necessary
	if e.needsQuoting(s) {
		sb.WriteString(e.quoteString(s))
	} else {
		sb.WriteString(s)
	}
	sb.WriteString("\n")
	return nil
}

// needsQuoting determines if a string needs quoting
func (e *TOONEncoder) needsQuoting(s string) bool {
	if s == "" {
		return true
	}

	// Quote if looks like boolean
	if s == "true" || s == "false" || s == "null" {
		return true
	}

	// Quote if looks like number
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return true
	}

	// Quote if contains special characters
	specialChars := []string{":", ",", "\"", "\\", "\t", "|", "[", "]", "{", "}"}
	for _, char := range specialChars {
		if strings.Contains(s, char) {
			return true
		}
	}

	// Quote if starts with "- " (list item pattern)
	if strings.HasPrefix(s, "- ") {
		return true
	}

	// Quote if has leading/trailing spaces
	if strings.TrimSpace(s) != s {
		return true
	}

	return false
}

// quoteString adds quotes and escapes special characters
func (e *TOONEncoder) quoteString(s string) string {
	var sb strings.Builder
	sb.WriteString("\"")
	for _, r := range s {
		switch r {
		case '"':
			sb.WriteString("\\\"")
		case '\\':
			sb.WriteString("\\\\")
		case '\n':
			sb.WriteString("\\n")
		case '\r':
			sb.WriteString("\\r")
		case '\t':
			sb.WriteString("\\t")
		default:
			sb.WriteRune(r)
		}
	}
	sb.WriteString("\"")
	return sb.String()
}

// needsQuotingAsKey checks if a map key needs to be quoted in YAML
// File paths and keys with special characters need quoting
func (e *TOONEncoder) needsQuotingAsKey(s string) bool {
	if s == "" {
		return true
	}

	// Quote if contains path separators or file extensions (file paths)
	if strings.Contains(s, "/") || strings.Contains(s, ".") {
		return true
	}

	// Quote if contains hyphens (common in filenames like "my-file.go")
	if strings.Contains(s, "-") {
		return true
	}

	// Quote if contains YAML special characters
	specialChars := []string{":", ",", "\"", "\\", "\t", "|", "[", "]", "{", "}", "#", "&", "*", "?", ">", "<", "=", "!", "%", "@"}
	for _, char := range specialChars {
		if strings.Contains(s, char) {
			return true
		}
	}

	// Quote if starts with "- " (list item pattern)
	if strings.HasPrefix(s, "- ") {
		return true
	}

	// Quote if has leading/trailing spaces
	if strings.TrimSpace(s) != s {
		return true
	}

	// Quote if it looks like a number or boolean
	if s == "true" || s == "false" || s == "null" || s == "yes" || s == "no" {
		return true
	}
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return true
	}

	return false
}

// writeKey writes a map key, quoting it if necessary
func (e *TOONEncoder) writeKey(sb *strings.Builder, key string) {
	if e.needsQuotingAsKey(key) {
		sb.WriteString(e.quoteString(key))
	} else {
		sb.WriteString(key)
	}
}

// encodeNumber encodes a numeric value
func (e *TOONEncoder) encodeNumber(sb *strings.Builder, num string, key string) error {
	e.writeIndent(sb)
	if key != "" {
		e.writeKey(sb, key)
		sb.WriteString(": ")
	}
	sb.WriteString(num)
	sb.WriteString("\n")
	return nil
}

// encodeBool encodes a boolean value
func (e *TOONEncoder) encodeBool(sb *strings.Builder, b bool, key string) error {
	e.writeIndent(sb)
	if key != "" {
		e.writeKey(sb, key)
		sb.WriteString(": ")
	}
	if b {
		sb.WriteString("true")
	} else {
		sb.WriteString("false")
	}
	sb.WriteString("\n")
	return nil
}

// encodeArray encodes an array or slice
func (e *TOONEncoder) encodeArray(sb *strings.Builder, v reflect.Value, key string) error {
	length := v.Len()

	// Empty array
	if length == 0 {
		e.writeIndent(sb)
		if key != "" {
			e.writeKey(sb, key)
		}
		sb.WriteString("[0]:\n")
		return nil
	}

	// Check if all elements are primitives (strings, numbers, bools)
	allPrimitives := true
	for i := 0; i < length; i++ {
		elem := v.Index(i)
		for elem.Kind() == reflect.Interface || elem.Kind() == reflect.Ptr {
			if elem.IsNil() {
				break
			}
			elem = elem.Elem()
		}
		if !e.isPrimitive(elem.Kind()) {
			allPrimitives = false
			break
		}
	}

	// Inline format for primitive arrays
	if allPrimitives {
		e.writeIndent(sb)
		if key != "" {
			e.writeKey(sb, key)
		}
		sb.WriteString(fmt.Sprintf("[%d]: ", length))

		for i := 0; i < length; i++ {
			if i > 0 {
				sb.WriteString(",")
			}
			elem := v.Index(i)
			for elem.Kind() == reflect.Interface || elem.Kind() == reflect.Ptr {
				if elem.IsNil() {
					sb.WriteString("null")
					continue
				}
				elem = elem.Elem()
			}
			e.encodePrimitiveValue(sb, elem)
		}
		sb.WriteString("\n")
		return nil
	}

	// Check if all elements are structs/maps with uniform keys (for tabular format)
	if e.canUseTabular(v) {
		return e.encodeTabularArray(sb, v, key)
	}

	// List format for mixed or non-uniform arrays
	e.writeIndent(sb)
	if key != "" {
		e.writeKey(sb, key)
	}
	sb.WriteString(fmt.Sprintf("[%d]:\n", length))

	e.indent++
	for i := 0; i < length; i++ {
		elem := v.Index(i)
		e.writeIndent(sb)

		// For primitives, write inline
		for elem.Kind() == reflect.Interface || elem.Kind() == reflect.Ptr {
			if elem.IsNil() {
				sb.WriteString("- null\n")
				continue
			}
			elem = elem.Elem()
		}

		if e.isPrimitive(elem.Kind()) {
			sb.WriteString("- ")
			e.encodePrimitiveValue(sb, elem)
			sb.WriteString("\n")
		} else {
			// For objects, write dash without trailing space, then content on next line
			sb.WriteString("-\n")
			e.indent++
			if err := e.encodeValue(sb, elem, ""); err != nil {
				e.indent--
				return err
			}
			e.indent--
		}
	}
	e.indent--

	return nil
}

// canUseTabular checks if array can use tabular format
func (e *TOONEncoder) canUseTabular(v reflect.Value) bool {
	if v.Len() == 0 {
		return false
	}

	// Get first element keys
	var firstKeys []string
	firstElem := v.Index(0)
	for firstElem.Kind() == reflect.Interface || firstElem.Kind() == reflect.Ptr {
		if firstElem.IsNil() {
			return false
		}
		firstElem = firstElem.Elem()
	}

	switch firstElem.Kind() {
	case reflect.Map:
		for _, key := range firstElem.MapKeys() {
			firstKeys = append(firstKeys, key.String())
		}
	case reflect.Struct:
		t := firstElem.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath == "" { // Exported field
				firstKeys = append(firstKeys, field.Name)
			}
		}
	default:
		return false
	}

	sort.Strings(firstKeys)

	// Check all elements have same keys and primitive values
	for i := 0; i < v.Len(); i++ {
		elem := v.Index(i)
		for elem.Kind() == reflect.Interface || elem.Kind() == reflect.Ptr {
			if elem.IsNil() {
				return false
			}
			elem = elem.Elem()
		}

		var elemKeys []string
		switch elem.Kind() {
		case reflect.Map:
			for _, key := range elem.MapKeys() {
				elemKeys = append(elemKeys, key.String())
				// Check value is primitive
				val := elem.MapIndex(key)
				for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
					if val.IsNil() {
						break
					}
					val = val.Elem()
				}
				if !e.isPrimitive(val.Kind()) {
					return false
				}
			}
		case reflect.Struct:
			t := elem.Type()
			for j := 0; j < t.NumField(); j++ {
				field := t.Field(j)
				if field.PkgPath == "" { // Exported field
					elemKeys = append(elemKeys, field.Name)
					// Check value is primitive
					val := elem.Field(j)
					for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
						if val.IsNil() {
							break
						}
						val = val.Elem()
					}
					if !e.isPrimitive(val.Kind()) {
						return false
					}
				}
			}
		default:
			return false
		}

		sort.Strings(elemKeys)
		if !e.stringSlicesEqual(firstKeys, elemKeys) {
			return false
		}
	}

	return true
}

// encodeTabularArray encodes array in tabular format
func (e *TOONEncoder) encodeTabularArray(sb *strings.Builder, v reflect.Value, key string) error {
	length := v.Len()

	// Get keys from first element
	var keys []string
	firstElem := v.Index(0)
	for firstElem.Kind() == reflect.Interface || firstElem.Kind() == reflect.Ptr {
		firstElem = firstElem.Elem()
	}

	switch firstElem.Kind() {
	case reflect.Map:
		for _, k := range firstElem.MapKeys() {
			keys = append(keys, k.String())
		}
		sort.Strings(keys)
	case reflect.Struct:
		t := firstElem.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath == "" {
				keys = append(keys, field.Name)
			}
		}
	}

	// Write header
	e.writeIndent(sb)
	if key != "" {
		e.writeKey(sb, key)
	}
	sb.WriteString(fmt.Sprintf("[%d]{%s}:\n", length, strings.Join(keys, ",")))

	// Write rows
	e.indent++
	for i := 0; i < length; i++ {
		elem := v.Index(i)
		for elem.Kind() == reflect.Interface || elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		e.writeIndent(sb)
		for j, key := range keys {
			if j > 0 {
				sb.WriteString(",")
			}

			var val reflect.Value
			switch elem.Kind() {
			case reflect.Map:
				val = elem.MapIndex(reflect.ValueOf(key))
			case reflect.Struct:
				val = elem.FieldByName(key)
			}

			for val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
				if val.IsNil() {
					sb.WriteString("null")
					continue
				}
				val = val.Elem()
			}

			if val.IsValid() {
				e.encodePrimitiveValue(sb, val)
			}
		}
		sb.WriteString("\n")
	}
	e.indent--

	return nil
}

// encodeMap encodes a map as nested object
func (e *TOONEncoder) encodeMap(sb *strings.Builder, v reflect.Value, key string) error {
	if v.Len() == 0 {
		e.writeIndent(sb)
		if key != "" {
			e.writeKey(sb, key)
			sb.WriteString(":\n")
		}
		return nil
	}

	// Write key if present
	if key != "" {
		e.writeIndent(sb)
		e.writeKey(sb, key)
		sb.WriteString(":\n")
		e.indent++
	}

	// Sort keys for consistent output
	keys := v.MapKeys()
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})

	// Encode each key-value pair
	for _, k := range keys {
		val := v.MapIndex(k)
		if err := e.encodeValue(sb, val, k.String()); err != nil {
			return err
		}
	}

	if key != "" {
		e.indent--
	}

	return nil
}

// encodeStruct encodes a struct as nested object
func (e *TOONEncoder) encodeStruct(sb *strings.Builder, v reflect.Value, key string) error {
	t := v.Type()

	// Write key if present
	if key != "" {
		e.writeIndent(sb)
		e.writeKey(sb, key)
		sb.WriteString(":\n")
		e.indent++
	}

	// Encode each exported field
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" { // Skip unexported fields
			continue
		}

		fieldValue := v.Field(i)
		fieldName := field.Name

		// Use struct tag if present
		if tag := field.Tag.Get("toon"); tag != "" {
			if tag == "-" {
				continue
			}
			fieldName = tag
		}

		if err := e.encodeValue(sb, fieldValue, fieldName); err != nil {
			return err
		}
	}

	if key != "" {
		e.indent--
	}

	return nil
}

// Helper methods

func (e *TOONEncoder) writeIndent(sb *strings.Builder) {
	for i := 0; i < e.indent; i++ {
		sb.WriteString(e.indentString)
	}
}

func (e *TOONEncoder) isPrimitive(kind reflect.Kind) bool {
	switch kind {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool:
		return true
	}
	return false
}

func (e *TOONEncoder) encodePrimitiveValue(sb *strings.Builder, v reflect.Value) {
	switch v.Kind() {
	case reflect.String:
		s := v.String()
		if e.needsQuoting(s) {
			sb.WriteString(e.quoteString(s))
		} else {
			sb.WriteString(s)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		sb.WriteString(fmt.Sprintf("%d", v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		sb.WriteString(fmt.Sprintf("%d", v.Uint()))
	case reflect.Float32, reflect.Float64:
		sb.WriteString(fmt.Sprintf("%g", v.Float()))
	case reflect.Bool:
		if v.Bool() {
			sb.WriteString("true")
		} else {
			sb.WriteString("false")
		}
	default:
		sb.WriteString("null")
	}
}

func (e *TOONEncoder) stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
