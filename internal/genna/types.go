package genna

import (
	"fmt"
	"html/template"

	"github.com/kukymbr/dbmodgen/internal/genna/model"
	"github.com/kukymbr/dbmodgen/internal/genna/util"
)

// Options for generator
type Options struct {
	// URL connection string
	URL string

	// Output file path
	Output string

	// List of Tables to generate
	// Default []string{"public.*"}
	Tables []string

	// Generate model for foreign keys,
	// even if Tables not listed in Tables param
	// will not generate fks if schema not listed
	FollowFKs bool

	// go-pg version
	GoPgVer int

	// Custom types goes here
	CustomTypes model.CustomTypeMapping

	// Package sets package name for model
	// Works only with SchemaPackage = false
	Package string

	// Do not replace primary key name to ID
	KeepPK bool

	// Soft delete column
	SoftDelete string

	// use sql.Null... instead of pointers
	UseSQLNulls bool

	// Do not generate alias tag
	NoAlias bool

	// Do not generate discard_unknown_columns tag
	NoDiscard bool

	// Override type for json/jsonb
	JSONTypes map[string]string

	// Add json tag to models
	AddJSONTag bool
}

// TemplatePackage stores package info
type TemplatePackage struct {
	Package string

	HasImports bool
	Imports    []string

	Entities []TemplateEntity
}

// NewTemplatePackage creates a package for template
func NewTemplatePackage(entities []model.Entity, options Options, tagName string) TemplatePackage {
	imports := util.NewSet()

	models := make([]TemplateEntity, len(entities))
	for i, entity := range entities {
		for _, imp := range entity.Imports {
			imports.Add(imp)
		}

		models[i] = NewTemplateEntity(entity, options, tagName)
	}

	return TemplatePackage{
		Package: options.Package,

		HasImports: imports.Len() > 0,
		Imports:    imports.Elements(),

		Entities: models,
	}
}

// TemplateEntity stores struct info
type TemplateEntity struct {
	model.Entity

	Tag template.HTML

	NoAlias bool
	Alias   string

	Columns []TemplateColumn

	HasRelations bool
	Relations    []TemplateRelation
}

// NewTemplateEntity creates an entity for template
func NewTemplateEntity(entity model.Entity, options Options, tagName string) TemplateEntity {
	if entity.HasMultiplePKs() {
		options.KeepPK = true
	}

	columns := make([]TemplateColumn, len(entity.Columns))
	for i, column := range entity.Columns {
		columns[i] = NewTemplateColumn(entity, column, options, tagName)
	}

	relations := make([]TemplateRelation, len(entity.Relations))
	for i, relation := range entity.Relations {
		relations[i] = NewTemplateRelation(relation, options, tagName)
	}

	tags := util.NewAnnotation()

	tags.AddTag(tagName, entity.PGFullName)

	return TemplateEntity{
		Entity: entity,
		Tag:    template.HTML(fmt.Sprintf("`%s`", tags.String())),

		NoAlias: options.NoAlias,
		Alias:   util.DefaultAlias,

		Columns: columns,

		HasRelations: len(relations) > 0,
		Relations:    relations,
	}
}

// TemplateColumn stores column info
type TemplateColumn struct {
	model.Column

	Tag     template.HTML
	Comment template.HTML
}

// NewTemplateColumn creates a column for template
func NewTemplateColumn(entity model.Entity, column model.Column, options Options, tagName string) TemplateColumn {
	if !options.KeepPK && column.IsPK {
		column.GoName = util.ID
	}

	if column.PGType == model.TypePGJSON || column.PGType == model.TypePGJSONB {
		if typ, ok := jsonType(options.JSONTypes, entity.PGSchema, entity.PGName, column.PGName); ok {
			column.Type = typ
		}
	}

	comment := ""
	tags := util.NewAnnotation()
	tags.AddTag(tagName, column.PGName)

	// ignore tag
	if column.GoType == model.TypeInterface {
		comment = "// unsupported"
		tags = util.NewAnnotation().AddTag(tagName, "-")
	}

	// add json tag
	if options.AddJSONTag {
		tags.AddTag("json", util.Underscore(column.PGName))
	}

	return TemplateColumn{
		Column: column,

		Tag:     template.HTML(fmt.Sprintf("`%s`", tags.String())),
		Comment: template.HTML(comment),
	}
}

// TemplateRelation stores relation info
type TemplateRelation struct {
	model.Relation

	Tag     template.HTML
	Comment template.HTML
}

// NewTemplateRelation creates relation for template
func NewTemplateRelation(relation model.Relation, options Options, tagName string) TemplateRelation {
	comment := ""
	tags := util.NewAnnotation()
	tags.AddTag(tagName, "-")

	if options.GoPgVer >= 10 {
		tags.AddTag("pg", "rel:has-one")
	}

	if len(relation.FKFields) > 1 {
		comment = "// unsupported"

	}

	// add json tag
	if options.AddJSONTag {
		tags.AddTag("json", util.Underscore(relation.GoName))
	}

	return TemplateRelation{
		Relation: relation,

		Tag:     template.HTML(fmt.Sprintf("`%s`", tags.String())),
		Comment: template.HTML(comment),
	}
}

func jsonType(mp map[string]string, schema, table, field string) (string, bool) {
	if mp == nil {
		return "", false
	}

	patterns := [][3]string{
		{schema, table, field},
		{schema, "*", field},
		{schema, table, "*"},
		{schema, "*", "*"},
	}

	var names []string
	for _, parts := range patterns {
		names = append(names, fmt.Sprintf("%s.%s", util.Join(parts[0], parts[1]), parts[2]))
		names = append(names, fmt.Sprintf("%s.%s", util.JoinF(parts[0], parts[1]), parts[2]))
	}
	names = append(names, util.Join(schema, table), "*")

	for _, name := range names {
		if v, ok := mp[name]; ok {
			return v, true
		}
	}

	return "", false
}

func tagName(options Options) string {
	if options.GoPgVer == 8 {
		return "sql"
	}
	return "pg"
}
