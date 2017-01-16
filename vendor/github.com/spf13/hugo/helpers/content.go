// Copyright 2015 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package helpers implements general utility functions that work with
// and on content.  The helper functions defined here lay down the
// foundation of how Hugo works with files and filepaths, and perform
// string operations on content.
package helpers

import (
	"bytes"
	"html/template"
	"os/exec"
	"unicode"
	"unicode/utf8"

	"github.com/miekg/mmark"
	"github.com/mitchellh/mapstructure"
	"github.com/russross/blackfriday"
	bp "github.com/spf13/hugo/bufferpool"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"

	"strings"
	"sync"
)

// SummaryLength is the length of the summary that Hugo extracts from a content.
var SummaryLength = 70

// SummaryDivider denotes where content summarization should end. The default is "<!--more-->".
var SummaryDivider = []byte("<!--more-->")

// Blackfriday holds configuration values for Blackfriday rendering.
type Blackfriday struct {
	Smartypants                      bool
	AngledQuotes                     bool
	Fractions                        bool
	HrefTargetBlank                  bool
	SmartDashes                      bool
	LatexDashes                      bool
	TaskLists                        bool
	PlainIDAnchors                   bool
	SourceRelativeLinksEval          bool
	SourceRelativeLinksProjectFolder string
	Extensions                       []string
	ExtensionsMask                   []string
}

// NewBlackfriday creates a new Blackfriday filled with site config or some sane defaults.
func NewBlackfriday(c ConfigProvider) *Blackfriday {

	defaultParam := map[string]interface{}{
		"smartypants":                      true,
		"angledQuotes":                     false,
		"fractions":                        true,
		"hrefTargetBlank":                  false,
		"smartDashes":                      true,
		"latexDashes":                      true,
		"plainIDAnchors":                   true,
		"taskLists":                        true,
		"sourceRelativeLinks":              false,
		"sourceRelativeLinksProjectFolder": "/docs/content",
	}

	ToLowerMap(defaultParam)

	siteParam := c.GetStringMap("blackfriday")

	siteConfig := make(map[string]interface{})

	for k, v := range defaultParam {
		siteConfig[k] = v
	}

	if siteParam != nil {
		for k, v := range siteParam {
			siteConfig[k] = v
		}
	}

	combinedConfig := &Blackfriday{}
	if err := mapstructure.Decode(siteConfig, combinedConfig); err != nil {
		jww.FATAL.Printf("Failed to get site rendering config\n%s", err.Error())
	}

	return combinedConfig
}

var blackfridayExtensionMap = map[string]int{
	"noIntraEmphasis":        blackfriday.EXTENSION_NO_INTRA_EMPHASIS,
	"tables":                 blackfriday.EXTENSION_TABLES,
	"fencedCode":             blackfriday.EXTENSION_FENCED_CODE,
	"autolink":               blackfriday.EXTENSION_AUTOLINK,
	"strikethrough":          blackfriday.EXTENSION_STRIKETHROUGH,
	"laxHtmlBlocks":          blackfriday.EXTENSION_LAX_HTML_BLOCKS,
	"spaceHeaders":           blackfriday.EXTENSION_SPACE_HEADERS,
	"hardLineBreak":          blackfriday.EXTENSION_HARD_LINE_BREAK,
	"tabSizeEight":           blackfriday.EXTENSION_TAB_SIZE_EIGHT,
	"footnotes":              blackfriday.EXTENSION_FOOTNOTES,
	"noEmptyLineBeforeBlock": blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK,
	"headerIds":              blackfriday.EXTENSION_HEADER_IDS,
	"titleblock":             blackfriday.EXTENSION_TITLEBLOCK,
	"autoHeaderIds":          blackfriday.EXTENSION_AUTO_HEADER_IDS,
	"backslashLineBreak":     blackfriday.EXTENSION_BACKSLASH_LINE_BREAK,
	"definitionLists":        blackfriday.EXTENSION_DEFINITION_LISTS,
}

var stripHTMLReplacer = strings.NewReplacer("\n", " ", "</p>", "\n", "<br>", "\n", "<br />", "\n")

