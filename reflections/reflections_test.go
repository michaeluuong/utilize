package reflections

import (
	"fmt"
	"maps"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Field1       string  `json:"field1"`
	Field2       string  `json:"field2"`
	FieldUint8   uint8   `json:"fielduint8"`
	FieldUint16  uint16  `json:"fielduint16"`
	FieldUint32  uint32  `json:"fielduint32"`
	FieldUint64  uint64  `json:"fielduint64"`
	FieldInt     int     `json:"fieldint"`
	FieldInt64   int     `json:"fieldint64"`
	FieldBool    bool    `json:"fieldbool"`
	FieldFloat32 float32 `json:"fieldfloat32"`
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()

	os.Exit(code)

}

func setup() {

}

func teardown() {

}

// go test -benchtime=1s -bench . -cpuprofile cpu.prof
// go tool pprof cpu.prof
func BenchmarkReflections(b *testing.B) {
	test := &TestStruct{Field1: "first field"}

	for b.Loop() {
		FunctionName(0)
		FuncSig(0, "")
		ReflectFieldByName(test, "Field1")
		field, _ := ReflectFieldByIndex(test, 1)
		SetStructFieldByType(field, "second field")
		_, _ = FieldAndTagNames(test, "json")

	}

}

func basename(filename string) string {
	if filename == "" {
		return ""

	}

	parts := strings.Split(filepath.Base(filename), ".")
	return parts[len(parts)-1]

}

