package generator

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	defaultPrimaryFieldName    = "ID"
	defaultSoftDeleteFieldName = "DeletedAt"
)

type (
	// Generator generates the CRUD methods.
	Generator struct {
		pkg                   *types.Package
		t                     *types.Struct
		structModel           string
		header, body          bytes.Buffer // Accumulated output.
		existingTypes         []cruderType
		TableName             string
		PkgName               string
		readFields            map[int]string
		writeFields           map[int]string
		primaryFieldOffset    int
		softDeleteFieldOffset int

		mx      sync.Mutex
		imports map[string]bool
	}
)

// New returns a Generator
func New(pkg *types.Package, structModel string) (*Generator, error) {
	o := pkg.Scope().Lookup(structModel)
	if o == nil {
		return nil, fmt.Errorf("the struct %s doesn't seem to exists in package %s", structModel, pkg.Name())
	}
	t, ok := o.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, fmt.Errorf("the type %s is not a struct", structModel)
	}

	gen := &Generator{
		pkg:                   pkg,
		t:                     t,
		structModel:           structModel,
		TableName:             structModel,
		PkgName:               pkg.Name(),
		readFields:            make(map[int]string, t.NumFields()),
		writeFields:           make(map[int]string, t.NumFields()),
		softDeleteFieldOffset: -1, // -1 disable soft deletion
		imports:               make(map[string]bool),
	}

	for i := 0; i < gen.t.NumFields(); i++ {
		// If the defaultSoftDeleteFieldName exists in the struct, use it
		// for soft deletion. Also don't include it in readFields or writeFields
		if defaultSoftDeleteFieldName == gen.t.Field(i).Name() {
			gen.softDeleteFieldOffset = i
			continue
		}

		gen.readFields[i] = gen.t.Field(i).Name()

		// If the defaultPrimaryFieldName exists in the struct, use it
		// as primary key. Also don't include it in writeFields
		if defaultPrimaryFieldName == gen.t.Field(i).Name() {
			gen.primaryFieldOffset = i
			continue
		}
		gen.writeFields[i] = gen.t.Field(i).Name()
	}

	gen.addImport("database/sql") // All methods use this package

	return gen, nil
}

// SetReadFields sets the fields that should be returned in reading operations.
// The passed in slice will be match against the fieldnames of the struct
func (g *Generator) SetReadFields(fields []string) error {
	var fieldNames map[int]string
	for _, f := range fields {
		for i := 0; i < g.t.NumFields(); i++ {
			if strings.TrimSpace(f) == g.t.Field(i).Name() {
				fieldNames[i] = g.t.Field(i).Name()
				continue
			}

			return fmt.Errorf("the field %s does not exists in struct %s", f, g.structModel)
		}
	}

	g.readFields = fieldNames
	return nil
}

// SetWriteFields sets the fields that should be returned in writing operations.
// The passed in slice will be match against the fieldnames of the struct
func (g *Generator) SetWriteFields(fields []string) error {
	var fieldNames map[int]string
	for _, f := range fields {
		for i := 0; i < g.t.NumFields(); i++ {
			if strings.TrimSpace(f) == g.t.Field(i).Name() {
				fieldNames[i] = g.t.Field(i).Name()
				continue
			}

			return fmt.Errorf("the field %s does not exists in struct %s", f, g.structModel)
		}
	}

	g.writeFields = fieldNames
	return nil
}

// SetPrimaryField sets the field that should be used for soft deletion.
// The field should be of type nullable datetime but this function does not check that.
func (g *Generator) SetPrimaryField(f string) error {
	for i := 0; i < g.t.NumFields(); i++ {
		if strings.TrimSpace(f) == g.t.Field(i).Name() {
			g.primaryFieldOffset = i

			return nil
		}

	}

	return fmt.Errorf("the field %s does not exists in struct %s", f, g.structModel)
}

// SetSoftDeleteField sets the field that should be used for soft deletion.
// The field should be of type nullable datetime but this function does not check that.
func (g *Generator) SetSoftDeleteField(f string) error {
	for i := 0; i < g.t.NumFields(); i++ {
		if strings.TrimSpace(f) == g.t.Field(i).Name() {
			g.softDeleteFieldOffset = i

			return nil
		}

	}

	return fmt.Errorf("the field %s does not exists in struct %s", f, g.structModel)
}

