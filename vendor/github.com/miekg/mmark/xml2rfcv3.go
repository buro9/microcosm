// XML2RFC v3 rendering backend

package mmark

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// XML renderer configuration options.
const (
	XML_STANDALONE = 1 << iota // create standalone document
)

var words2119 = map[string]bool{
	"MUST":        true,
	"MUST NOT":    true,
	"REQUIRED":    true,
	"SHALL":       true,
	"SHALL NOT":   true,
	"SHOULD":      true,
	"SHOULD NOT":  true,
	"RECOMMENDED": true,
	"MAY":         true,
	"OPTIONAL":    true,
}

// Xml is a type that implements the Renderer interface for XML2RFV3 output.
//
// Do not create this directly, instead use the XmlRenderer function.
type xml struct {
	flags          int  // XML_* options
	sectionLevel   int  // current section level
	docLevel       int  // frontmatter/mainmatter or backmatter
	part           bool // parts cannot nest, if true a part has been opened
	specialSection int
	para           bool // when true we're in a para, artworks need to close it first then.

	// Store the IAL we see for this block element
	ial *inlineAttr

	// TitleBlock in TOML
	titleBlock *title
}

// XmlRenderer creates and configures a Xml object, which
// satisfies the Renderer interface.
//
// flags is a set of XML_* options ORed together
func XmlRenderer(flags int) Renderer { return &xml{flags: flags} }
func (options *xml) Flags() int      { return options.flags }
func (options *xml) State() int      { return 0 }

func (options *xml) SetAttr(i *inlineAttr) {
	options.ial = i
}

func (options *xml) Attr() *inlineAttr {
	if options.ial == nil {
		return newInlineAttr()
	}
	return options.ial
}

func (options *xml) AttrString(i *inlineAttr) string {
	if i == nil {
		return ""
	}
	s := ""
	if i.id != "" {
		s = " anchor=\"" + i.id + "\""
	}

	keys := i.SortClasses()
	if len(keys) > 0 {
		s += " class=\"" + strings.Join(keys, " ") + "\""
	}

	keys = i.SortAttributes()
	attr := make([]string, len(keys))
	for j, k := range keys {
		v := i.attr[k]
		attr[j] = k + "=\"" + v + "\""
	}
	if len(keys) > 0 {
		s += " " + strings.Join(attr, " ")
	}
	return s
}

// render code chunks using verbatim, or listings if we have a language
func (options *xml) BlockCode(out *bytes.Buffer, text []byte, lang string, caption []byte, subfigure, callout bool) {
	if options.para {
		// close it
		out.WriteString("</t>")
		defer out.WriteString("<t>")
	}

	// Tick of language for sourcecode...
	ial := options.Attr()
	if lang != "" {
		ial.GetOrDefaultAttr("type", lang)
	}
	prefix := ial.Value("prefix")
	ial.DropAttr("prefix") // it's a fake attribute, so drop it

	s := options.AttrString(ial)

	text = blockCodePrefix(prefix, text)

	// if in a figure quote suppress <figure> and caption use
	if !subfigure && len(caption) > 0 {
		out.WriteString("<figure" + s + ">\n")
		s = ""
		out.WriteString("<name>")
		out.Write(caption)
		out.WriteString("</name>\n")
	}

	if lang != "" {
		out.WriteString("\n<sourcecode" + s + ">\n")
	} else {
		out.WriteString("<artwork" + s + ">\n")
	}
	writeEntity(out, text)

	if lang != "" {
		out.WriteString("</sourcecode>\n")
	} else {
		out.WriteString("</artwork>\n")
	}
	if !subfigure && len(caption) > 0 {
		out.WriteString("</figure>\n")
	}
}

func (options *xml) CalloutCode(out *bytes.Buffer, index, id string) {
	printf(nil, "TODO implement: CalloutCode")
}

func (options *xml) CalloutText(out *bytes.Buffer, index string, id []string) {
	printf(nil, "TODO implement: CalloutText")
}