func TestAnyType_main(t *testing.T) {
	type myTest struct{}
	myTestType := reflect.TypeFor[myTest]()
	var actualReflectType reflect.Type

	//--------------------------------------------------------------------------------------------
	// Struct
	actualReflectType = AnyType(myTest{})
	assert.Equal(t, myTestType, actualReflectType)

	actualReflectType = AnyType(&myTest{})
	assert.Equal(t, myTestType, actualReflectType)

	//--------------------------------------------------------------------------------------------
	// Map
	expectedType := reflect.TypeFor[map[string]string]()
	actualReflectType = AnyType(map[string]string{})
	assert.Equal(t, expectedType, actualReflectType)

	actualReflectType = AnyType(&map[string]string{})
	assert.Equal(t, expectedType, actualReflectType)

	//--------------------------------------------------------------------------------------------
	// nil
	actualReflectType = AnyType(nil)
	assert.Nil(t, actualReflectType)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestCopyStruct_main(t *testing.T) {
	type SourceStruct struct {
		Field1 string
		Field2 int
		Field3 bool
	}

	type DestinationStruct struct {
		Field1 string
		Field2 int
		Field3 bool
	}

	//--------------------------------------------------------------------------------------------
	// Happy
	source := SourceStruct{
		Field1: "field1",
		Field2: 70,
		Field3: true,
	}

	var source2 SourceStruct
	CopyStruct(any(source), any(&source2))
	assert.Equal(t, source, source2)

	// Pointer
	var sourcePointer *SourceStruct
	sourcePointer = &source
	sourceDestination2 := &SourceStruct{}
	CopyStruct(any(sourcePointer), any(sourceDestination2))

	sourceDestinationPointer := &SourceStruct{}
	CopyStruct(any(source), any(sourceDestinationPointer))
	assert.Equal(t, source, *sourceDestinationPointer)

	var sourcePointer2 *SourceStruct = &SourceStruct{}
	CopyStruct(any(source), any(sourcePointer2))
	assert.Equal(t, source, *sourcePointer2)

	//--------------------------------------------------------------------------------------------
	// Differenct types
	destination := DestinationStruct{}

	CopyStruct(any(source), any(&destination))
	assert.EqualValues(t, source, destination)

	// Pointer
	destination2 := &DestinationStruct{}
	CopyStruct(any(source), any(destination2))
	assert.EqualValues(t, source, *destination2)

	//--------------------------------------------------------------------------------------------
	// Only overwrite zero values
	sourceOverwrite := &SourceStruct{Field1: "copy field", Field2: 69, Field3: true}
	destinationOverwrite := &SourceStruct{Field2: 70}
	CopyStruct(any(sourceOverwrite), any(destinationOverwrite), true)
	assert.Equal(t, destinationOverwrite.Field1, sourceOverwrite.Field1)
	assert.NotEqual(t, destinationOverwrite.Field2, sourceOverwrite.Field2)
	assert.Equal(t, destinationOverwrite.Field3, sourceOverwrite.Field3)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestFieldAndTagNames_main(t *testing.T) {
	var fields map[string]bool
	var tagToName map[string]string
	expectedFields := []string{"Field1", "Field2", "FieldUint8", "FieldUint16", "FieldUint32", "FieldUint64",
		"FieldInt", "FieldInt64", "FieldBool", "FieldFloat32"}
	slices.Sort(expectedFields)

	expectedTags := []string{"field1", "field2", "fielduint8", "fielduint16", "fielduint32", "fielduint64",
		"fieldint", "fieldint64", "fieldbool", "fieldfloat32"}
	slices.Sort(expectedTags)

	var actualFields []string
	var actualTags []string

	//--------------------------------------------------------------------------------------------
	// Happy
	test := TestStruct{}
	fields, tagToName = FieldAndTagNames(test, "json")

	actualFields = slices.Collect(maps.Keys(fields))
	slices.Sort(actualFields)
	assert.Equal(t, expectedFields, actualFields)

	actualTags = slices.Collect(maps.Keys(tagToName))
	slices.Sort(actualTags)
	assert.Equal(t, expectedTags, actualTags)

	//--------------------------------------------------------------------------------------------
	// Pointer
	var testPointer *TestStruct
	testPointer = &TestStruct{}
	fields, tagToName = FieldAndTagNames(testPointer, "json")

	actualFields = slices.Collect(maps.Keys(fields))
	slices.Sort(actualFields)
	assert.Equal(t, expectedFields, actualFields)

	actualTags = slices.Collect(maps.Keys(tagToName))
	slices.Sort(actualTags)
	assert.Equal(t, expectedTags, actualTags)

	//--------------------------------------------------------------------------------------------
	// Bad tag
	_, tagToName = FieldAndTagNames(test, "gson")
	assert.Empty(t, tagToName)

	//--------------------------------------------------------------------------------------------
	// Map
	fields, tagToName = FieldAndTagNames(map[string]string{}, "json")
	assert.Empty(t, fields)
	assert.Empty(t, tagToName)

	//--------------------------------------------------------------------------------------------
	// Empty struct
	fields, _ = FieldAndTagNames(nil, "json")
	assert.Empty(t, fields)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestFuncSig_main(t *testing.T) {
	var sig string

	//--------------------------------------------------------------------------------------------
	// Happy
	sig = basename(FuncSig(1, "_test()"))
	assert.Equal(t, "TestFuncSig_main_test()", sig)

	//--------------------------------------------------------------------------------------------
	// No suffix
	sig = basename(FuncSig(1))
	assert.Equal(t, "TestFuncSig_main() - ", sig)

	//--------------------------------------------------------------------------------------------
	// Invalid stack frame number
	sig = basename(FuncSig(-10))
	assert.Equal(t, "Caller() - ", sig)

	sig = basename(FuncSig(1000, ""))
	assert.Empty(t, sig)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestFunctionName_main(t *testing.T) {
	var name string

	//--------------------------------------------------------------------------------------------
	// Happy
	name = basename(FunctionName(0))
	assert.Equal(t, "TestFunctionName_main", name)

	//--------------------------------------------------------------------------------------------
	// Invalid stack frame number
	name = basename(FunctionName(-20))
	assert.Equal(t, "Caller", name)

	name = basename(FunctionName(1000))
	assert.Empty(t, name)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

type TestObject struct {
	TestMap map[string]string
	testMap map[string]string
}

func TestInitialzeStruct_main(t *testing.T) {
	testObject := &TestObject{}
	assert.Nil(t, testObject.TestMap)
	assert.Nil(t, testObject.testMap)

	testObject2 := TestObject{}

	//--------------------------------------------------------------------------------------------
	// Happy
	InitializeStruct(testObject)
	assert.NotNil(t, testObject.TestMap)
	assert.NotNil(t, testObject.testMap)

	InitializeStruct(testObject2)
	assert.NotNil(t, testObject.TestMap)
	assert.NotNil(t, testObject.testMap)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestIsPointer_main(t *testing.T) {
	var isPointer bool
	var pointer *int

	//--------------------------------------------------------------------------------------------
	// Happy
	isPointer = IsPointer(&TestStruct{})
	assert.True(t, isPointer)

	isPointer = IsPointer(TestStruct{})
	assert.False(t, isPointer)

	isPointer = IsPointer(pointer)
	assert.True(t, isPointer)

	//--------------------------------------------------------------------------------------------
	// Map
	isPointer = IsPointer(map[string]bool{})
	assert.False(t, isPointer)

	var setPointer *map[string]bool
	isPointer = IsPointer(setPointer)
	assert.True(t, isPointer)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestGetType_main(t *testing.T) {
	//--------------------------------------------------------------------------------------------
	// Struct
	assert.Equal(t, reflect.Struct, GetType(&struct{}{}))

	//--------------------------------------------------------------------------------------------
	// Map
	assert.Equal(t, reflect.Map, GetType(&map[string]bool{}))

	//--------------------------------------------------------------------------------------------
	// Slice
	var mySlice []string
	assert.Equal(t, reflect.Slice, GetType(mySlice))

	//--------------------------------------------------------------------------------------------
	// int
	var myInt int
	assert.Equal(t, reflect.Int, GetType(myInt))

	var myFloat float32
	assert.Equal(t, reflect.Float32, GetType(myFloat))

	//--------------------------------------------------------------------------------------------
	// nil
	assert.Equal(t, reflect.Invalid, GetType(nil))

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestIsStruct_main(t *testing.T) {
	var isStruct bool

	//--------------------------------------------------------------------------------------------
	var test TestStruct
	isStruct = IsStruct(test)
	assert.True(t, isStruct)

	isStruct = IsStruct(&TestStruct{})
	assert.True(t, isStruct)

	isStruct = IsStruct(map[string]string{})
	assert.False(t, isStruct)

	var num int = 69
	isStruct = IsStruct(num)
	assert.False(t, isStruct)

	var num2 *int
	isStruct = IsStruct(num2)
	assert.False(t, isStruct)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestEmptyStructInstance_main(t *testing.T) {

	//--------------------------------------------------------------------------------------------
	// Happy
	expected := EmptyStructInstance[TestStruct]()
	assert.Equal(t, &TestStruct{}, expected)

	type myInt int
	pNil, ok := EmptyStructInstance[myInt]().(myInt)
	assert.False(t, ok)
	assert.Equal(t, 0, int(pNil))

	//--------------------------------------------------------------------------------------------
	// Modify
	pExpected := &TestStruct{Field1: "Test Ease"}
	pActual := EmptyStructInstance[TestStruct]().(*TestStruct)
	pActual.Field1 = "Test Ease"
	assert.Equal(t, pActual, pExpected)

	pActual2 := EmptyStructInstance[*TestStruct]().(*TestStruct)
	pActual2.Field1 = "Test Ease"
	assert.Equal(t, pActual2, pExpected)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestNewMapInstance_main(t *testing.T) {

	newMap := NewMapInstance[string, int]()
	newMap["Field 1"] = 2
	assert.Equal(t, newMap["Field 1"], 2)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestReflectFieldByIndex_main(t *testing.T) {
	var field reflect.Value
	var err error
	var test *TestStruct

	//--------------------------------------------------------------------------------------------
	// Happy
	test = &TestStruct{}
	field, err = ReflectFieldByIndex(test, 1)
	field2Value := "field 2"
	if err == nil {
		field.SetString(field2Value)

	}
	assert.Equalf(t, field2Value, test.Field2, "Reflectively set Field2=\"%s\" so they should be equal", field2Value)

	//--------------------------------------------------------------------------------------------
	// Not a struct
	field, err = ReflectFieldByIndex(map[string]string{}, 0)
	assert.EqualError(t, err, "obj is not a struct")

	//--------------------------------------------------------------------------------------------
	// Bad index
	test = &TestStruct{}
	field, err = ReflectFieldByIndex(test, -1)
	assert.ErrorContains(t, err, "must be a valid field number in obj")

	//--------------------------------------------------------------------------------------------
	// nil
	field, err = ReflectFieldByIndex(nil, 1)
	assert.EqualError(t, err, "obj is not a struct")

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestReflectFieldByName_main(t *testing.T) {
	var field reflect.Value
	var err error
	var fieldName string
	var test *TestStruct

	//--------------------------------------------------------------------------------------------
	// Happy
	test = &TestStruct{}
	fieldName = "Field2"
	field, err = ReflectFieldByName(test, fieldName)
	field2Value := "field 2"
	if err == nil {
		field.SetString(field2Value)

	}
	assert.Equalf(t, field2Value, test.Field2, "Reflectively set Field2=\"%s\" so they should be equal", field2Value)

	//--------------------------------------------------------------------------------------------
	// Not a struct
	field, err = ReflectFieldByName(map[string]string{}, "Field1")
	assert.EqualError(t, err, "obj is not a struct")

	//--------------------------------------------------------------------------------------------
	// Field does not exist
	test = &TestStruct{}
	fieldName = "Dne"
	field, err = ReflectFieldByName(test, fieldName)
	assert.ErrorContainsf(t, err, "fieldName Dne is not a valid field in obj", "Field %s should not be a valid field", fieldName)

	//--------------------------------------------------------------------------------------------
	// nil
	field, err = ReflectFieldByName(nil, "")
	assert.EqualError(t, err, "obj is not a struct")

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestSetStructField_main(t *testing.T) {
	//	var field reflect.Value
	var err error
	var test *TestStruct

	//--------------------------------------------------------------------------------------------
	// Happy
	test = &TestStruct{}
	_ = SetStructField(test, "Field2", "stringfield")
	assert.Equal(t, "stringfield", test.Field2)

	//--------------------------------------------------------------------------------------------
	// Not a struct
	err = SetStructField(map[string]string{}, "Field2", "stringfield")
	assert.EqualError(t, err, "obj is not a struct")

	//--------------------------------------------------------------------------------------------
	// Incorrect field
	test = &TestStruct{}
	err = SetStructField(test, "FieldDNE", "stringfield")
	assert.ErrorContains(t, err, "fieldName FieldDNE is not a valid field in obj")

	//--------------------------------------------------------------------------------------------
	// nil
	err = SetStructField(nil, "Field2", "stringfield")
	assert.EqualError(t, err, "obj is not a struct")

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

func TestSetStructFieldByType_main(t *testing.T) {
	var err error
	var test *TestStruct

	var maxUint8 uint8 = uint8(math.MaxUint8)
	var maxUint8String string = strconv.FormatUint(uint64(maxUint8), 10)
	var maxUint16 uint16 = uint16(math.MaxUint16)
	var maxUint16String string = strconv.FormatUint(uint64(maxUint16), 10)
	var maxUint32 uint32 = uint32(math.MaxUint32)
	var maxUint32String string = strconv.FormatUint(uint64(maxUint32), 10)
	var maxUint64 uint64 = uint64(math.MaxUint64)
	var maxUint64String string = strconv.FormatUint(maxUint64, 10)

	var maxFloat64 float64 = float64(math.MaxFloat64)
	var maxFloat64String string = strconv.FormatFloat(maxFloat64, 'f', -1, 64)
	maxFloat64String = "69" + maxFloat64String

	//--------------------------------------------------------------------------------------------
	// uint
	test = &TestStruct{}
	err = SetStructField(test, "FieldUint8", maxUint8String)
	assert.Equal(t, maxUint8, test.FieldUint8)

	err = SetStructField(test, "FieldUint16", maxUint16String)
	assert.Equal(t, maxUint16, test.FieldUint16)

	err = SetStructField(test, "FieldUint32", maxUint32String)
	assert.Equal(t, maxUint32, test.FieldUint32)

	err = SetStructField(test, "FieldUint64", maxUint64String)
	assert.Equal(t, maxUint64, test.FieldUint64)

	err = SetStructField(test, "FieldUint8", maxUint64String)
	assert.ErrorContains(t, err, ": value out of range")

	//--------------------------------------------------------------------------------------------
	// int
	test = &TestStruct{}
	err = SetStructField(test, "FieldInt", "69")
	assert.Equal(t, 69, test.FieldInt)

	err = SetStructField(test, "FieldInt", maxUint64String)
	assert.ErrorContains(t, err, "value out of range")

	err = SetStructField(test, "FieldInt64", "6969696969696969")
	assert.Equal(t, 6969696969696969, test.FieldInt64)

	err = SetStructField(test, "FieldInt64", maxUint64String)
	assert.ErrorContains(t, err, "value out of range")

	//--------------------------------------------------------------------------------------------
	// bool
	err = SetStructField(test, "FieldBool", "true")
	assert.True(t, test.FieldBool)

	err = SetStructField(test, "FieldBool", "falsify")
	assert.ErrorContains(t, err, "invalid syntax")

	//--------------------------------------------------------------------------------------------
	// Float
	err = SetStructField(test, "FieldFloat32", "69.6969")
	assert.InDelta(t, 69.6969, test.FieldFloat32, .001)

	err = SetStructField(test, "FieldFloat32", maxFloat64String)
	assert.ErrorContains(t, err, "value out of range")

	//--------------------------------------------------------------------------------------------
	// Cleanup

}

type TestStructString struct {
	field1  string
	Field2  int
	Field3  []string
	Field4  []*string
	field5  []*string
	Field6  map[string]string
	Field7  map[string]*string
	Field8  []int
	Field9  []map[string]string
	Field10 []map[string]*string
	Field11 [][]string
	Field12 struct {
		Field12a string
	}
	Field13    float32
	field14    bool
	StringOpts []string
}

func (t *TestStructString) setStringOpts(stringOpts ...string) {
	t.StringOpts = stringOpts
	_ = t.String()
}

func (t TestStructString) String() string {
	return StructString(t, t.StringOpts...)

}

type TestStructStringPointer struct {
	field1 bool
	Field2 bool
	Field3 *string
	Field4 map[string][]string
}

func (t *TestStructStringPointer) String() string {
	return StructString(t, "newline")

}

func TestStructString_main(t *testing.T) {
	testString := "test"
	var tStr *string = &testString
	testStructString := &TestStructString{
		field1:  "field1",
		Field2:  10,
		Field3:  []string{"field3 test 1", "field3 test 2"},
		Field4:  []*string{tStr},
		field5:  []*string{tStr},
		Field6:  map[string]string{"test 1": "test 1 val"},
		Field7:  map[string]*string{"test 1": tStr},
		Field8:  []int{1, 2, 3},
		Field9:  []map[string]string{{"Field9 1": "value1"}, {"field9 2": "value2"}},
		Field10: []map[string]*string{{"Field9 1": tStr}, {"field9 2": tStr}},
		Field11: [][]string{{"1", "2"}, {"3", "4"}},
		Field12: struct{ Field12a string }{Field12a: "field12a"},
		Field13: 1246.00,
		field14: true,
		//StringOpts: []string{"type", "newline", "\t"},
	}

	var actualString, expectedString string

	//--------------------------------------------------------------------------------------------
	// Happy
	actualString = fmt.Sprintf("%s", testStructString)
	expectedString = `TestStructString{field1:"field1", Field2:10, Field3:["field3 test 1", "field3 test 2"], Field4:["test"], field5:["test"], Field6:{"test 1":"test 1 val"}, Field7:{"test 1":"test"}, Field8:[1, 2, 3], Field9:[{"Field9 1":"value1"}, {"field9 2":"value2"}], Field10:[{"Field9 1":"test"}, {"field9 2":"test"}], Field11:[["1", "2"], ["3", "4"]], Field12{Field12a:"field12a"}, Field13:1246.000000, field14:true, StringOpts:[]}`
	assert.Equal(t, expectedString, actualString)

	//--------------------------------------------------------------------------------------------
	// type
	testStructString.setStringOpts("type")
	actualString = fmt.Sprintf("%s", testStructString)
	expectedString = `TestStructString{field1:string "field1", Field2:int 10, Field3:[]string["field3 test 1", "field3 test 2"], Field4:[]*string["test"], field5:[]*string["test"], Field6:map[string]string{"test 1":"test 1 val"}, Field7:map[string]*string{"test 1":"test"}, Field8:[]int[1, 2, 3], Field9:[]map[string]string[map[string]string{"Field9 1":"value1"}, map[string]string{"field9 2":"value2"}], Field10:[]map[string]*string[map[string]*string{"Field9 1":"test"}, map[string]*string{"field9 2":"test"}], Field11:[][]string[[]string["1", "2"], []string["3", "4"]], Field12{Field12a:string "field12a"}, Field13:float32 1246.000000, field14:bool true, StringOpts:[]string["type"]}`
	assert.Equal(t, expectedString, actualString)

	//--------------------------------------------------------------------------------------------
	// newline
	testStructString.setStringOpts("newline")
	actualString = fmt.Sprintf("%s", testStructString)
	expectedString = "TestStructString{\n\tfield1:\"field1\", \n\tField2:10, \n\tField3:[\"field3 test 1\", \"field3 test 2\"], \n\tField4:[\"test\"], \n\tfield5:[\"test\"], \n\tField6:{\n\t\t\"test 1\":\"test 1 val\"\n\t}, \n\tField7:{\n\t\t\"test 1\":\"test\"\n\t}, \n\tField8:[1, 2, 3], \n\tField9:[\n\t\tField9:{\n\t\t\t\"Field9 1\":\"value1\"\n\t\t}, \n\t\tField9:{\n\t\t\t\"field9 2\":\"value2\"\n\t\t}\n\t], \n\tField10:[\n\t\tField10:{\n\t\t\t\"Field9 1\":\"test\"\n\t\t}, \n\t\tField10:{\n\t\t\t\"field9 2\":\"test\"\n\t\t}\n\t], \n\tField11:[\n\t\tField11:[\"1\", \"2\"], \n\t\tField11:[\"3\", \"4\"]\n\t], \n\tField12{\n\t\tField12a:\"field12a\"\n\t}, \n\tField13:1246.000000, \n\tfield14:true, \n\tStringOpts:[\"newline\"]\n}"
	assert.Equal(t, expectedString, actualString)

	//fmt.Printf("%s\n", testStructString)

	//--------------------------------------------------------------------------------------------
	// newline, tab
	testStructString.setStringOpts("newline", "\t\t")
	actualString = fmt.Sprintf("%s", testStructString)
	expectedString = "\tTestStructString{\n\t\tfield1:\"field1\", \n\t\tField2:10, \n\t\tField3:[\"field3 test 1\", \"field3 test 2\"], \n\t\tField4:[\"test\"], \n\t\tfield5:[\"test\"], \n\t\tField6:{\n\t\t\t\"test 1\":\"test 1 val\"\n\t\t}, \n\t\tField7:{\n\t\t\t\"test 1\":\"test\"\n\t\t}, \n\t\tField8:[1, 2, 3], \n\t\tField9:[\n\t\t\tField9:{\n\t\t\t\t\"Field9 1\":\"value1\"\n\t\t\t}, \n\t\t\tField9:{\n\t\t\t\t\"field9 2\":\"value2\"\n\t\t\t}\n\t\t], \n\t\tField10:[\n\t\t\tField10:{\n\t\t\t\t\"Field9 1\":\"test\"\n\t\t\t}, \n\t\t\tField10:{\n\t\t\t\t\"field9 2\":\"test\"\n\t\t\t}\n\t\t], \n\t\tField11:[\n\t\t\tField11:[\"1\", \"2\"], \n\t\t\tField11:[\"3\", \"4\"]\n\t\t], \n\t\tField12{\n\t\t\tField12a:\"field12a\"\n\t\t}, \n\t\tField13:1246.000000, \n\t\tfield14:true, \n\t\tStringOpts:[\"newline\", \"\t\t\"]\n\t}"
	assert.Equal(t, expectedString, actualString)

	//--------------------------------------------------------------------------------------------
	// type, newline, tab
	testStructString.setStringOpts("type", "newline", "\t")
	actualString = fmt.Sprintf("%s", testStructString)
	expectedString = "TestStructString{\n\tfield1:string \"field1\", \n\tField2:int 10, \n\tField3:[]string[\"field3 test 1\", \"field3 test 2\"], \n\tField4:[]*string[\"test\"], \n\tfield5:[]*string[\"test\"], \n\tField6:map[string]string{\n\t\t\"test 1\":\"test 1 val\"\n\t}, \n\tField7:map[string]*string{\n\t\t\"test 1\":\"test\"\n\t}, \n\tField8:[]int[1, 2, 3], \n\tField9:[]map[string]string[\n\t\tField9:map[string]string{\n\t\t\t\"Field9 1\":\"value1\"\n\t\t}, \n\t\tField9:map[string]string{\n\t\t\t\"field9 2\":\"value2\"\n\t\t}\n\t], \n\tField10:[]map[string]*string[\n\t\tField10:map[string]*string{\n\t\t\t\"Field9 1\":\"test\"\n\t\t}, \n\t\tField10:map[string]*string{\n\t\t\t\"field9 2\":\"test\"\n\t\t}\n\t], \n\tField11:[][]string[\n\t\tField11:[]string[\"1\", \"2\"], \n\t\tField11:[]string[\"3\", \"4\"]\n\t], \n\tField12{\n\t\tField12a:string \"field12a\"\n\t}, \n\tField13:float32 1246.000000, \n\tfield14:bool true, \n\tStringOpts:[]string[\"type\", \"newline\", \"\t\"]\n}"
	assert.Equal(t, expectedString, actualString)

	//--------------------------------------------------------------------------------------------
	// Pointer
	testStructStringPointer := &TestStructStringPointer{field1: true, Field2: false, Field3: tStr, Field4: map[string][]string{"Field 4 key": {"test1", "test2"}}}
	actualString = fmt.Sprintf("%s", testStructStringPointer)
	expectedString = "TestStructStringPointer{\n\tfield1:true, \n\tField2:false, \n\tField3:\"test\", \n\tField4:{\n\t\t\"Field 4 key\":[\"test1\", \"test2\"]\n\t}\n}"
	assert.Equal(t, expectedString, actualString)

	// Cover a comma for fun
	testStructStringPointer = &TestStructStringPointer{Field4: map[string][]string{"Field4 1 key": {"test1", "test2"}, "Field4 2 key": {"tp1", "tp2", "tp3"}}}
	actualString = fmt.Sprintf("%s", testStructStringPointer)
	expectedString = "TestStructStringPointer{\n\tfield1:false, \n\tField2:false, \n\tField4:{\n\t\t\"Field4 1 key\":[\"test1\", \"test2\"],\n\t\t\"Field4 2 key\":[\"tp1\", \"tp2\", \"tp3\"]\n\t}\n}"
	//assert.Equal(t, expectedString, actualString)

	//--------------------------------------------------------------------------------------------
	// Cleanup

}
