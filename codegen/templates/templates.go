//go:generate go run ./inliner/inliner.go

package templates

import (
	"bytes"
	"fmt"
	"go/types"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/tinhtran24/gqlgen/internal/imports"

	"github.com/pkg/errors"
)

// CurrentImports this is done with a global because subtemplates currently get called in functions. Lets aim to remove this eventually.
var CurrentImports *Imports

func Run(name string, tpldata interface{}) (*bytes.Buffer, error) {
	// load path relative to calling source file
	_, callerFile, _, _ := runtime.Caller(1)
	rootDir := filepath.Dir(callerFile)

	t := template.New("").Funcs(template.FuncMap{
		"ucFirst":       ucFirst,
		"lcFirst":       lcFirst,
		"quote":         strconv.Quote,
		"rawQuote":      rawQuote,
		"toCamel":       ToCamel,
		"dump":          dump,
		"prefixLines":   prefixLines,
		"reserveImport": CurrentImports.Reserve,
		"lookupImport":  CurrentImports.Lookup,
	})
	var roots []string

	for filename, data := range data {
		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			name = filepath.ToSlash(strings.TrimPrefix(path, rootDir+string(os.PathSeparator))) + filename
			t, err = t.New(name).Parse(data)
			if err != nil {
				panic(err)
			}
			roots = append(roots, name)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	buf := &bytes.Buffer{}
	err := t.Lookup(name).Execute(buf, tpldata)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func ucFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = unicode.ToUpper(r[0])
	return string(r)
}

func lcFirst(s string) string {
	if s == "" {
		return ""
	}

	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}

func isDelimiter(c rune) bool {
	return c == '-' || c == '_' || unicode.IsSpace(c)
}

func ToCamel(s string) string {
	buffer := make([]rune, 0, len(s))
	upper := true
	lastWasUpper := false

	for _, c := range s {
		if isDelimiter(c) {
			upper = true
			continue
		}
		if !lastWasUpper && unicode.IsUpper(c) {
			upper = true
		}

		if upper {
			buffer = append(buffer, unicode.ToUpper(c))
		} else {
			buffer = append(buffer, unicode.ToLower(c))
		}
		upper = false
		lastWasUpper = unicode.IsUpper(c)
	}

	return string(buffer)
}

func rawQuote(s string) string {
	return "`" + strings.Replace(s, "`", "`+\"`\"+`", -1) + "`"
}

func dump(val interface{}) string {
	switch val := val.(type) {
	case int:
		return strconv.Itoa(val)
	case int64:
		return fmt.Sprintf("%d", val)
	case float64:
		return fmt.Sprintf("%f", val)
	case string:
		return strconv.Quote(val)
	case bool:
		return strconv.FormatBool(val)
	case nil:
		return "nil"
	case []interface{}:
		var parts []string
		for _, part := range val {
			parts = append(parts, dump(part))
		}
		return "[]interface{}{" + strings.Join(parts, ",") + "}"
	case map[string]interface{}:
		buf := bytes.Buffer{}
		buf.WriteString("map[string]interface{}{")
		var keys []string
		for key := range val {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			data := val[key]

			buf.WriteString(strconv.Quote(key))
			buf.WriteString(":")
			buf.WriteString(dump(data))
			buf.WriteString(",")
		}
		buf.WriteString("}")
		return buf.String()
	default:
		panic(fmt.Errorf("unsupported type %T", val))
	}
}

func prefixLines(prefix, s string) string {
	return prefix + strings.Replace(s, "\n", "\n"+prefix, -1)
}

func RenderToFile(tpl string, filename string, data interface{}) error {
	if CurrentImports != nil {
		panic(fmt.Errorf("recursive or concurrent call to RenderToFile detected"))
	}
	CurrentImports = &Imports{destDir: filepath.Dir(filename)}
	// load path relative to calling source file
	_, callerFile, _, _ := runtime.Caller(1)
	rootDir := filepath.Dir(callerFile)
	t := template.New("").Funcs(template.FuncMap{
		"ucFirst":       ucFirst,
		"lcFirst":       lcFirst,
		"quote":         strconv.Quote,
		"rawQuote":      rawQuote,
		"toCamel":       ToCamel,
		"dump":          dump,
		"prefixLines":   prefixLines,
		"reserveImport": CurrentImports.Reserve,
		"lookupImport":  CurrentImports.Lookup,
	})
	var roots []string
	// load all the templates in the directory
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := filepath.ToSlash(strings.TrimPrefix(path, rootDir+string(os.PathSeparator)))
		if !strings.HasSuffix(info.Name(), tpl) {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		t, err = t.New(name).Parse(string(b))
		if err != nil {
			return errors.Wrap(err, filename)
		}

		roots = append(roots, name)
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "locating templates")
	}
	// then execute all the important looking ones in order, adding them to the same file
	sort.Slice(roots, func(i, j int) bool {
		// important files go first
		if strings.HasSuffix(roots[i], "!.gotpl") {
			return true
		}
		if strings.HasSuffix(roots[j], "!.gotpl") {
			return false
		}
		return roots[i] < roots[j]
	})

	var buf bytes.Buffer
	for _, root := range roots {
		err = t.Lookup(root).Execute(&buf, data)
		if err != nil {
			return errors.Wrap(err, root)
		}
	}

	b := bytes.Replace(buf.Bytes(), []byte("%%%IMPORTS%%%"), []byte(CurrentImports.String()), -1)
	CurrentImports = nil

	return write(filename, b)
}

func write(filename string, b []byte) error {
	err := os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create directory")
	}

	formatted, err := imports.Prune(filename, b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "gofmt failed on %s: %s\n", filepath.Base(filename), err.Error())
		formatted = b
	}

	err = ioutil.WriteFile(filename, formatted, 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write %s", filename)
	}

	return nil
}

var pkgReplacer = strings.NewReplacer(
	"/", "ᚋ",
	".", "ᚗ",
	"-", "ᚑ",
)

func TypeIdentifier(t types.Type) string {
	res := ""
	for {
		switch it := t.(type) {
		case *types.Pointer:
			t.Underlying()
			res += "ᚖ"
			t = it.Elem()
		case *types.Slice:
			res += "ᚕ"
			t = it.Elem()
		case *types.Named:
			res += pkgReplacer.Replace(it.Obj().Pkg().Path())
			res += "ᚐ"
			res += it.Obj().Name()
			return res
		case *types.Basic:
			res += it.Name()
			return res
		case *types.Map:
			res += "map"
			return res
		case *types.Interface:
			res += "interface"
			return res
		default:
			panic(fmt.Errorf("unexpected type %T", it))
		}
	}
}