// readFieldNames returns a slice of the read field names.
// A prefix can be passed which would be added before each name.
func (g *Generator) readFieldNames(prefix string) []string {
	var fieldNames []string
	// Ordered iteration of the map
	var keys []int
	for k := range g.readFields {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fieldNames = append(fieldNames, prefix+g.readFields[k])
	}

	return fieldNames
}

// readFieldDBNames returns a slice of the read field names, but in their DB forms (if any).
// The DB form is taken from the "db" struct tag if defined, or it will be the same as the field
// name. A prefix can be passed which would be added before each name.
func (g *Generator) readFieldDBNames(prefix string) []string {
	var fieldNames []string
	// Ordered iteration of the map
	var keys []int
	for k := range g.readFields {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fieldNames = append(fieldNames, prefix+g.fieldDBName(k))
	}

	return fieldNames
}

// writeFieldNames returns a slice of the read field names.
// A prefix can be passed which would be added before each name.
func (g *Generator) writeFieldNames(prefix string) []string {
	var fieldNames []string
	// Ordered iteration of the map
	var keys []int
	for k := range g.writeFields {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fieldNames = append(fieldNames, prefix+g.writeFields[k])
	}

	return fieldNames
}

// writeFieldDBNames returns a slice of the read field names, but in their DB forms (if any).
// The DB form is taken from the "db" struct tag if defined, or it will be the same as the field
// name. A prefix can be passed which would be added before each name.
func (g *Generator) writeFieldDBNames(prefix string) []string {
	var fieldNames []string
	// Ordered iteration of the map
	var keys []int
	for k := range g.writeFields {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		fieldNames = append(fieldNames, prefix+g.fieldDBName(k))
	}

	return fieldNames
}

// placeholderStrings returns a slice of strings for n placeholders
func (g *Generator) placeholderStrings(n int) []string {
	var placeholders []string
	for i := 1; i <= n; i++ {
		placeholders = append(placeholders, "$"+strconv.Itoa(i))
	}

	return placeholders
}

// HeaderPrintf writes the input to the Generator's header buffer
func (g *Generator) HeaderPrintf(in string, args ...interface{}) {
	fmt.Fprintf(&g.header, in, args...)
}

// Printf writes the input to the Generator's body buffer
func (g *Generator) Printf(in string, args ...interface{}) {
	fmt.Fprintf(&g.body, in, args...)
}

func (g *Generator) pkgDecl() []byte {
	return []byte(fmt.Sprintf("package %s\n", g.PkgName))
}

func (g *Generator) importsDecl() []byte {
	var imports []string
	for i := range g.imports {
		imports = append(imports, "\""+i+"\"")
	}

	return []byte(fmt.Sprintf(`
import (
	%s
)
`, strings.Join(imports, "\n\t")))
}

// Format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) Format() ([]byte, error) {
	return format.Source(append(append(g.pkgDecl(), g.importsDecl()...), append(g.header.Bytes(), g.body.Bytes()...)...))
}

// String output all buffers as string
func (g *Generator) String() string {
	return string(g.pkgDecl()) + string(g.importsDecl()) + g.header.String() + g.body.String()
}

func (g *Generator) typeExist(t cruderType) bool {
	for _, x := range g.existingTypes {
		if x == t {
			return true
		}
	}

	// Check if type already exist in package
	return g.pkg.Scope().Lookup(string(t)) != nil
}

func (g *Generator) addImport(pkg string) {
	g.mx.Lock()
	g.imports[pkg] = true
	g.mx.Unlock()
}

// Return the name of the field in their DB form. The DB form is taken from the
// "db" struct tag if defined, otherwise the field name.
func (g *Generator) fieldDBName(i int) string {
	if i >= g.t.NumFields() {
		return ""
	}

	st := reflect.StructTag(g.t.Tag(i))
	if st.Get("db") != "" {
		return st.Get("db")
	}

	return g.t.Field(i).Name()
}