var mmarkExtensionMap = map[string]int{
	"tables":                 mmark.EXTENSION_TABLES,
	"fencedCode":             mmark.EXTENSION_FENCED_CODE,
	"autolink":               mmark.EXTENSION_AUTOLINK,
	"laxHtmlBlocks":          mmark.EXTENSION_LAX_HTML_BLOCKS,
	"spaceHeaders":           mmark.EXTENSION_SPACE_HEADERS,
	"hardLineBreak":          mmark.EXTENSION_HARD_LINE_BREAK,
	"footnotes":              mmark.EXTENSION_FOOTNOTES,
	"noEmptyLineBeforeBlock": mmark.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK,
	"headerIds":              mmark.EXTENSION_HEADER_IDS,
	"autoHeaderIds":          mmark.EXTENSION_AUTO_HEADER_IDS,
}

// StripHTML accepts a string, strips out all HTML tags and returns it.
func StripHTML(s string) string {

	// Shortcut strings with no tags in them
	if !strings.ContainsAny(s, "<>") {
		return s
	}
	s = stripHTMLReplacer.Replace(s)

	// Walk through the string removing all tags
	b := bp.GetBuffer()
	defer bp.PutBuffer(b)
	var inTag, isSpace, wasSpace bool
	for _, r := range s {
		if !inTag {
			isSpace = false
		}

		switch {
		case r == '<':
			inTag = true
		case r == '>':
			inTag = false
		case unicode.IsSpace(r):
			isSpace = true
			fallthrough
		default:
			if !inTag && (!isSpace || (isSpace && !wasSpace)) {
				b.WriteRune(r)
			}
		}

		wasSpace = isSpace

	}
	return b.String()
}

// stripEmptyNav strips out empty <nav> tags from content.
func stripEmptyNav(in []byte) []byte {
	return bytes.Replace(in, []byte("<nav>\n</nav>\n\n"), []byte(``), -1)
}

// BytesToHTML converts bytes to type template.HTML.
func BytesToHTML(b []byte) template.HTML {
	return template.HTML(string(b))
}

// getHTMLRenderer creates a new Blackfriday HTML Renderer with the given configuration.
func getHTMLRenderer(defaultFlags int, ctx *RenderingContext) blackfriday.Renderer {
	renderParameters := blackfriday.HtmlRendererParameters{
		FootnoteAnchorPrefix:       viper.GetString("footnoteAnchorPrefix"),
		FootnoteReturnLinkContents: viper.GetString("footnoteReturnLinkContents"),
	}

	b := len(ctx.DocumentID) != 0

	if b && !ctx.getConfig().PlainIDAnchors {
		renderParameters.FootnoteAnchorPrefix = ctx.DocumentID + ":" + renderParameters.FootnoteAnchorPrefix
		renderParameters.HeaderIDSuffix = ":" + ctx.DocumentID
	}

	htmlFlags := defaultFlags
	htmlFlags |= blackfriday.HTML_USE_XHTML
	htmlFlags |= blackfriday.HTML_FOOTNOTE_RETURN_LINKS

	if ctx.getConfig().Smartypants {
		htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	}

	if ctx.getConfig().AngledQuotes {
		htmlFlags |= blackfriday.HTML_SMARTYPANTS_ANGLED_QUOTES
	}

	if ctx.getConfig().Fractions {
		htmlFlags |= blackfriday.HTML_SMARTYPANTS_FRACTIONS
	}

	if ctx.getConfig().HrefTargetBlank {
		htmlFlags |= blackfriday.HTML_HREF_TARGET_BLANK
	}

	if ctx.getConfig().SmartDashes {
		htmlFlags |= blackfriday.HTML_SMARTYPANTS_DASHES
	}

	if ctx.getConfig().LatexDashes {
		htmlFlags |= blackfriday.HTML_SMARTYPANTS_LATEX_DASHES
	}

	return &HugoHTMLRenderer{
		RenderingContext: ctx,
		Renderer:         blackfriday.HtmlRendererWithParameters(htmlFlags, "", "", renderParameters),
	}
}

