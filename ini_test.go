package ini

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestLoad(t *testing.T) {
	src := `
  # Comments are ignored

  herp = derp

  [foo]
  hello=world
  whitespace should   =   not matter   
  ; sneaky semicolon-style comment
  multiple = equals = signs

  [bar]
  this = that`

	file, err := Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	check := func(section, key, expect string) {
		checkStr(t, &file, section, key, expect)
	}

	check("", "herp", "derp")
	check("foo", "hello", "world")
	check("foo", "whitespace should", "not matter")
	check("foo", "multiple", "equals = signs")
	check("bar", "this", "that")
}

func TestBoolFalse(t *testing.T) {
	src := `
  # Comments are ignored

  foo =
  bar = 0
  fox = no
  baz = False
  `

	file, err := Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	check := func(key string) {
		var (
			val bool
			ok  bool
		)
		val, ok = file.GetBool("", key)
		if !ok {
			t.Errorf("GetBool(%q): not read successfully", key)
		}
		if val {
			t.Errorf("GetBool(%q): expected false not true", key)
		}
	}

	check("foo")
	check("bar")
	check("fox")
	check("baz")
}

func TestBoolTrue(t *testing.T) {
	src := `
  # Comments are ignored

  foo = 1
  bar = yes
  fox = true
  baz = True
  `

	file, err := Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	check := func(key string) {
		var (
			val bool
			ok  bool
		)
		val, ok = file.GetBool("", key)
		if !ok {
			t.Errorf("GetBool(%q): not read successfully", key)
		}
		if !val {
			t.Errorf("GetBool(%q): expected true not false", key)
		}
	}

	check("foo")
	check("bar")
	check("fox")
	check("baz")
}

func TestBoolInvalid(t *testing.T) {
	src := `
  # Comments are ignored

  foo = 11
  bar = si
  fox = 1.0
  baz = 0.0
  `

	file, err := Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	check := func(key string) {
		_, ok := file.GetBool("", key)
		if ok {
			t.Errorf("GetBool(%q): should not be parsed", key)
		}
	}

	check("foo")
	check("bar")
	check("fox")
	check("baz")
}

func TestInteger(t *testing.T) {
	src := `
  foo = 1
  bar = 0100
  fox = -67
  baz = 8001
  `

	file, err := Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	check := func(key string, expected int) {
		val, ok := file.GetInt("", key)
		if !ok {
			t.Errorf("GetInt(%q): not read successfully", key)
		}
		if val != expected {
			t.Errorf("GetInt(%q): expected %d not %d", key, expected, val)
		}
	}

	check("foo", 1)
	check("bar", 100)
	check("fox", -67)
	check("baz", 8001)
}

func TestIntegerInvalid(t *testing.T) {
	src := `
  foo = 1.0
  bar = 1,000
  fox = 0x10
  baz = 1.
  `

	file, err := Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}

	check := func(key string) {
		_, ok := file.GetInt("", key)
		if ok {
			t.Errorf("GetInt(%q): should not have been readable", key)
		}
	}

	check("foo")
	check("bar")
	check("fox")
	check("baz")
}

func checkStr(t *testing.T, file *File, section, key, expect string) {
	if value, _ := file.Get(section, key); value != expect {
		t.Errorf("Get(%q, %q): expected %q, got %q", section, key, expect, value)
	}
}

func checkArr(t *testing.T, file *File, section, key string, expect []string) {
	value, ok := file.GetArr(section, key)
	if !ok {
		t.Errorf("Get(%q, %q): expected value but not found", section, key)
	}
	if len(value) != len(expect) {
		t.Errorf("Get(%q, %q): expected %d items found, got %d", section, key, len(expect), len(value))
	}
	for curKey, thisVal := range expect {
		if thisVal != value[curKey] {
			t.Errorf("Get(%q, %q): expected %s at %d, got %s", section, key, thisVal, curKey, value[curKey])
		}
	}
}

func TestArray(t *testing.T) {
	var (
		file File
		src  string
		err  error
	)
	check := func(section, key string, expect []string) {
		checkArr(t, &file, section, key, expect)
	}

	checkStr := func(section, key, expect string) {
		checkStr(t, &file, section, key, expect)
	}

	src = `
foo [] = bar`
	file, err = Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	check("", "foo", []string{"bar"})

	src = `
[section]
foo[] = bar
foo[] = fox`
	file, err = Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	check("section", "foo", []string{"bar", "fox"})
	src = `
[section]
foo = baz
foo[] = fox
foo[] = bar
herp[] = derp
`
	file, err = Load(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	check("section", "foo", []string{"fox", "bar"})
	checkStr("section", "foo", "baz")
}

func TestSyntaxError(t *testing.T) {
	src := `
  # Line 2
  [foo]
  bar = baz
  # Here's an error on line 6:
  wut?
  herp = derp`
	_, err := Load(strings.NewReader(src))
	t.Logf("%T: %v", err, err)
	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	syntaxErr, ok := err.(ErrSyntax)
	if !ok {
		t.Fatal("expected an error of type ErrSyntax")
	}
	if syntaxErr.Line != 6 {
		t.Fatal("incorrect line number")
	}
	if syntaxErr.Source != "wut?" {
		t.Fatal("incorrect source")
	}
}

func TestDefinedSectionBehaviour(t *testing.T) {
	check := func(src string, expect File) {
		file, err := Load(strings.NewReader(src))
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(file, expect) {
			t.Errorf("expected %v, got %v", expect, file)
		}
	}
	testFile := File{}
	// No sections for an empty file
	check("", testFile)
	// Default section only if there are actually values for it
	testFile.Set("", "foo","bar")
	check("foo=bar", testFile)
	// User-defined sections should always be present, even if empty
	check("[a]\n[b]\nfoo=bar", File{
		"a": makeSection(StringSection{}),
		"b": makeSection(StringSection{"foo": "bar"}),
	})
	check("foo=bar\n[a]\nthis=that", File{
		"":  makeSection(StringSection{"foo": "bar"}),
		"a": makeSection(StringSection{"this": "that"}),
	})
}

func TestWrite(t *testing.T) {
	testIni := File{}
	testIni.Set("section1", "option1", "value1")
	testIni.Set("section1", "option2", "value2")
	testIni.Set("section2", "option3", "value3")
	testIni.Set("section2", "option4", "value4")

	err := testIni.WriteFile("test_write_out.ini")
	if err != nil {
		t.Fatal("Unable to write ini file")
	}
	in, err := os.Open("test_write.ini")
	if err != nil {
		t.Fatal("Unable to open comparison file")
	}
	defer in.Close()
	out, err := os.Open("test_write_out.ini")
	if err != nil {
		t.Fatal("Unable to open comparison file")
	}
	defer out.Close()
	sampleStr := make([]byte, 1024)
	actualStr := make([]byte, 1024)
	sourceBytesRead, err := in.Read(sampleStr)
	if err != nil {
		t.Fatal("Unable to read comparison file")
	}
	newBytesRead, err := out.Read(actualStr)
	if err != nil {
		t.Fatal("Unable to read new file")
	}
	if sourceBytesRead != newBytesRead {
		t.Fatal("Written file differs in length from expected")
	}
	for curPos, curChar := range sampleStr {
		if curChar != actualStr[curPos] {
			t.Error(fmt.Sprintf("Mismatch %q vs %q as char %d", curChar, actualStr[curChar], curPos))
		}
	}
}
