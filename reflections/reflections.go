// Package reflections collects common functions that use reflection.
package reflections

import (
	"fmt"
	"log/slog"
	"reflect"
	"regexp"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"unsafe"
)

// AnyType gets the type of an object.
//   - obj the object to get the type for.
//
// Return the type of obj or nil (if obj is nil).
func AnyType(obj any) reflect.Type {
	t := reflect.TypeOf(obj)
	if IsPointer(obj) {
		t = t.Elem()

	}

	return t

}

// CopyStruct copies field values from a source struct to a destination struct.
//   - src is the struct to copy field values from
//   - dst is the struct to copy field values to
//   - isZeroOpt optionally set to true to only copy values to fields with a zero value.
func CopyStruct(src any, dst any, isZeroOpt ...bool) {
	srcVal := reflect.ValueOf(src)
	if srcVal.Kind() == reflect.Pointer {
		srcVal = srcVal.Elem()

	}

	dstVal := reflect.ValueOf(dst).Elem()
	for i := 0; i < srcVal.NumField(); i++ {
		srcField := srcVal.Type().Field(i)
		if dstField := dstVal.FieldByName(srcField.Name); dstField.IsValid() {
			if len(isZeroOpt) > 0 && isZeroOpt[0] {
				if dstField.IsZero() {
					dstField.Set(srcVal.Field(i))

				}

			} else {
				dstField.Set(srcVal.Field(i))

			}

		}

	}

}

// GetFieldAndTagNames get the set of field names for obj along with a map of obj tags valued by the field name.
//   - obj is the object to get field and tag names from
//   - tag is the object tag to get from obj (e.g. "json")
//
// Returns
//   - A set of obj field names
//   - A map keyed by obj tags valued by their corresponding field name
func FieldAndTagNames(obj any, tag string) (map[string]bool, map[string]string) {
	if obj == nil || !IsStruct(obj) {
		return map[string]bool{}, map[string]string{}

	}

	val := reflect.Indirect(reflect.ValueOf(obj))
	t := val.Type()
	names, tagToName := map[string]bool{}, map[string]string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		names[field.Name] = true

		tagName := field.Tag.Get(tag)
		if tagName != "" && tagName != "-" {
			tagToName[tagName] = field.Name

		}

	}

	return names, tagToName

}

// FuncSig gets the calling functions name with an added suffix.
//   - skip is the number of stack frames to ascend minus the current one
//   - suffix is the string to add to the end of the method name; "() - " is the default)
//
// Return the calling functions name with suffix (e.g. "main() - ").
func FuncSig(skip int, suffixOpt ...string) string {
	suffix := "() - "
	if len(suffixOpt) > 0 {
		suffix = suffixOpt[0]

	}

	funcSig := FunctionName(skip) + suffix

	return funcSig

}

// FunctionName returns the name of the calling function.
//   - skip is the number of stack frames to ascend minus the current one (i.e. the current stack frame is always skipped)
func FunctionName(skip int) string {
	// Always skip this function name (i.e. GetFunctionName)
	pc, _, _, ok := runtime.Caller(1 + skip)
	if !ok {
		return ""

	}

	return runtime.FuncForPC(pc).Name()

}

// GetType returns the object's type.
//   - obj is the object to get the type of
func GetType(obj any) reflect.Kind {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()

	}

	return v.Kind()

}

// IsPointer determines if the object is a pointer.
//   - obj the type to check
//
// Return true if obj is a pointer else false.
func IsPointer(obj any) bool {
	return reflect.ValueOf(obj).Kind() == reflect.Pointer

}

// IsStruct determines if the object is a struct.
//   - obj is the object to check for a struct
//
// Return true if obj is a struct else false.
func IsStruct(obj any) bool {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Pointer {
		return v.Elem().Kind() == reflect.Struct

	}

	return reflect.ValueOf(obj).Kind() == reflect.Struct

}

// EmptyStructInstance creates a new instance of T.
func EmptyStructInstance[T any]() any {
	var tAny T
	t := reflect.TypeOf(tAny)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()

	}

	if t.Kind() != reflect.Struct {
		return nil

	}

	return reflect.New(t).Interface()

}

