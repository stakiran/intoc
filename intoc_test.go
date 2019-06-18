package main

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"strings"
)

func TestDuplicator(t *testing.T) {
	dup := NewDuplicator()

	assert.Equal(t, dup.Add("key1"), 0)
	assert.Equal(t, dup.Add("key1"), 1)
	assert.Equal(t, dup.Add("key1"), 2)
	assert.Equal(t, dup.Add("key22"), 0)
	assert.Equal(t, dup.Add("3key3"), 0)
	assert.Equal(t, dup.Add("3key3"), 1)
	assert.Equal(t, dup.Add("key22 "), 0)
}

func TestSectioname2anchor(t *testing.T) {
	dup := NewDuplicator()
	actual := ""
	expect := ""

	expect = "aaa"
	actual = sectionname2anchor("aaa", &dup)
	assert.Equal(t, actual, expect)

	expect = "あいうえお"
	actual = sectionname2anchor("あいうえお", &dup)
	assert.Equal(t, actual, expect)

	expect = "you-have-the-wrong-number"
	actual = sectionname2anchor("You have the wrong number.", &dup)
	assert.Equal(t, actual, expect)

	expect = "なんかnameerror-name-infile-is-not-definedとか出たんだけど"
	actual = sectionname2anchor("なんか「NameError: name 'infile' is not defined」とか出たんだけど", &dup)
	assert.Equal(t, actual, expect)

	expect = "aaa-1"
	actual = sectionname2anchor("aaa", &dup)
	assert.Equal(t, actual, expect, "GFM duplicated anchor numbering test.")
}

func TestLine2sectioninfo(t *testing.T) {
	lv := 0
	body := ""

	lv, body = line2sectioninfo("# 大見出し")
	assert.Equal(t, 1, lv)
	assert.Equal(t, " 大見出し", body)

	lv, body = line2sectioninfo("## 中見出し")
	assert.Equal(t, 2, lv)
	assert.Equal(t, " 中見出し", body)

	lv, body = line2sectioninfo("#### 小見出しよりもさらに小さい")
	assert.Equal(t, 4, lv)
	assert.Equal(t, " 小見出しよりもさらに小さい", body)

	lv, body = line2sectioninfo("- 見出しじゃない行")
	assert.Equal(t, 0, lv)
	assert.Equal(t, "", body)
}

func TestIsTargetLine(t *testing.T) {
	editTarget := "<!-- TOC"
	testeeLine := ""

	testeeLine = "<!-- TOC -->"
	assert.True(t, isEditTargetLine(testeeLine, editTarget))

	testeeLine = "<!-- TOC-->"
	assert.True(t, isEditTargetLine(testeeLine, editTarget))

	testeeLine = "<!--TOC-->"
	assert.False(t, isEditTargetLine(testeeLine, editTarget))

	testeeLine = "    <!-- TOC -->    "
	assert.True(t, isEditTargetLine(testeeLine, editTarget))
}

func NewSectionLineTestee(level int, body string, dup *Duplicator, listmark string, indentDepth int) SectionLine {
	sl := NewSectionLine(level, body, dup, listmark)
	if indentDepth != 0 {
		sl.SetIndentDepth(indentDepth)
	}
	return sl
}

func TestSectionLine(t *testing.T) {
	dup := NewDuplicator()
	NoIndentDepth := 0

	sl := NewSectionLineTestee(1, "大見出し", &dup, "-", 2)
	assert.Equal(t, "大見出し", sl.GetPlainLine())
	assert.Equal(t, "- [大見出し](#大見出し)", sl.GetTocLine())
	assert.Equal(t, "- 大見出し", sl.GetTocLineWithoutLinkformat())

	sl = NewSectionLineTestee(2, "中見出し", &dup, "-", 2)
	assert.Equal(t, "中見出し", sl.GetPlainLine())
	assert.Equal(t, "  - [中見出し](#中見出し)", sl.GetTocLine())
	assert.Equal(t, "  - 中見出し", sl.GetTocLineWithoutLinkformat())

	sl = NewSectionLineTestee(2, "中見出し", &dup, "-", 2)
	assert.Equal(t, "中見出し", sl.GetPlainLine())
	assert.Equal(t, "  - [中見出し](#中見出し-1)", sl.GetTocLine(), "Test duplicated body")
	assert.Equal(t, "  - 中見出し", sl.GetTocLineWithoutLinkformat())

	sl = NewSectionLineTestee(1, "中見出しインデント無し", &dup, "-", NoIndentDepth)
	assert.Equal(t, "中見出しインデント無し", sl.GetPlainLine())
	assert.Equal(t, "- [中見出しインデント無し](#中見出しインデント無し)", sl.GetTocLine())
	assert.Equal(t, "- 中見出しインデント無し", sl.GetTocLineWithoutLinkformat())

	sl = NewSectionLineTestee(3, "小見出し", &dup, "-", 2)
	assert.Equal(t, "小見出し", sl.GetPlainLine())
	assert.Equal(t, "    - [小見出し](#小見出し)", sl.GetTocLine())

	sl = NewSectionLineTestee(3, "小見出し4spaces", &dup, "-", 4)
	assert.Equal(t, "小見出し4spaces", sl.GetPlainLine())
	assert.Equal(t, "        - [小見出し4spaces](#小見出し4spaces)", sl.GetTocLine())

	sl = NewSectionLineTestee(3, "asterisk mark", &dup, "*", 4)
	assert.Equal(t, "asterisk mark", sl.GetPlainLine())
	assert.Equal(t, "        * [asterisk mark](#asterisk-mark)", sl.GetTocLine(), "Test duplicated and including-space body")
	assert.Equal(t, "        * asterisk mark", sl.GetTocLineWithoutLinkformat())
}

