package utils

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// markdown parser instance with common extensions
var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM, // GitHub Flavored Markdown
		extension.Footnote,
		extension.Typographer,
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
		html.WithUnsafe(), // 允许原始HTML
	),
)

// RenderMarkdown 将Markdown文本渲染为HTML
func RenderMarkdown(source string) (string, error) {
	if source == "" {
		return "", nil
	}

	var buf bytes.Buffer
	if err := md.Convert([]byte(source), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
} 