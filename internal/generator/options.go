package generator

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kukymbr/dbmodgen/internal/formatter"
	gennautil "github.com/kukymbr/dbmodgen/internal/genna/util"
	"github.com/kukymbr/dbmodgen/internal/util"
	"gopkg.in/yaml.v3"
)

const (
	DefaultPackageName = "dbmodel"
	DefaultTargetFile  = "internal/" + DefaultPackageName + "/model.gen.go"
	DefaultFieldTag    = "db"

	EnvDSN = "DBMODGEN_DSN"
)

func ReadOptions(path string, getEnv func(string) string) (Options, error) {
	dsn := getEnv(EnvDSN)
	if dsn == "" {
		return Options{}, fmt.Errorf("no DBMODGEN_DSN environment variable is provided")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return Options{}, fmt.Errorf("failed to read options file: %w", err)
	}

	var opt Options

	if err := yaml.Unmarshal(data, &opt); err != nil {
		return Options{}, fmt.Errorf("failed to parse options file: %w", err)
	}

	opt.DSN = dsn

	return opt, nil
}

type Options struct {
	// DSN is a database connection URL.
	// Is set from the DBMODGEN_DSN environment variable.
	DSN string `json:"-" yaml:"-"`

	// TargetFile is a path to a target generated code file.
	// Default is "internal/dbmodel/model.gen.go".
	TargetFile string `json:"target_file" yaml:"target_file"`

	// PackageName is a target package name of the generated code.
	// Default is "dbmodel".
	PackageName string `json:"package_name" yaml:"package_name"`

	// Tables is a list of the table names to generate models for.
	Tables []string `json:"tables" yaml:"tables"`

	// FieldTag is a name tag, added to the model fields.
	// Default is "db".
	FieldTag string `json:"field_tag" yaml:"field_tag"`

	// Formatter is a name of the formatter for the generated code files.
	// Available options: gofmt (default), none.
	Formatter string `json:"formatter" yaml:"formatter"`

	// UseSQLNulls enables sql.Null types instead of pointers.
	UseSQLNulls bool `json:"use_sql_nulls" yaml:"use_sql_nulls"`
}

func (opt Options) Debug() string {
	encoded, err := yaml.Marshal(opt)
	if err != nil {
		return fmt.Sprintf("%#v", opt)
	}

	return string(encoded)
}

func prepareOptions(opt *Options) error {
	if opt.TargetFile == "" {
		opt.TargetFile = DefaultTargetFile
	}

	if opt.PackageName == "" {
		opt.PackageName = filepath.Base(filepath.Dir(opt.TargetFile))
	}

	if opt.PackageName == "" {
		opt.PackageName = DefaultPackageName
	}

	if opt.FieldTag == "" {
		opt.FieldTag = DefaultFieldTag
	}

	if opt.Formatter == "" {
		opt.Formatter = formatter.DefaultFormatter
	}

	if len(opt.Tables) == 0 {
		opt.Tables = []string{gennautil.PublicSchema + ".*"}
	}

	if err := util.ValidatePackageName(opt.PackageName); err != nil {
		return err
	}

	if err := util.ValidateTag(opt.FieldTag); err != nil {
		return err
	}

	if err := util.EnsureDir(filepath.Dir(opt.TargetFile)); err != nil {
		return err
	}

	return nil
}
