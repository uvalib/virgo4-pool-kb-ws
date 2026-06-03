package bedrock

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

var (
	valueFieldRE    = regexp.MustCompile(`Value:\s*"?([^"}]*)"?`)
	wrappedStructRE = regexp.MustCompile(`^&?\{([^}]*)\}$`)
)

func extractMetadataString(metadata interface{}, keys ...string) string {
	mv := reflect.ValueOf(metadata)
	if mv.Kind() != reflect.Map {
		return ""
	}
	for _, key := range keys {
		kv := reflect.ValueOf(key)
		val := mv.MapIndex(kv)
		if val.IsValid() {
			return metadataValueToString(val.Interface())
		}
		for _, k := range mv.MapKeys() {
			if strings.EqualFold(fmt.Sprintf("%v", k.Interface()), key) {
				return metadataValueToString(mv.MapIndex(k).Interface())
			}
		}
	}
	return ""
}

func metadataValueToString(val interface{}) string {
	if val == nil {
		return ""
	}

	switch v := val.(type) {
	case string:
		return cleanMetadataString(v)
	case *string:
		if v != nil {
			return cleanMetadataString(*v)
		}
		return ""
	}

	var decoded string
	if err := unmarshalSmithyDocument(val, &decoded); err == nil && decoded != "" {
		return cleanMetadataString(decoded)
	}

	if extracted := extractValueField(val); extracted != "" {
		return cleanMetadataString(extracted)
	}

	return cleanMetadataString(fmt.Sprintf("%v", val))
}

func extractValueField(val interface{}) string {
	rv := reflect.ValueOf(val)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return ""
		}
		rv = rv.Elem()
	}

	switch rv.Kind() {
	case reflect.String:
		return rv.String()
	case reflect.Struct:
		if fv := rv.FieldByName("Value"); fv.IsValid() {
			switch fv.Kind() {
			case reflect.String:
				return fv.String()
			case reflect.Pointer:
				if !fv.IsNil() && fv.Elem().Kind() == reflect.String {
					return fv.Elem().String()
				}
			}
		}
	}
	return ""
}

func unmarshalSmithyDocument(v interface{}, target interface{}) error {
	if v == nil {
		return fmt.Errorf("nil document")
	}
	rv := reflect.ValueOf(v)
	method := rv.MethodByName("UnmarshalSmithyDocument")
	if method.IsValid() {
		results := method.Call([]reflect.Value{reflect.ValueOf(target)})
		if len(results) > 0 && !results[0].IsNil() {
			return results[0].Interface().(error)
		}
		return nil
	}
	genericMethod := rv.MethodByName("Unmarshal")
	if genericMethod.IsValid() {
		results := genericMethod.Call([]reflect.Value{reflect.ValueOf(target)})
		if len(results) > 0 && !results[0].IsNil() {
			return results[0].Interface().(error)
		}
		return nil
	}
	return fmt.Errorf("no unmarshal method on %T", v)
}

func cleanMetadataString(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	if matches := wrappedStructRE.FindStringSubmatch(s); len(matches) > 1 {
		s = strings.TrimSpace(matches[1])
	}

	if strings.Contains(s, "Value:") {
		if matches := valueFieldRE.FindStringSubmatch(s); len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	return s
}