func TestGetTocRange(t *testing.T) {
	// origin0 := 0
	const origin1 = 1
	const notFound = -1
	thisOrigin := origin1

	editTarget := "<!-- TOC"
	listmark := "-"

	testeeLinesNew := []string{
		"# intoc", // 0
		"",        // 1
		"## Table of contents",        // 2
		"<!-- toc -->",                // 3
		"",                            // 4
		"# Overview",                  // 5
		"intoc is a TOC Generator...", // 6
		"...",       // 7
		"",          // 8
		"# Feature", // 9
	}
	yStart, yEnd := GetTocRange(testeeLinesNew, editTarget, listmark)
	assert.Equal(t, thisOrigin+3, yStart, "onepass new case.")
	assert.Equal(t, notFound, yEnd)

	testeeLinesOverwrite := []string{
		"# intoc",                               // 0
		"",                                      // 1
		"<!-- toc -->",                          // 2
		"- [intoc](#intoc)",                     // 3
		"  - [Overview](#overview)",             // 4
		"  - [Feature](#feature)",               // 5
		"  - [Install](#install)",               // 6
		"  - [Samples](#samples)",               // 7
		"    - [Basic](#basic)",                 // 8
		"    - [Depth control](#depth-control)", // 9
		"  - [License](#license)",               // 10
		"",                            // 11
		"# Overview",                  // 12
		"intoc is a TOC Generator...", // 13
		"...",       // 14
		"",          // 15
		"# Feature", // 16
	}
	yStart, yEnd = GetTocRange(testeeLinesOverwrite, editTarget, listmark)
	assert.Equal(t, thisOrigin+2, yStart, "onepass overwrite case.")
	assert.Equal(t, thisOrigin+10, yEnd)

	testeeLinesNotFound := []string{
		"# intoc", // 0
		"",        // 1
		"## Table of contents", // 2
		"",                            // 3
		"# Overview",                  // 4
		"intoc is a TOC Generator...", // 5
		"...",       // 6
		"",          // 7
		"# Feature", // 8
	}
	yStart, yEnd = GetTocRange(testeeLinesNotFound, editTarget, listmark)
	assert.Equal(t, notFound, yStart, "onepass notfound case.")
	assert.Equal(t, notFound, yEnd)
}

func ____tocLines_and_outLines_tests____(){
}

const NoEditTargetPosition = -1
const useDefaultGetting = true
var args = argparse(useDefaultGetting)

var testdata string = `# title

## sec1
aaa

### sec1-2

## せくしょん２

### セクション2-1
### セクション2-2
aiueo
### セクション2-3
#### セクション2-3-1
eof`

var testdataForEdit string = `# title

## sec1
aaa

### sec1-2

## table of contents
<!-- TOC -->

## せくしょん２

### セクション2-1
### セクション2-2
aiueo
### セクション2-3
#### セクション2-3-1
eof`

var testdataForEditOverwrite string = `# title

## sec1
aaa

### sec1-2

## table of contents
<!-- TOC -->
- TOC
  - already
    - exists.

## せくしょん２

### セクション2-1
### セクション2-2
aiueo
### セクション2-3
#### セクション2-3-1
eof`

var expectForEdit string = `# title

## sec1
aaa

### sec1-2

## table of contents
<!-- TOC -->
- title
  - sec1
    - sec1-2
  - table of contents
  - せくしょん２
    - セクション2-1
    - セクション2-2
    - セクション2-3
      - セクション2-3-1

## せくしょん２

### セクション2-1
### セクション2-2
aiueo
### セクション2-3
#### セクション2-3-1
eof`

