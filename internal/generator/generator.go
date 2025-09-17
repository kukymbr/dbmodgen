package generator

import (
	"context"

	"github.com/kukymbr/dbmodgen/internal/formatter"
	"github.com/kukymbr/dbmodgen/internal/generator/types"
	"github.com/kukymbr/dbmodgen/internal/util"
	"github.com/kukymbr/dbmodgen/internal/version"
)

func New(opt Options) (*Generator, error) {
	if err := prepareOptions(&opt); err != nil {
		return nil, err
	}

	f, err := formatter.Factory(opt.Formatter)
	if err != nil {
		return nil, err
	}

	util.PrintHellof("Hi, this is dbmodgen generator.")
	util.PrintDebugf("Options: " + opt.Debug())

	return &Generator{
		opt:       opt,
		formatter: f,
	}, nil
}

type Generator struct {
	opt       Options
	formatter formatter.Formatter
}

func (g *Generator) Generate(ctx context.Context) error {
	panic("not implemented")
}

func (g *Generator) newGenericData() types.GenericData {
	d := types.GenericData{
		Package:  g.opt.PackageName,
		FieldTag: g.opt.FieldTag,
		Version:  version.GetVersion(),
	}

	return d
}

func (g *Generator) format(ctx context.Context) {
	util.PrintDebugf("Formatting %s...", g.opt.TargetDir)

	if err := g.formatter.Format(ctx, g.opt.TargetDir); err != nil {
		util.PrintWarningf("Failed to format generated code: %s", err.Error())
	}
}
