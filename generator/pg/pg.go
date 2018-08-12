package pg

import (
	"bytes"
	"fmt"
	"go/format"
	"go/types"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/pengux/cruder/generator"
)

const (
	defaultPrimaryFieldName    = "ID"
	defaultSoftDeleteFieldName = "DeletedAt"
)

type (
	// PG generates the CRUD methods for Postgresql using lib/pg.
	PG struct {
		pkg                   *types.Package
		t                     *types.Struct
		structModel           string
		header, body          bytes.Buffer // Accumulated output.
		existingTypes         []cruderType
		TableName             string
		PkgName               string
		SkipSuffix            bool
		readFields            map[int]string
		writeFields           map[int]string
		primaryFieldOffset    int
		softDeleteFieldOffset int
		sqlImportAdded        bool

		mx      sync.Mutex
		imports map[string]bool
	}
)

// New returns a PG
func New(pkg *types.Package, t *types.Struct, structModel string) (*PG, error) {
	gen := &PG{
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

	return gen, nil
}

// Generate generates CRUD code for the passed in functions
func (g *PG) Generate(w io.Writer, fns ...generator.Function) error {
	return nil
}

// SetReadFields sets the fields that should be returned in reading operations.
// The passed in slice will be match against the fieldnames of the struct
func (g *PG) SetReadFields(fields []string) error {
	fieldNames := make(map[int]string)
	for _, f := range fields {
		found := false
		for i := 0; i < g.t.NumFields(); i++ {
			if strings.TrimSpace(f) == g.t.Field(i).Name() {
				fieldNames[i] = g.t.Field(i).Name()
				found = true
				break
			}

		}

		if !found {
			return fmt.Errorf("the field %s does not exists in struct %s", f, g.structModel)
		}
	}

	g.readFields = fieldNames
	return nil
}

// SetWriteFields sets the fields that should be returned in writing operations.
// The passed in slice will be match against the fieldnames of the struct
func (g *PG) SetWriteFields(fields []string) error {
	fieldNames := make(map[int]string)
	for _, f := range fields {
		found := false
		for i := 0; i < g.t.NumFields(); i++ {
			if strings.TrimSpace(f) == g.t.Field(i).Name() {
				fieldNames[i] = g.t.Field(i).Name()
				found = true
				break
			}

		}
		if !found {
			return fmt.Errorf("the field %s does not exists in struct %s", f, g.structModel)
		}
	}

	g.writeFields = fieldNames
	return nil
}

// SetPrimaryField sets the field that are used as primary key in lookups
func (g *PG) SetPrimaryField(f string) error {
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
func (g *PG) SetSoftDeleteField(f string) error {
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
func (g *PG) readFieldNames(prefix string) []string {
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
func (g *PG) readFieldDBNames(prefix string) []string {
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

// writeFieldNames returns a slice of the write field names.
// A prefix can be passed which would be added before each name.
func (g *PG) writeFieldNames(prefix string) []string {
	var fieldNames []string
	// Ordered iteration of the map
	var keys []int
	for k := range g.writeFields {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, k := range keys {
		name := prefix + g.writeFields[k]

		// If the field is a struct, then return a pointer expression
		if _, ok := g.t.Field(k).Type().(*types.Named); ok {
			name = "&" + name
		}
		fieldNames = append(fieldNames, name)
	}

	return fieldNames
}

// writeFieldDBNames returns a slice of the read field names, but in their DB forms (if any).
// The DB form is taken from the "db" struct tag if defined, or it will be the same as the field
// name. A prefix can be passed which would be added before each name.
func (g *PG) writeFieldDBNames(prefix string) []string {
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
func (g *PG) placeholderStrings(n int) []string {
	var placeholders []string
	for i := 1; i <= n; i++ {
		placeholders = append(placeholders, "$"+strconv.Itoa(i))
	}

	return placeholders
}

// HeaderPrintf writes the input to the PG's header buffer
func (g *PG) HeaderPrintf(in string, args ...interface{}) {
	fmt.Fprintf(&g.header, in, args...)
}

// Printf writes the input to the PG's body buffer
func (g *PG) Printf(in string, args ...interface{}) {
	fmt.Fprintf(&g.body, in, args...)
}

func (g *PG) pkgDecl() []byte {
	return []byte(fmt.Sprintf("package %s\n", g.PkgName))
}

func (g *PG) importsDecl() []byte {
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

// Format returns the gofmt-ed contents of the PG's buffer.
func (g *PG) Format() ([]byte, error) {
	return format.Source(append(append(g.pkgDecl(), g.importsDecl()...), append(g.header.Bytes(), g.body.Bytes()...)...))
}

// String output all buffers as string
func (g *PG) String() string {
	return string(g.pkgDecl()) + string(g.importsDecl()) + g.header.String() + g.body.String()
}

func (g *PG) typeExist(t cruderType) bool {
	for _, x := range g.existingTypes {
		if x == t {
			return true
		}
	}

	// Check if type already exist in package
	return g.pkg.Scope().Lookup(string(t)) != nil
}

func (g *PG) addImport(pkg string) {
	g.mx.Lock()
	g.imports[pkg] = true
	g.mx.Unlock()
}

// Return the name of the field in their DB form. The DB form is taken from the
// "db" struct tag if defined, otherwise the field name.
func (g *PG) fieldDBName(i int) string {
	if i >= g.t.NumFields() {
		return ""
	}

	st := reflect.StructTag(g.t.Tag(i))
	if st.Get("db") != "" {
		return st.Get("db")
	}

	return g.t.Field(i).Name()
}
