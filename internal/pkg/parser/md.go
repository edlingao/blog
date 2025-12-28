package parser

import (
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/toc"
)

var MDParser = goldmark.New(
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	goldmark.WithExtensions(
		extension.GFM,
		extension.DefinitionList,
		extension.Table,
		extension.TaskList,
		extension.Linkify,
		meta.Meta,
		&toc.Extender{},
	),
)