func getMarkdownExtensions(ctx *RenderingContext) int {
	// Default Blackfriday common extensions
	commonExtensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS

	// Extra Blackfriday extensions that Hugo enables by default
	flags := commonExtensions |
		blackfriday.EXTENSION_AUTO_HEADER_IDS |
		blackfriday.EXTENSION_FOOTNOTES

	for _, extension := range ctx.getConfig().Extensions {
		if flag, ok := blackfridayExtensionMap[extension]; ok {
			flags |= flag
		}
	}
	for _, extension := range ctx.getConfig().ExtensionsMask {
		if flag, ok := blackfridayExtensionMap[extension]; ok {
			flags &= ^flag
		}
	}
	return flags
}

func markdownRender(ctx *RenderingContext) []byte {
	if ctx.RenderTOC {
		return blackfriday.Markdown(ctx.Content,
			getHTMLRenderer(blackfriday.HTML_TOC, ctx),
			getMarkdownExtensions(ctx))
	}
	return blackfriday.Markdown(ctx.Content, getHTMLRenderer(0, ctx),
		getMarkdownExtensions(ctx))
}

// getMmarkHTMLRenderer creates a new mmark HTML Renderer with the given configuration.
func getMmarkHTMLRenderer(defaultFlags int, ctx *RenderingContext) mmark.Renderer {
	renderParameters := mmark.HtmlRendererParameters{
		FootnoteAnchorPrefix:       viper.GetString("footnoteAnchorPrefix"),
		FootnoteReturnLinkContents: viper.GetString("footnoteReturnLinkContents"),
	}

	b := len(ctx.DocumentID) != 0

	if b && !ctx.getConfig().PlainIDAnchors {
		renderParameters.FootnoteAnchorPrefix = ctx.DocumentID + ":" + renderParameters.FootnoteAnchorPrefix
		// renderParameters.HeaderIDSuffix = ":" + ctx.DocumentId
	}

	htmlFlags := defaultFlags
	htmlFlags |= mmark.HTML_FOOTNOTE_RETURN_LINKS

	return &HugoMmarkHTMLRenderer{
		mmark.HtmlRendererWithParameters(htmlFlags, "", "", renderParameters),
	}
}

func getMmarkExtensions(ctx *RenderingContext) int {
	flags := 0
	flags |= mmark.EXTENSION_TABLES
	flags |= mmark.EXTENSION_FENCED_CODE
	flags |= mmark.EXTENSION_AUTOLINK
	flags |= mmark.EXTENSION_SPACE_HEADERS
	flags |= mmark.EXTENSION_CITATION
	flags |= mmark.EXTENSION_TITLEBLOCK_TOML
	flags |= mmark.EXTENSION_HEADER_IDS
	flags |= mmark.EXTENSION_AUTO_HEADER_IDS
	flags |= mmark.EXTENSION_UNIQUE_HEADER_IDS
	flags |= mmark.EXTENSION_FOOTNOTES
	flags |= mmark.EXTENSION_SHORT_REF
	flags |= mmark.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK
	flags |= mmark.EXTENSION_INCLUDE

	for _, extension := range ctx.getConfig().Extensions {
		if flag, ok := mmarkExtensionMap[extension]; ok {
			flags |= flag
		}
	}
	return flags
}

func mmarkRender(ctx *RenderingContext) []byte {
	return mmark.Parse(ctx.Content, getMmarkHTMLRenderer(0, ctx),
		getMmarkExtensions(ctx)).Bytes()
}

// ExtractTOC extracts Table of Contents from content.
func ExtractTOC(content []byte) (newcontent []byte, toc []byte) {
	origContent := make([]byte, len(content))
	copy(origContent, content)
	first := []byte(`<nav>
<ul>`)

	last := []byte(`</ul>
</nav>`)

	replacement := []byte(`<nav id="TableOfContents">
<ul>`)

	startOfTOC := bytes.Index(content, first)

	peekEnd := len(content)
	if peekEnd > 70+startOfTOC {
		peekEnd = 70 + startOfTOC
	}

	if startOfTOC < 0 {
		return stripEmptyNav(content), toc
	}
	// Need to peek ahead to see if this nav element is actually the right one.
	correctNav := bytes.Index(content[startOfTOC:peekEnd], []byte(`<li><a href="#`))
	if correctNav < 0 { // no match found
		return content, toc
	}
	lengthOfTOC := bytes.Index(content[startOfTOC:], last) + len(last)
	endOfTOC := startOfTOC + lengthOfTOC

	newcontent = append(content[:startOfTOC], content[endOfTOC:]...)
	toc = append(replacement, origContent[startOfTOC+len(first):endOfTOC]...)
	return
}