func (options *xml) TitleBlockTOML(out *bytes.Buffer, block *title) {
	if options.flags&XML_STANDALONE == 0 {
		return
	}
	// Processing Instructions are attribute of <rfc> know.
	options.titleBlock = block
	out.WriteString("<rfc xmlns:xi=\"http://www.w3.org/2001/XInclude\"")
	out.WriteString(" ipr=\"" + options.titleBlock.Ipr + "\"")
	out.WriteString(" category=\"" + options.titleBlock.Category + "\"")
	out.WriteString(" docName=\"" + options.titleBlock.DocName + "\">")
	if len(options.titleBlock.Updates) > 0 {
		updates := make([]string, len(options.titleBlock.Updates))
		for i := range updates {
			updates[i] = strconv.Itoa(options.titleBlock.Updates[i])
		}
		out.WriteString(" updates=\"" + strings.Join(updates, ", ") + "\"")
	}
	if len(options.titleBlock.Obsoletes) > 0 {
		obsoletes := make([]string, len(options.titleBlock.Obsoletes))
		for i := range obsoletes {
			obsoletes[i] = strconv.Itoa(options.titleBlock.Obsoletes[i])
		}
		out.WriteString(" obsoletes=\"" + strings.Join(obsoletes, ", ") + "\"")
	}
	out.WriteString("\n")
	out.WriteString("<front>\n")
	out.WriteString("<title abbrev=\"" + options.titleBlock.Abbrev + "\">")
	out.WriteString(options.titleBlock.Title + "</title>\n\n")

	for _, a := range options.titleBlock.Author {
		titleBlockTOMLAuthor(out, a)
	}

	titleBlockTOMLDate(out, options.titleBlock.Date)

	out.WriteString("<area>" + options.titleBlock.Area + "</area>\n")
	out.WriteString("<workgroup>" + options.titleBlock.Workgroup + "</workgroup>\n")

	titleBlockTOMLKeyword(out, options.titleBlock.Keyword)
	out.WriteString("\n")
}

func (options *xml) BlockQuote(out *bytes.Buffer, text []byte, attribution []byte) {
	// check for "person -- URI" syntax use those if found
	// need to strip tags because these are attributes
	ial := options.Attr()
	if len(attribution) != 0 {
		parts := bytes.Split(attribution, []byte(" -- "))
		if len(parts) == 2 {
			cite := string(bytes.TrimSpace(parts[0]))
			quotedFrom := sanitizeXML(bytes.TrimSpace(parts[1]))
			ial.GetOrDefaultAttr("cite", cite)
			ial.GetOrDefaultAttr("quotedFrom", string(quotedFrom))
		}
	}

	out.WriteString("<blockquote" + options.AttrString(ial) + ">\n")
	out.Write(text)
	out.WriteString("</blockquote>\n")
}

func (options *xml) Aside(out *bytes.Buffer, text []byte) {
	ial := options.Attr()
	s := options.AttrString(ial)
	out.WriteString("<aside" + s + ">\n")
	out.Write(text)
	out.WriteString("</aside>\n")
}

func (options *xml) CommentHtml(out *bytes.Buffer, text []byte) {
	// nothing fancy any left of the first `:` will be used as the source="..."
	// if the syntax is different, don't output anything.
	i := bytes.Index(text, []byte("-->"))
	if i > 0 {
		text = text[:i]
	}
	// strip, <!--
	text = text[4:]

	var source []byte
	l := len(text)
	if l > 20 {
		l = 20
	}
	for i := 0; i < l; i++ {
		if text[i] == '-' && text[i+1] == '-' {
			source = text[:i]
			text = text[i+2:]
			break
		}
	}
	// don't output a cref if it is not name: remark
	if len(source) != 0 {
		// sanitize source here
		source = bytes.TrimSpace(source)
		text = bytes.TrimSpace(text)
		out.WriteString("<t><cref source=\"")
		out.Write(source)
		out.WriteString("\">")
		out.Write(text)
		out.WriteString("</cref></t>\n")
	}
	return
}

func (options *xml) BlockHtml(out *bytes.Buffer, text []byte) {
	// not supported, don't know yet if this is useful
	return
}