var testdataForHilightInclusion string = `# including hilight syntax

## python

` + "```" + `py
# hello
print("hello world.")
` + "```" + `

## markdown

` + "```" + `
# level1
## level2
### level3
- list
    - list2
        - list3
` + "```" + `
`

// PrepareDefaultPrameters prepare default parameters for outLines testing.
// This parameters do not mean intoc's default.
func PrepareDefaultPrameters(testdata string) ([]string, Duplicator) {
	const useDefaultGetting = true

	lines := strings.Split(testdata, "\n")
	duplicator := NewDuplicator()
	ClearArgsParameter(&args)

	return lines, duplicator
}

func ClearArgsParameter(args *Args) () {
	*args.useAsterisk = false
	*args.noLinkformat = false
	*args.usePlainEnum = false
	*args.useEdit = false
	*args.editTarget = "<!-- TOC"
	*args.parseDepth = -1
	*args.indentDepth = 2
}

func TestOutLinesCase_DepthOptions(t *testing.T) {
	expect := `- [title](#title)
    - [sec1](#sec1)
    - [せくしょん２](#せくしょん２)`

	lines, duplicator := PrepareDefaultPrameters(testdata)
	*args.parseDepth = 2
	*args.indentDepth = 4

	tocLines, editTargetPos := lines2toclines(lines, args, &duplicator)
	actual := strings.Join(tocLines, "\n")

	assert.Equal(t, NoEditTargetPosition, editTargetPos)
	assert.Equal(t, expect, actual)
}

func TestOutLinesCase_NolinkformatAndUseAsterisk(t *testing.T) {
	expect := `* title
  * sec1
    * sec1-2
  * せくしょん２
    * セクション2-1
    * セクション2-2
    * セクション2-3
      * セクション2-3-1`

	lines, duplicator := PrepareDefaultPrameters(testdata)
	*args.noLinkformat = true
	*args.useAsterisk = true

	tocLines, editTargetPos := lines2toclines(lines, args, &duplicator)
	actual := strings.Join(tocLines, "\n")

	assert.Equal(t, NoEditTargetPosition, editTargetPos)
	assert.Equal(t, expect, actual)
}

func TestOutLinesCase_PlainEnumeration(t *testing.T) {
	expect := `title
sec1
sec1-2
せくしょん２
セクション2-1
セクション2-2
セクション2-3
セクション2-3-1`

	lines, duplicator := PrepareDefaultPrameters(testdata)
	*args.usePlainEnum = true

	tocLines, editTargetPos := lines2toclines(lines, args, &duplicator)
	actual := strings.Join(tocLines, "\n")

	assert.Equal(t, NoEditTargetPosition, editTargetPos)
	assert.Equal(t, expect, actual)
}

func TestOutLinesCase_EditNew(t *testing.T) {
	lines, duplicator := PrepareDefaultPrameters(testdataForEdit)
	*args.useEdit = true
	*args.noLinkformat = true // For simplifing expect data

	tocLines, editTargetPos := lines2toclines(lines, args, &duplicator)
	assert.NotEqual(t, NoEditTargetPosition, editTargetPos)

	outLines := createOutLines(lines, tocLines, args, editTargetPos)
	actual := strings.Join(outLines, "\n")
	expect := expectForEdit
	assert.Equal(t, expect, actual)
}

func TestOutLinesCase_EditOverwrite(t *testing.T) {
	lines, duplicator := PrepareDefaultPrameters(testdataForEditOverwrite)
	*args.useEdit = true
	*args.noLinkformat = true

	tocLines, editTargetPos := lines2toclines(lines, args, &duplicator)
	assert.NotEqual(t, NoEditTargetPosition, editTargetPos)

	outLines := createOutLines(lines, tocLines, args, editTargetPos)
	actual := strings.Join(outLines, "\n")
	expect := expectForEdit
	assert.Equal(t, expect, actual)
}

func TestOutLinesCase_EditNotFound(t *testing.T) {
	lines, duplicator := PrepareDefaultPrameters(testdata)
	*args.useEdit = true
	*args.noLinkformat = true

	_, editTargetPos := lines2toclines(lines, args, &duplicator)
	assert.Equal(t, NoEditTargetPosition, editTargetPos)
}

func TestOutLinesCase_HilightInclusion(t *testing.T) {
	expect := `- including hilight syntax
  - python
  - markdown`

	lines, duplicator := PrepareDefaultPrameters(testdataForHilightInclusion)
	*args.noLinkformat = true

	tocLines, editTargetPos := lines2toclines(lines, args, &duplicator)
	actual := strings.Join(tocLines, "\n")

	assert.Equal(t, NoEditTargetPosition, editTargetPos)
	assert.Equal(t, expect, actual)
}
