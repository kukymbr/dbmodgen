package generator

import (
	"strings"

	"github.com/kukymbr/dbmodgen/internal/formatter"
	"github.com/kukymbr/dbmodgen/internal/util"
)

const (
	DefaultPackageName = "dbmodel"
	DefaultTargetDir   = "internal/" + DefaultPackageName
	DefaultFieldTag    = "db"
)

type Options struct {
	// PackageName is a target package name of the generated code.
	// Default is "queries".
	PackageName string

	// TargetDir is a target go code directory.
	// Default is "internal/queries".
	TargetDir string

	// FieldTag is a name tag, added to the model fields.
	// Default is "db".
	FieldTag string

	// Formatter is a name of the formatter for the generated code files.
	// Available options: gofmt (default), none.
	Formatter string
}

func (opt Options) Debug() string {
	values := []string{
		"package_name" + "=" + opt.PackageName,
		"target_dir" + "=" + opt.TargetDir,
		"field_tag" + "=" + opt.FieldTag,
	}

	return strings.Join(values, "; ")
}

func prepareOptions(opt *Options) error {
	if opt.PackageName == "" {
		opt.PackageName = DefaultPackageName
	}

	if opt.TargetDir == "" {
		opt.TargetDir = DefaultTargetDir
	}

	if opt.FieldTag == "" {
		opt.FieldTag = DefaultFieldTag
	}

	if opt.Formatter == "" {
		opt.Formatter = formatter.DefaultFormatter
	}

	if err := util.ValidatePackageName(opt.PackageName); err != nil {
		return err
	}

	if err := util.ValidateTag(opt.FieldTag); err != nil {
		return err
	}

	if err := util.EnsureDir(opt.TargetDir); err != nil {
		return err
	}

	return nil
}