// InitializeStruct initializes struct fields that need to be initialized (i.e. maps).
//   - sPtr is a pointer to a struct
//
// Return an error if sPtr is not a pointer to a struct or is nil.
func InitializeStruct(sPtr any) error {
	v := reflect.ValueOf(sPtr)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct || v.IsNil() {
		return fmt.Errorf("%s must be a pointer to a struct, kind is %s", reflect.TypeOf(sPtr).Name(), v.Kind())

	}

	val := v.Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Field(i)
		fieldType := typ.Field(i)

		if fieldVal.Kind() == reflect.Map && fieldVal.IsNil() {
			newMap := reflect.MakeMap(fieldType.Type)
			if !fieldVal.CanSet() {
				ptr := unsafe.Pointer(fieldVal.UnsafeAddr())
				reflect.NewAt(fieldType.Type, ptr).Elem().Set(newMap)

			} else {
				fieldVal.Set(newMap)

			}

		}

	}

	return nil

}

// NewMapInstance creates a new map instance of [K]V.
func NewMapInstance[K comparable, V any]() map[K]V {
	return make(map[K]V)

}

// ReflectFieldByIndex gets a struct field's reflect.Value.
//   - obj is the struct to get the reflect.Value from
//   - fieldIndex is the index/number of the field to get the reflect.Value for
//
// Return the struct field's reflect.Value or an error if obj is not a struct.
func ReflectFieldByIndex(obj any, fieldIndex int) (reflect.Value, error) {
	if !IsStruct(obj) {
		return reflect.ValueOf(nil), fmt.Errorf("obj is not a struct")

	}

	v := reflect.ValueOf(obj).Elem()
	var field reflect.Value
	if fieldIndex >= 0 && fieldIndex < v.NumField() {
		field = v.Field(fieldIndex)

	} else {
		return reflect.ValueOf(nil), fmt.Errorf("fieldIndex %d must be a valid field number in obj", fieldIndex)

	}

	return field, nil

}

// ReflectFieldByName gets a struct field's reflect.Value.
//   - obj is the struct to get the reflect.Value from
//   - fieldName is the name of the field in the struct to get the reflect.Value for
//
// Return the struct field's reflect.Value or an error if fieldName does not exist in obj.
func ReflectFieldByName(obj any, fieldName string) (reflect.Value, error) {
	if !IsStruct(obj) {
		return reflect.ValueOf(nil), fmt.Errorf("obj is not a struct")

	}

	var v reflect.Value = reflect.ValueOf(obj)
	if v.Kind() == reflect.Pointer {
		v = reflect.ValueOf(obj).Elem()

	}

	var field reflect.Value
	field = v.FieldByName(fieldName)
	if !field.IsValid() {
		return reflect.ValueOf(nil), fmt.Errorf("fieldName %s is not a valid field in obj", fieldName)

	}

	return field, nil

}

// SetStructField sets the value of a struct's field.
//   - obj contains the field to set the value for
//   - fieldName is the name of the field in obj to set the value for
//   - fieldValue is the value to set for fieldName
//
// Return error if
//   - field does not exist
//   - value type isn't supported
func SetStructField(obj any, fieldName string, fieldValue string) error {
	if !IsStruct(obj) {
		return fmt.Errorf("obj is not a struct")

	}

	field, err := ReflectFieldByName(obj, fieldName)
	if err != nil {
		return err

	}

	return SetStructFieldByType(field, fieldValue)

}

// SetStructFieldByType dynamically sets a struct field to a string value converted to the correct type
// (String, Int, Int64, Bool, Float32).
//
//   - field is the struct field to set.
//   - stringValue is the value to convert to the field's type.
//
// Return an error if the stringValue is not convertible to the struct's field type.
func SetStructFieldByType(field reflect.Value, stringValue string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(stringValue)

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bitSize := 64
		switch field.Kind() {
		case reflect.Uint8:
			bitSize = 8

		case reflect.Uint16:
			bitSize = 16

		case reflect.Uint32:
			bitSize = 32

		}

		intVal, err := strconv.ParseUint(stringValue, 10, bitSize)
		if err != nil {
			return err

		}
		field.SetUint(intVal)

	case reflect.Int, reflect.Int64:
		intVal, err := strconv.ParseInt(stringValue, 10, 64)
		if err != nil {
			return err

		}

		field.SetInt(intVal)

	case reflect.Bool:
		boolVal, err := strconv.ParseBool(stringValue)
		if err != nil {
			return err

		}

		field.SetBool(boolVal)

	case reflect.Float32:
		floatVal, err := strconv.ParseFloat(stringValue, 64)
		if err != nil {
			return err

		}

		field.SetFloat(floatVal)

	}

	return nil

}