// RenderingContext holds contextual information, like content and configuration,
// for a given content rendering.
type RenderingContext struct {
	Content        []byte
	PageFmt        string
	DocumentID     string
	DocumentName   string
	Config         *Blackfriday
	RenderTOC      bool
	FileResolver   FileResolverFunc
	LinkResolver   LinkResolverFunc
	ConfigProvider ConfigProvider
	configInit     sync.Once
}

func newViperProvidedRenderingContext() *RenderingContext {
	return &RenderingContext{ConfigProvider: viper.GetViper()}
}

func (c *RenderingContext) getConfig() *Blackfriday {
	c.configInit.Do(func() {
		if c.Config == nil {
			c.Config = NewBlackfriday(c.ConfigProvider)
		}
	})
	return c.Config
}

// RenderBytes renders a []byte.
func RenderBytes(ctx *RenderingContext) []byte {
	switch ctx.PageFmt {
	default:
		return markdownRender(ctx)
	case "markdown":
		return markdownRender(ctx)
	case "asciidoc":
		return getAsciidocContent(ctx)
	case "mmark":
		return mmarkRender(ctx)
	case "rst":
		return getRstContent(ctx)
	}
}

// TotalWords counts instance of one or more consecutive white space
// characters, as defined by unicode.IsSpace, in s.
// This is a cheaper way of word counting than the obvious len(strings.Fields(s)).
func TotalWords(s string) int {
	n := 0
	inWord := false
	for _, r := range s {
		wasInWord := inWord
		inWord = !unicode.IsSpace(r)
		if inWord && !wasInWord {
			n++
		}
	}
	return n
}

// Old implementation only kept for benchmark comparison.
// TODO(bep) remove
func totalWordsOld(s string) int {
	return len(strings.Fields(s))
}

// TruncateWordsByRune truncates words by runes.
func TruncateWordsByRune(words []string, max int) (string, bool) {
	count := 0
	for index, word := range words {
		if count >= max {
			return strings.Join(words[:index], " "), true
		}
		runeCount := utf8.RuneCountInString(word)
		if len(word) == runeCount {
			count++
		} else if count+runeCount < max {
			count += runeCount
		} else {
			for ri := range word {
				if count >= max {
					truncatedWords := append(words[:index], word[:ri])
					return strings.Join(truncatedWords, " "), true
				}
				count++
			}
		}
	}

	return strings.Join(words, " "), false
}

// TruncateWordsToWholeSentence takes content and truncates to whole sentence
// limited by max number of words. It also returns whether it is truncated.
func TruncateWordsToWholeSentence(s string, max int) (string, bool) {

	var (
		wordCount     = 0
		lastWordIndex = -1
	)

	for i, r := range s {
		if unicode.IsSpace(r) {
			wordCount++
			lastWordIndex = i

			if wordCount >= max {
				break
			}

		}
	}

	if lastWordIndex == -1 {
		return s, false
	}

	endIndex := -1

	for j, r := range s[lastWordIndex:] {
		if isEndOfSentence(r) {
			endIndex = j + lastWordIndex + utf8.RuneLen(r)
			break
		}
	}

	if endIndex == -1 {
		return s, false
	}

	return strings.TrimSpace(s[:endIndex]), endIndex < len(s)
}

func isEndOfSentence(r rune) bool {
	return r == '.' || r == '?' || r == '!' || r == '"' || r == '\n'
}

// Kept only for benchmark.
func truncateWordsToWholeSentenceOld(content string, max int) (string, bool) {
	words := strings.Fields(content)

	if max >= len(words) {
		return strings.Join(words, " "), false
	}

	for counter, word := range words[max:] {
		if strings.HasSuffix(word, ".") ||
			strings.HasSuffix(word, "?") ||
			strings.HasSuffix(word, ".\"") ||
			strings.HasSuffix(word, "!") {
			upper := max + counter + 1
			return strings.Join(words[:upper], " "), (upper < len(words))
		}
	}

	return strings.Join(words[:max], " "), true
}

