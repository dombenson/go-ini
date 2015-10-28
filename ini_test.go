package ini

import (
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
	// No sections for an empty file
	check("", File{})
	// Default section only if there are actually values for it
	check("foo=bar", File{"": MakeSection(StringSection{"foo": "bar"})})
	// User-defined sections should always be present, even if empty
	check("[a]\n[b]\nfoo=bar", File{
		"a": MakeSection(StringSection{}),
		"b": MakeSection(StringSection{"foo": "bar"}),
	})
	check("foo=bar\n[a]\nthis=that", File{
		"":  MakeSection(StringSection{"foo": "bar"}),
		"a": MakeSection(StringSection{"this": "that"}),
	})
}
