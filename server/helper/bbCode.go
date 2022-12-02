package helper

import (
	"github.com/frustra/bbcode"
	"strconv"
)

func NewBBCodeCompiler() bbcode.Compiler {
	// set up compiler
	bbcodeCompiler := bbcode.NewCompiler(true, true)
	bbcodeCompiler.SetTag("size", nil)
	bbcodeCompiler.SetTag("color", nil)
	bbcodeCompiler.SetTag("center", nil)

	for i := 1; i <= 6; i++ {
		bbCodeToHTMLSameName(bbcodeCompiler, "h"+strconv.Itoa(i))
	}

	bbCodeToHTMLSameName(bbcodeCompiler, "u")
	bbCodeToHTMLSameName(bbcodeCompiler, "i")
	bbCodeToHTMLSameName(bbcodeCompiler, "hr")

	bbcodeCompiler.SetTag("strike", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "s"
		return out, true
	})
	bbcodeCompiler.SetTag("list", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "ul"
		return out, true
	})
	bbcodeCompiler.SetTag("olist", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "ol"
		return out, true
	})

	bbcodeCompiler.SetTag("spoiler", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "span"
		out.Attrs = map[string]string{
			"style": "background-color: black",
		}
		return out, true
	})
	bbcodeCompiler.SetTag("*", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = "li"
		return out, true
	})

	return bbcodeCompiler
}

func NewBBCodeToTextCompiler() bbcode.Compiler {
	bbcodeCompiler := bbcode.NewCompiler(true, true)
	bbcodeCompiler.SetTag("size", nil)
	bbcodeCompiler.SetTag("color", nil)
	bbcodeCompiler.SetTag("center", nil)
	removeTagComplete(bbcodeCompiler, "img")
	removeTagComplete(bbcodeCompiler, "url")

	for i := 1; i <= 6; i++ {
		removeTagComplete(bbcodeCompiler, "h"+strconv.Itoa(i))
	}

	removeTag(bbcodeCompiler, "strike")
	removeTag(bbcodeCompiler, "list")
	removeTag(bbcodeCompiler, "olist")
	removeTag(bbcodeCompiler, "spoiler")
	removeTag(bbcodeCompiler, "*")

	return bbcodeCompiler
}

func bbCodeToHTMLSameName(c bbcode.Compiler, name string) {
	c.SetTag(name, func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out := bbcode.NewHTMLTag("")
		out.Name = name
		return out, true
	})
}

func removeTag(c bbcode.Compiler, name string) {
	c.SetTag(name, func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		return bbcode.NewHTMLTag(""), true
	})
}

func removeTagComplete(c bbcode.Compiler, name string) {
	c.SetTag(name, func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		return nil, false
	})
}