// Need to fix this to print any e.g. DefaultConfigger.DefaultParentValues   any    `json:"-"` // Default values for the parent
//
// StructString gets all struct fields and values (struct_name {field_name:[(field_type) ]field_value}).
//   - anyStruct is the struct to print
//   - stringOpts "type" includes the struct field's type, "newline" tries to put fields on separate lines, add tabs i.e. "\t\t" to prefix lines
//
// Return a string containing a string representation of the struct.
func StructString(anyStruct any, stringOpts ...string) string {
	objValue := reflect.ValueOf(anyStruct)
	if objValue.Kind() == reflect.Pointer {
		objValue = objValue.Elem()

	}

	if objValue.Kind() == reflect.Struct && !objValue.CanAddr() {
		newObjValue := reflect.New(objValue.Type())
		newObjValue.Elem().Set(objValue)
		objValue = newObjValue.Elem()

	}

	newline, tab, endTab := "", "", ""
	if len(stringOpts) > 0 && slices.Contains(stringOpts, "newline") {
		newline = "\n"
		if tab == "" {
			tab = "\t"

		}

	}

	// use tabs if also using newline
	if len(stringOpts) > 0 && newline != "" {
		tabIndex := slices.IndexFunc(stringOpts, func(s string) bool {
			return strings.Contains(s, "\t")

		})
		if tabIndex > -1 && newline != "" {
			tab = stringOpts[tabIndex]
			tabNum := strings.Count(tab, "\t")
			endTab = strings.Repeat("\t", tabNum-1)

		}

	}

	var sb strings.Builder
	structStringHelper(&sb, objValue, reflect.TypeOf(objValue).Name(), newline, tab, endTab, stringOpts...)
	return sb.String()

}