func (options *xml) Part(out *bytes.Buffer, text func() bool, id string) {
	printf(nil, "syntax not supported: Part")
}

func (options *xml) Note(out *bytes.Buffer, text func() bool, id string) {
	switch options.specialSection {
	case _ABSTRACT:
		out.WriteString("</abstract>\n\n")
	case _NOTE:
		out.WriteString("</note>\n\n")
	}
	level := 1
	if level <= options.sectionLevel {
		// close previous ones
		for i := options.sectionLevel - level + 1; i > 0; i-- {
			out.WriteString("</section>\n")
		}
	}

	ial := options.Attr()

	out.WriteString("\n<note" + options.AttrString(ial) + ">\n")
	out.WriteString("<name>")
	text()
	out.WriteString("</name>\n")
	options.sectionLevel = 0
	options.specialSection = _NOTE
	return
}

func (options *xml) SpecialHeader(out *bytes.Buffer, what []byte, text func() bool, id string) {
	if string(what) == "preface" {
		printf(nil, "handling preface like abstract")
		what = []byte("abstract")
	}
	switch options.specialSection {
	case _ABSTRACT:
		out.WriteString("</abstract>\n\n")
	case _NOTE:
		out.WriteString("</note>\n\n")
	}
	level := 1
	if level <= options.sectionLevel {
		// close previous ones
		for i := options.sectionLevel - level + 1; i > 0; i-- {
			out.WriteString("</section>\n")
		}
	}

	ial := options.Attr()

	out.WriteString("\n<abstract" + options.AttrString(ial) + ">\n")
	options.sectionLevel = 0
	options.specialSection = _ABSTRACT
	return
}

func (options *xml) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	switch options.specialSection {
	case _ABSTRACT:
		out.WriteString("</abstract>\n\n")
	case _NOTE:
		out.WriteString("</note>\n\n")
	}
	if level <= options.sectionLevel {
		// close previous ones
		for i := options.sectionLevel - level + 1; i > 0; i-- {
			out.WriteString("</section>\n")
		}
	}

	ial := options.Attr()
	ial.GetOrDefaultId(id)

	// new section
	out.WriteString("\n<section" + options.AttrString(ial) + ">\n")
	out.WriteString("<name>")
	text()
	out.WriteString("</name>\n")
	options.sectionLevel = level
	options.specialSection = 0
	return
}

func (options *xml) HRule(out *bytes.Buffer) {
	printf(nil, "syntax not supported: HRule")
}

func (options *xml) List(out *bytes.Buffer, text func() bool, flags, start int, group []byte) {
	marker := out.Len()

	ial := options.Attr()
	ial.KeepAttr([]string{"type", "start", "group", "spacing", "empty"})

	if start > 1 {
		ial.GetOrDefaultAttr("start", strconv.Itoa(start))
	}
	if group != nil {
		ial.GetOrDefaultAttr("group", string(group))
	}

	if flags&_LIST_TYPE_ORDERED != 0 {
		switch {
		case flags&_LIST_TYPE_ORDERED_ALPHA_LOWER != 0:
			ial.GetOrDefaultAttr("type", "%c")
		case flags&_LIST_TYPE_ORDERED_ALPHA_UPPER != 0:
			ial.GetOrDefaultAttr("type", "%C")
		case flags&_LIST_TYPE_ORDERED_ROMAN_LOWER != 0:
			ial.GetOrDefaultAttr("type", "%i")
		case flags&_LIST_TYPE_ORDERED_ROMAN_UPPER != 0:
			ial.GetOrDefaultAttr("type", "%I")
		case flags&_LIST_TYPE_ORDERED_GROUP != 0:
			// ? TODO(miek)
			// Handled above, don't need to do anything because v3 format will take care of this.
		}
	}

	s := options.AttrString(ial)

	switch {
	case flags&_LIST_TYPE_ORDERED != 0:
		out.WriteString("<ol" + s + ">\n")
	case flags&_LIST_TYPE_DEFINITION != 0:
		out.WriteString("<dl" + s + ">\n")
	default:
		out.WriteString("<ul" + s + ">\n")
	}

	if !text() {
		out.Truncate(marker)
		return
	}
	switch {
	case flags&_LIST_TYPE_ORDERED != 0:
		out.WriteString("</ol>\n")
	case flags&_LIST_TYPE_DEFINITION != 0:
		out.WriteString("</dl>\n")
	default:
		out.WriteString("</ul>\n")
	}
}

