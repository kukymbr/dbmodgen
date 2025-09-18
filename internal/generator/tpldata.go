package generator

import "github.com/kukymbr/dbmodgen/internal/genna"

type modelTplData struct {
	genna.TemplatePackage

	Version string
}