func structStringHelper(sb *strings.Builder, objValue reflect.Value, objName, newline, tab, endTab string, stringOpts ...string) {
	/*if objValue.Kind() == reflect.Pointer {
		objValue = objValue.Elem()

	}*/
	objType := objValue.Type()

	colon := ":"
	if objName == "" {
		colon = ""

	}

	checkTab := ""
	if matched, _ := regexp.MatchString("\n$", sb.String()); matched {
		checkTab = tab

	}

	if objValue.Kind() == reflect.Struct {
		for i := 0; i < objValue.NumField(); i++ {
			fieldValue := objValue.Field(i)
			field := objType.Field(i)
			fieldName := field.Name

			fType := ""
			if len(stringOpts) > 0 && slices.Contains(stringOpts, "type") {
				fType = fmt.Sprintf("%v ", field.Type)

			}

			if fieldValue.Kind() == reflect.Pointer {
				fieldValue = fieldValue.Elem()

			}

			itemPtr := fieldValue
			if !field.IsExported() {
				itemPtr = reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()

			}

			if !itemPtr.IsValid() {
				continue

			}

			if i == 0 {
				newName := objType.Name()
				if newName == "" {
					newName = objName

				}

				structStart := fmt.Sprintf("%s%s{%s", endTab, newName, newline)
				sb.WriteString(structStart)

			} else {
				sb.WriteString(", ")
				sb.WriteString(newline)

			}

			if sValue := stringValue(itemPtr); sValue != "" {
				fmt.Fprintf(sb, "%s%s:%s%s", tab, fieldName, fType, sValue) //, newline)

			} else {
				newTab, newEndTab := tab, endTab
				if fieldValue.Kind() == reflect.Struct && newline != "" {
					newTab, newEndTab = tab+"\t", endTab+"\t"

				}

				structStringHelper(sb, fieldValue, fieldName, newline, newTab, newEndTab, stringOpts...)
				continue

			}

		}

		sb.WriteString(newline)
		sb.WriteString(endTab)
		sb.WriteString("}")

	} else if objValue.Kind() == reflect.Slice {
		fType := ""
		if len(stringOpts) > 0 && slices.Contains(stringOpts, "type") {
			fType = fmt.Sprintf("%s", objValue.Type())

		}

		var sbPtr strings.Builder
		typeVerb := "%s"
		for i := 0; i < objValue.Len(); i++ {
			val := ""
			if !objValue.IsNil() {
				objPtrIndex := objValue.Index(i)
				if objPtrIndex.Kind() == reflect.Pointer {
					objPtrIndex = objPtrIndex.Elem()

				}

				if sValue := stringValue(objPtrIndex); sValue != "" {
					val = fmt.Sprintf("%s", sValue)

				} else {
					newName := ""
					newTab, newEndTab := tab, endTab
					//if objPtrIndex.Kind() == reflect.Struct {
					if newline != "" {
						newName = objName
						newTab, newEndTab = tab+"\t", endTab+"\t"

					}

					if i == 0 {
						sbPtr.WriteString(newline)

					}

					structStringHelper(&sbPtr, objPtrIndex, newName, newline, newTab, newEndTab, stringOpts...)

				}

			}

			if val != "" {
				sbPtr.WriteString(val)

			}

			if i != objValue.Len()-1 {
				sbPtr.WriteString(", ")

			}

			if strings.Contains(sbPtr.String(), "\n") {
				sbPtr.WriteString(newline)

			}

		}

		if objValue.Len() > 0 && strings.Contains(sbPtr.String(), "\n") {
			sbPtr.WriteString(tab)

		}

		fmt.Fprintf(sb, "%s%s%s%s["+typeVerb+"]", checkTab, objName, colon, fType, sbPtr.String())

	} else if objValue.Kind() == reflect.Map {
		mapIter := objValue.MapRange()

		fType := ""
		if len(stringOpts) > 0 && slices.Contains(stringOpts, "type") {
			fType = fmt.Sprintf("%s", objValue.Type())

		}

		fmt.Fprintf(sb, "%s%s%s%s{%s", checkTab, objName, colon, fType, newline)
		thisTab := ""
		if newline != "" {
			thisTab = "\t"

		}

		i := 0
		keyLen := len(objValue.MapKeys())
		for mapIter.Next() {
			k := mapIter.Key() //.Interface()
			v := mapIter.Value()
			if v.Kind() == reflect.Pointer {
				v = v.Elem()

			}

			comma := ""
			if i != keyLen-1 {
				comma = ","

			}

			val := ""
			if sValue := stringValue(v); sValue != "" {
				val = fmt.Sprintf("%s", sValue)

			} else {
				checkName := objName
				if v.Kind() == reflect.Slice {
					fmt.Fprintf(sb, "%s\"%s\":", tab+thisTab, k)
					checkName = ""

				}
				structStringHelper(sb, v, checkName, newline, tab+"\t", endTab+"\t", stringOpts...)

			}

			if v.Kind() != reflect.Slice {
				fmt.Fprintf(sb, "%s\"%s\":%s%s%s", tab+thisTab, k, val, comma, newline)

			} else {
				//fmt.Fprintf(sb, "%s%s%s", tab+thisTab, comma, newline)
				fmt.Fprintf(sb, "%s%s", comma, newline)

			}

			i++

		}

		fmt.Fprintf(sb, "%s}", tab)

	}

}

// stringValue return the value as a string.
//   - objValue is the value to convert to a string
//
// Return
//   - "objValue" if a string (i.e. with quotes "")
//   - objValue if it is a number
//   - true|false if it is a bool
func stringValue(objValue reflect.Value) string {
	var val string
	slog.Debug("POES3", "NumberGeneral(objValue)", NumberGeneral(objValue))
	if objValue.Kind() == reflect.String {
		//val = fmt.Sprintf("%s", objValue)
		val = fmt.Sprintf("\"%s\"", objValue)

	} else if NumberGeneral(objValue) == reflect.Int {
		val = fmt.Sprintf("%d", objValue.Int())

	} else if NumberGeneral(objValue) == reflect.Uint {
		val = fmt.Sprintf("%d", objValue.Uint())

	} else if NumberGeneral(objValue) == reflect.Float32 {
		val = fmt.Sprintf("%f", objValue.Float())

	} else if NumberGeneral(objValue) == reflect.Bool {
		val = fmt.Sprintf("%t", objValue.Bool())

	}

	return val

}

// NumberGeneral generalizes integers and floats.
//   - objValue is the value to generalize
//
// Return
//   - reflect.Int if any of: reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64
//   - reflect.Uint if any of: reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr
//   - reflect.Float32 if any of: reflect.Float32, reflect.Float64
//   - objValue.Kind()
func NumberGeneral(objValue reflect.Value) reflect.Kind {
	switch objValue.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return reflect.Uint

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return reflect.Int

	case reflect.Float32, reflect.Float64:
		return reflect.Float32

	}

	return objValue.Kind()

}