func (options *xml) ListItem(out *bytes.Buffer, text []byte, flags int) {
	if flags&_LIST_TYPE_DEFINITION != 0 && flags&_LIST_TYPE_TERM == 0 {
		out.WriteString("<dd>")
		out.Write(text)
		out.WriteString("</dd>\n")
		return
	}
	if flags&_LIST_TYPE_TERM != 0 {
		out.WriteString("<dt>")
		out.Write(text)
		out.WriteString("</dt>\n")
		return
	}
	out.WriteString("<li>")
	out.Write(text)
	out.WriteString("</li>\n")
}

func (options *xml) Example(out *bytes.Buffer, index int) {
	out.WriteByte('(')
	out.WriteString(strconv.Itoa(index))
	out.WriteByte(')')
}

func (options *xml) Paragraph(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()
	options.para = true
	defer func() { options.para = false }()
	out.WriteString("<t>\n")
	if !text() {
		out.Truncate(marker)
		return
	}
	if marker+3 == out.Len() { // empty paragraph, suppress
		out.Truncate(marker)
		return
	}
	out.WriteByte('\n')
	out.WriteString("</t>\n")
}

func (options *xml) Math(out *bytes.Buffer, text []byte, display bool) {
	printf(nil, "syntax not supported: Math")
}

func (options *xml) Table(out *bytes.Buffer, header []byte, body []byte, footer []byte, columnData []int, caption []byte) {
	ial := options.Attr()
	s := options.AttrString(ial)
	out.WriteString("<table" + s + ">\n")
	if caption != nil {
		out.WriteString("<name>")
		out.Write(caption)
		out.WriteString("</name>\n")
	}
	out.WriteString("<thead>\n")
	out.Write(header)
	out.WriteString("</thead>\n")
	out.Write(body)
	out.WriteString("<tfoot>\n")
	out.Write(header)
	out.WriteString("</tfoot>\n")
	out.WriteString("</table>\n")
}

func (options *xml) TableRow(out *bytes.Buffer, text []byte) {
	out.WriteString("<tr>")
	out.Write(text)
	out.WriteString("</tr>\n")
}

func (options *xml) TableHeaderCell(out *bytes.Buffer, text []byte, align, colspan int) {
	a := ""

	if colspan > 1 {
		a = fmt.Sprintf(" colspan=\"%d\"", colspan)
	}

	switch align {
	case _TABLE_ALIGNMENT_LEFT:
		a += " align=\"left\""
	case _TABLE_ALIGNMENT_RIGHT:
		a += " align=\"right\""
	default:
		a += " align=\"center\""
	}
	out.WriteString("<th" + a + ">")
	out.Write(text)
	out.WriteString("</th>")
}

func (options *xml) TableCell(out *bytes.Buffer, text []byte, align, colspan int) {
	col := ""
	if colspan > 1 {
		col = fmt.Sprintf(" colspan=\"%d\"", colspan)
	}
	out.WriteString("<td" + col + ">")
	out.Write(text)
	out.WriteString("</td>")
}

func (options *xml) Footnotes(out *bytes.Buffer, text func() bool) {
	printf(nil, "syntax not supported: Footnotes")
}

func (options *xml) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	printf(nil, "syntax not supported: FootnoteItem")
}

func (options *xml) Index(out *bytes.Buffer, primary, secondary []byte, prim bool) {
	p := ""
	if prim {
		p = " primary=\"true\""
	}
	out.WriteString("<iref item=\"" + string(primary) + "\"" + p)
	out.WriteString(" subitem=\"" + string(secondary) + "\"" + "/>")
}

func (options *xml) Citation(out *bytes.Buffer, link, title []byte) {
	if len(title) == 0 {
		out.WriteString("<xref target=\"" + string(link) + "\"/>")
		return
	}
	out.WriteString("<xref target=\"" + string(link) + "\" section=\"" + string(title) + "\"/>")
}

