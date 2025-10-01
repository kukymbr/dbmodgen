package generator

import (
	"bytes"
	"context"
	"fmt"

	"github.com/kukymbr/dbmodgen/internal/formatter"
	"github.com/kukymbr/dbmodgen/internal/generator/templates"
	"github.com/kukymbr/dbmodgen/internal/genna"
	"github.com/kukymbr/dbmodgen/internal/genna/model"
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
		genna:     genna.New(opt.DSN, newGennaLogger()),
	}, nil
}

type Generator struct {
	opt       Options
	formatter formatter.Formatter

	genna *genna.Genna
}

func (g *Generator) Generate(ctx context.Context) error {
	entities, err := g.genna.Read(g.opt.Tables, false, g.opt.UseSQLNulls, model.CustomTypeMapping{} /* TODO */)
	if err != nil {
		return fmt.Errorf("failed to read entities from database: %w", err)
	}

	tplData := modelTplData{
		TemplatePackage: genna.NewTemplatePackage(entities, genna.Options{
			URL:     g.opt.DSN,
			Output:  g.opt.TargetFile,
			Package: g.opt.PackageName,
			Tables:  g.opt.Tables,

			UseSQLNulls: g.opt.UseSQLNulls,
			CustomTypes: model.CustomTypeMapping{}, // TODO
			NoAlias:     true,
			NoDiscard:   true,
			AddJSONTag:  true,
			FollowFKs:   true,
			KeepPK:      true,
			JSONTypes:   nil, // TODO
		}, g.opt.FieldTag),

		Version: version.GetVersion(),
	}

	var buffer bytes.Buffer

	if err := templates.ExecuteModelTemplate(&buffer, tplData); err != nil {
		return err
	}

	content := g.format(ctx, buffer.Bytes())

	if err := util.WriteFile(g.opt.TargetFile, content); err != nil {
		return fmt.Errorf("failed to write target file %s: %w", g.opt.TargetFile, err)
	}

	return nil
}

func (g *Generator) format(ctx context.Context, content []byte) []byte {
	util.PrintDebugf("Formatting generated code...")

	formatted, err := g.formatter.Format(ctx, content)
	if err != nil {
		util.PrintWarningf("Failed to format generated code: %s", err.Error())

		return content
	}

	return formatted
}
