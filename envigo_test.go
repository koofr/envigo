package envigo

import (
	"os"
	"testing"
)

func TestEnvigo(t *testing.T) {
	type SubConfig struct {
		Int       int
		Uint      uint
		Float     float64
		Bool      bool
		BoolOne   bool
		Default   string
		anonymous string
	}

	type Config struct {
		Str    string
		Sub    SubConfig
		SubPtr *SubConfig
	}

	getenv := func(key string) (string, bool) {
		switch key {
		case "STR":
			return "string", true
		case "SUB_INT":
			return "-42", true
		case "SUB_UINT":
			return "42", true
		case "SUB_FLOAT":
			return "3.14", true
		case "SUB_BOOL":
			return "true", true
		case "SUB_BOOLONE":
			return "1", true
		case "SUB_ANONYMOUS":
			return "value", true
		case "SUBPTR_INT":
			return "52", true
		default:
			return "", false
		}
	}

	c := Config{}

	err := Envigo(&c, "", getenv)

	if err != nil {
		t.Errorf("Envigo error: %s", err)
	}

	if c.Str != "string" {
		t.Fail()
	}

	if c.Sub.Int != -42 {
		t.Fail()
	}

	if c.Sub.Uint != 42 {
		t.Fail()
	}

	if c.Sub.Float != 3.14 {
		t.Fail()
	}

	if c.Sub.Bool != true {
		t.Fail()
	}

	if c.Sub.BoolOne != true {
		t.Fail()
	}

	if c.Sub.Default != "" {
		t.Fail()
	}

	if c.Sub.anonymous != "" {
		t.Fail()
	}

	if c.SubPtr.Int != 52 {
		t.Fail()
	}
}

func TestEnvigoReal(t *testing.T) {
	type Config struct {
		Path string
	}

	c := Config{}

	err := Envigo(&c, "", EnvironGetter())

	if err != nil {
		t.Errorf("Envigo error: %s", err)
	}

	if c.Path != os.Getenv("PATH") {
		t.Fail()
	}
}

func TestEnvigoPrefix(t *testing.T) {
	type Config struct {
		Str string
	}

	getenv := func(key string) (string, bool) {
		switch key {
		case "PREFIX_STR":
			return "string", true
		default:
			return "", false
		}
	}

	c := Config{}

	err := Envigo(&c, "PREFIX", getenv)

	if err != nil {
		t.Errorf("Envigo error: %s", err)
	}

	if c.Str != "string" {
		t.Fail()
	}
}

func TestEnvigoParseError(t *testing.T) {
	type SubConfig struct {
		Int int
	}

	type Config struct {
		Sub    SubConfig
		SubPtr *SubConfig
		Uint   uint
		Float  float64
	}

	getenv := func(key string) (string, bool) {
		switch key {
		case "SUB_INT":
			return "invalid", true
		default:
			return "", false
		}
	}

	c := Config{}

	err := Envigo(&c, "", getenv)

	if err == nil {
		t.Errorf("Envigo parse sub int error should occur")
	}

	getenv = func(key string) (string, bool) {
		switch key {
		case "SUBPTR_INT":
			return "invalid", true
		default:
			return "", false
		}
	}

	err = Envigo(&c, "", getenv)

	if err == nil {
		t.Errorf("Envigo parse subptr int error should occur")
	}

	getenv = func(key string) (string, bool) {
		switch key {
		case "UINT":
			return "invalid", true
		default:
			return "", false
		}
	}

	err = Envigo(&c, "", getenv)

	if err == nil {
		t.Errorf("Envigo parse uint error should occur")
	}

	getenv = func(key string) (string, bool) {
		switch key {
		case "FLOAT":
			return "invalid", true
		default:
			return "", false
		}
	}

	err = Envigo(&c, "", getenv)

	if err == nil {
		t.Errorf("Envigo parse float error should occur")
	}
}

func TestEnvigoNonstructError(t *testing.T) {
	getenv := func(key string) (string, bool) {
		return "", false
	}

	c := ""

	err := Envigo(&c, "", getenv)

	if err == nil {
		t.Errorf("Envigo error should occur")
	}
}

func TestEnvironGetter(t *testing.T) {
	getter := EnvironGetter()

	path, ok := getter("PATH")

	if p := os.Getenv("PATH"); !ok || path != p {
		t.Errorf("EnvironGetter failed for PATH: %s != %s", path, p)
	}
}