func (options *xml) References(out *bytes.Buffer, citations map[string]*citation) {
	if options.flags&XML_STANDALONE == 0 {
		return
	}
	// close any option section tags
	for i := options.sectionLevel; i > 0; i-- {
		out.WriteString("</section>\n")
		options.sectionLevel--
	}
	switch options.docLevel {
	case _DOC_FRONT_MATTER:
		out.WriteString("</front>\n")
		out.WriteString("<back>\n")
	case _DOC_MAIN_MATTER:
		out.WriteString("</middle>\n")
		out.WriteString("<back>\n")
	case _DOC_BACK_MATTER:
		// nothing to do
	}
	options.docLevel = _DOC_BACK_MATTER

	refi, refn, keys := countCitationsAndSort(citations)

	// output <xi:include href="<references file>.xml"/>, we use file it its not empty, otherwise
	// we construct one for RFCNNNN and I-D.something something.
	if refi+refn > 0 {
		if refn > 0 {
			out.WriteString("<references>\n")
			out.WriteString("<name>Normative References</name>\n")
			for _, k := range keys {
				c := citations[k]
				if c.typ == 'n' {
					if c.xml != nil {
						out.Write(c.xml)
						out.WriteByte('\n')
						continue
					}
					f := referenceFile(c)
					out.WriteString("<xi:include href=\"" + f + "\"/>\n")
				}
			}
			out.WriteString("</references>\n")
		}
		if refi > 0 {
			// This needs an anchor
			out.WriteString("<references>\n")
			out.WriteString("<name>Informative References</name>\n")
			for _, k := range keys {
				c := citations[k]
				if c.typ == 'i' {
					// if we have raw xml, output that
					if c.xml != nil {
						out.Write(c.xml)
						out.WriteByte('\n')
						continue
					}
					f := referenceFile(c)
					out.WriteString("<xi:include href=\"" + f + "\"/>\n")
				}
			}
			out.WriteString("</references>\n")
		}
	}
}

func (options *xml) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	out.WriteString("<eref target=\"")
	if kind == _LINK_TYPE_EMAIL {
		out.WriteString("mailto:")
	}
	out.Write(link)
	out.WriteString("\"/>")
}

func (options *xml) CodeSpan(out *bytes.Buffer, text []byte) {
	out.WriteString("<tt>")
	writeEntity(out, text)
	out.WriteString("</tt>")
}

func (options *xml) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	// Check for 2119 Keywords
	s := string(text)
	if _, ok := words2119[s]; ok {
		out.WriteString("<bcp14>")
		out.Write(text)
		out.WriteString("</bcp14>")
		return
	}
	out.WriteString("<strong>")
	out.Write(text)
	out.WriteString("</strong>")
}

func (options *xml) Emphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("<em>")
	out.Write(text)
	out.WriteString("</em>")
}

func (options *xml) Subscript(out *bytes.Buffer, text []byte) {
	out.WriteString("<sub>")
	out.Write(text)
	out.WriteString("</sub>")
}

func (options *xml) Superscript(out *bytes.Buffer, text []byte) {
	out.WriteString("<sup>")
	out.Write(text)
	out.WriteString("</sup>")
}

func (options *xml) Figure(out *bytes.Buffer, text []byte, caption []byte) {
	// add figure and typeset the caption
	ial := options.Attr()
	s := options.AttrString(ial)
	out.WriteString("<figure" + s + ">\n")
	out.WriteString("<name>")
	out.Write(caption)
	out.WriteString("</name>\n")
	out.Write(text)
	out.WriteString("</figure>\n")
}