func getAsciidocExecPath() string {
	path, err := exec.LookPath("asciidoctor")
	if err != nil {
		path, err = exec.LookPath("asciidoc")
		if err != nil {
			return ""
		}
	}
	return path
}

// HasAsciidoc returns whether Asciidoctor or Asciidoc is installed on this computer.
func HasAsciidoc() bool {
	return getAsciidocExecPath() != ""
}

// getAsciidocContent calls asciidoctor or asciidoc as an external helper
// to convert AsciiDoc content to HTML.
func getAsciidocContent(ctx *RenderingContext) []byte {
	content := ctx.Content
	cleanContent := bytes.Replace(content, SummaryDivider, []byte(""), 1)

	path := getAsciidocExecPath()
	if path == "" {
		jww.ERROR.Println("asciidoctor / asciidoc not found in $PATH: Please install.\n",
			"                 Leaving AsciiDoc content unrendered.")
		return content
	}

	jww.INFO.Println("Rendering", ctx.DocumentName, "with", path, "...")
	cmd := exec.Command(path, "--no-header-footer", "--safe", "-")
	cmd.Stdin = bytes.NewReader(cleanContent)
	var out, cmderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &cmderr
	err := cmd.Run()
	// asciidoctor has exit code 0 even if there are errors in stderr
	// -> log stderr output regardless of state of err
	for _, item := range strings.Split(string(cmderr.Bytes()), "\n") {
		item := strings.TrimSpace(item)
		if item != "" {
			jww.ERROR.Println(strings.Replace(item, "<stdin>", ctx.DocumentName, 1))
		}
	}
	if err != nil {
		jww.ERROR.Printf("%s rendering %s: %v", path, ctx.DocumentName, err)
	}

	return normalizeExternalHelperLineFeeds(out.Bytes())
}

// HasRst returns whether rst2html is installed on this computer.
func HasRst() bool {
	return getRstExecPath() != ""
}

func getRstExecPath() string {
	path, err := exec.LookPath("rst2html")
	if err != nil {
		path, err = exec.LookPath("rst2html.py")
		if err != nil {
			return ""
		}
	}
	return path
}

func getPythonExecPath() string {
	path, err := exec.LookPath("python")
	if err != nil {
		path, err = exec.LookPath("python.exe")
		if err != nil {
			return ""
		}
	}
	return path
}

// getRstContent calls the Python script rst2html as an external helper
// to convert reStructuredText content to HTML.
func getRstContent(ctx *RenderingContext) []byte {
	content := ctx.Content
	cleanContent := bytes.Replace(content, SummaryDivider, []byte(""), 1)

	python := getPythonExecPath()
	path := getRstExecPath()

	if path == "" {
		jww.ERROR.Println("rst2html / rst2html.py not found in $PATH: Please install.\n",
			"                 Leaving reStructuredText content unrendered.")
		return content

	}

	jww.INFO.Println("Rendering", ctx.DocumentName, "with", path, "...")
	cmd := exec.Command(python, path, "--leave-comments")
	cmd.Stdin = bytes.NewReader(cleanContent)
	var out, cmderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &cmderr
	err := cmd.Run()
	// By default rst2html exits w/ non-zero exit code only if severe, i.e.
	// halting errors occurred. -> log stderr output regardless of state of err
	for _, item := range strings.Split(string(cmderr.Bytes()), "\n") {
		item := strings.TrimSpace(item)
		if item != "" {
			jww.ERROR.Println(strings.Replace(item, "<stdin>", ctx.DocumentName, 1))
		}
	}
	if err != nil {
		jww.ERROR.Printf("%s rendering %s: %v", path, ctx.DocumentName, err)
	}

	result := normalizeExternalHelperLineFeeds(out.Bytes())

	// TODO(bep) check if rst2html has a body only option.
	bodyStart := bytes.Index(result, []byte("<body>\n"))
	if bodyStart < 0 {
		bodyStart = -7 //compensate for length
	}

	bodyEnd := bytes.Index(result, []byte("\n</body>"))
	if bodyEnd < 0 || bodyEnd >= len(result) {
		bodyEnd = len(result) - 1
		if bodyEnd < 0 {
			bodyEnd = 0
		}
	}

	return result[bodyStart+7 : bodyEnd]
}