func (options *xml) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte, subfigure bool) {
	// use title as caption is we have it and wrap everything in a figure
	// check the extension of the local include to set the type of the thing.
	if options.para {
		// close it
		out.WriteString("</t>")
		defer out.WriteString("<t>")
	}

	// if subfigure, no <figure>
	s := options.AttrString(options.Attr())
	if bytes.HasPrefix(link, []byte("http://")) || bytes.HasPrefix(link, []byte("https://")) {
		// link to external entity
		out.WriteString("<artwork" + s)
		out.WriteString(" alt=\"")
		out.Write(alt)
		out.WriteString("\"")
		out.WriteString(" src=\"")
		out.Write(link)
		out.WriteString("\"/>")
	} else {
		// local file, xi:include it
		out.WriteString("<artwork" + s)
		out.WriteString(" alt=\"")
		out.Write(alt)
		out.WriteString("\">")
		out.WriteString("<xi:include href=\"")
		out.Write(link)
		out.WriteString("\"/>\n")
		out.WriteString("</artwork>\n")
	}
}

func (options *xml) LineBreak(out *bytes.Buffer) {
	out.WriteString("\n<br/>\n")
}

func (options *xml) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	if link[0] == '#' {
		out.WriteString("<xref target=\"")
		out.Write(link[1:])
		out.WriteString("\"/>")
		return
	}
	out.WriteString("<eref target=\"")
	out.Write(link)
	out.WriteString("\">")
	out.Write(content)
	out.WriteString("</eref>")
}

func (options *xml) Abbreviation(out *bytes.Buffer, abbr, title []byte) {
	out.Write(abbr)
}

func (options *xml) RawHtmlTag(out *bytes.Buffer, tag []byte) {
	switch {
	case bytes.Compare(tag, []byte("<br/>")) == 0:
		out.WriteString("<vspace/>")
		return
	}
	printf(nil, "syntax not supported: RawHtmlTag: %s", string(tag))
}

func (options *xml) TripleEmphasis(out *bytes.Buffer, text []byte) {
	out.WriteString("<strong><em>")
	out.Write(text)
	out.WriteString("</em></strong>")
}

func (options *xml) StrikeThrough(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (options *xml) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	printf(nil, "syntax not supported: FootnoteRef")
}

func (options *xml) Entity(out *bytes.Buffer, entity []byte) {
	out.Write(entity)
}

func (options *xml) NormalText(out *bytes.Buffer, text []byte) {
	attrEscape(out, text)
}

// header and footer
func (options *xml) DocumentHeader(out *bytes.Buffer, first bool) {
	if !first || options.flags&XML_STANDALONE == 0 {
		return
	}
	out.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
}

func (options *xml) DocumentFooter(out *bytes.Buffer, first bool) {
	if !first {
		return
	}
	switch options.specialSection {
	case _ABSTRACT:
		out.WriteString("</abstract>\n\n")
	case _NOTE:
		out.WriteString("</note>\n\n")
	}
	// close any option section tags
	for i := options.sectionLevel; i > 0; i-- {
		out.WriteString("</section>\n")
		options.sectionLevel--
	}
	if options.flags&XML_STANDALONE == 0 {
		return
	}
	switch options.docLevel {
	case _DOC_FRONT_MATTER:
		out.WriteString("\n</front>\n")
	case _DOC_MAIN_MATTER:
		out.WriteString("\n</middle>\n")
	case _DOC_BACK_MATTER:
		out.WriteString("\n</back>\n")
	}
	out.WriteString("</rfc>\n")
}

func (options *xml) DocumentMatter(out *bytes.Buffer, matter int) {
	if options.flags&XML_STANDALONE == 0 {
		return
	}
	switch options.specialSection {
	case _ABSTRACT:
		out.WriteString("</abstract>\n\n")
	case _NOTE:
		out.WriteString("</note>\n\n")
	}
	// we default to frontmatter already openened in the documentHeader
	for i := options.sectionLevel; i > 0; i-- {
		out.WriteString("</section>\n")
		options.sectionLevel--
	}
	switch matter {
	case _DOC_FRONT_MATTER:
		// already open
	case _DOC_MAIN_MATTER:
		out.WriteString("</front>\n")
		out.WriteString("\n<middle>\n")
	case _DOC_BACK_MATTER:
		out.WriteString("\n</middle>\n")
		out.WriteString("<back>\n")
	}
	options.docLevel = matter
	options.specialSection = 0
}
