/*
	intoc is a TOC Generator for GitHub Flarvored Markdown.
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const ProductName = "intoc"
const ProductVersion = "0.0.1"

func ____util___() {
}

func success() {
	os.Exit(0)
}

func abort(msg string) {
	fmt.Printf("[Error!] %s\n", msg)
	os.Exit(1)
}

func file2list(filepath string) []string {
	fp, err := os.Open(filepath)
	if err != nil {
		abort(err.Error())
	}
	defer fp.Close()

	lines := []string{}

	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	return lines
}

func list2file(filepath string, lines []string) {
	fp, err := os.Create(filepath)
	if err != nil {
		abort(err.Error())
	}
	defer fp.Close()

	writer := bufio.NewWriter(fp)
	for _, line := range lines {
		writer.WriteString(line + "\n")
	}
	writer.Flush()
}

func isExist(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

func ____funcs____() {
}

func determinListPrefix (args Args) string{
	listPrefix := "-"
	if *args.useAsterisk {
		listPrefix = "*"
	}
	return listPrefix
}

type Duplicator struct {
	dict map[string]int
}

func NewDuplicator() Duplicator {
	dup := Duplicator{}
	dup.dict = map[string]int{} // map needs explicit initialization
	return dup
}

func (dup *Duplicator) Add(key string) int {
	_, isExist := dup.dict[key]
	if !isExist {
		dup.dict[key] = 1
		count_before_adding := 0
		return count_before_adding
	}
	count_before_adding := dup.dict[key]
	dup.dict[key] += 1
	return count_before_adding
}

func sectionname2anchor(sectionname string, dup *Duplicator) string {
	ret := sectionname

	ret = strings.ToLower(ret)
	ret = strings.Replace(ret, " ", "-", -1)

	reAsciiMarksWithoutHypenAndUnderscore := regexp.MustCompile("[!\"#$%&'\\(\\)\\*\\+,\\./:;<=>?@\\[\\\\\\]\\^`\\{\\|\\}~]")
	ret = reAsciiMarksWithoutHypenAndUnderscore.ReplaceAllString(ret, "")

	reJapaneseMarks := regexp.MustCompile("[、。，．・：；？！゛゜´｀¨＾￣＿ヽヾゝゞ〃仝々〆〇‐／＼～∥｜…‥‘’“”（）〔〕［］｛｝〈〉《》「」『』【】＋－±×÷＝≠＜＞≦≧∞∴♂♀°′″℃￥＄￠￡％＃＆＊＠§☆★○●◎◇◆□■△▲▽▼※〒→←↑↓〓∈∋⊆⊇⊂⊃∪∩∧∨￢⇒⇔∀∃∠⊥⌒∂∇≡≒≪≫√∽∝∵∫∬Å‰♯♭♪ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩαβγδεζηθικλμνξοπρστυφχψωАБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюя─│┌┐┘└├┬┤┴┼━┃┏┓┛┗┣┳┫┻╋┠┯┨┷┿┝┰┥┸╂｡｢｣､･ｦｧｨｩｪｫｬｭｮｯｰｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝﾞﾟ①②③④⑤⑥⑦⑧⑨⑩⑪⑫⑬⑭⑮⑯⑰⑱⑲⑳ⅠⅡⅢⅣⅤⅥⅦⅧⅨⅩ㍉㌔㌢㍍㌘㌧㌃㌶㍑㍗㌍㌦㌣㌫㍊㌻㎜㎝㎞㎎㎏㏄㎡　㍻〝〟№㏍℡㊤㊥㊦㊧㊨㈱㈲㈹㍾㍽㍼≒≡∫∮∑√⊥∠∟⊿∵∩∪]")
	ret = reJapaneseMarks.ReplaceAllString(ret, "")

	dupCount := dup.Add(ret)
	if dupCount > 0 {
		ret = fmt.Sprintf("%s-%d", ret, dupCount)
	}

	return ret
}

func line2sectioninfo(line string) (int, string) {
	sectionlevel := 0
	body := ""
	runeLine := []rune(line)

	for {
		cnt := sectionlevel + 1
		comparer := strings.Repeat("#", cnt)
		runeComparer := []rune(comparer)

		if string(runeLine[:cnt]) == string(runeComparer) {
			sectionlevel += 1
			body = string(runeLine[sectionlevel:])
			continue
		}
		break
	}

	return sectionlevel, body
}

type SectionLine struct {
	level       int
	body        string
	duplicator  *Duplicator
	listPrefix    string
	indentDepth int
}

func NewSectionLine(level int, body string, duplicator *Duplicator, listPrefix string) SectionLine {
	sl := SectionLine{}
	sl.level = level
	sl.body = body
	sl.duplicator = duplicator
	sl.listPrefix = listPrefix
	return sl
}

func (sl *SectionLine) SetIndentDepth(indentDepth int) {
	sl.indentDepth = indentDepth
}

func (sl *SectionLine) GetPlainLine() string {
	ret := strings.TrimSpace(sl.body)
	return ret
}

func (sl *SectionLine) GetTocLine() string {
	trimmedBody := strings.TrimSpace(sl.body)

	indentSpaces := strings.Repeat(" ", (sl.level-1)*sl.indentDepth)
	mark := sl.listPrefix
	text := trimmedBody
	anchor := sectionname2anchor(text, sl.duplicator)

	ret := fmt.Sprintf("%s%s [%s](#%s)", indentSpaces, mark, text, anchor)
	return ret
}

func (sl *SectionLine) GetTocLineWithoutLinkformat() string {
	trimmedBody := strings.TrimSpace(sl.body)

	indentSpaces := strings.Repeat(" ", (sl.level-1)*sl.indentDepth)
	mark := sl.listPrefix
	text := trimmedBody

	ret := fmt.Sprintf("%s%s %s", indentSpaces, mark, text)
	return ret
}

// isEditTargetLine is for judging the line is whether edit target or not.
//
// about editTarget
// - like this: "<!-- TOC"
// - Always be used head matching
// - No case-sensitive about editTarget
func isEditTargetLine(line string, editTarget string) bool {
	trimmedLine := strings.TrimSpace(line)
	trimmedLowerLine := strings.ToLower(trimmedLine)
	return strings.HasPrefix(trimmedLowerLine, strings.ToLower(editTarget))
}

// GetTocRange returns line numbers about the range of toc generated by intoc.
// But, the origins of returned values are 1, not 0.
// If not found, return -1.
func GetTocRange(lines []string, editTarget string, listPrefix string) (int, int) {
	notSet := -1
	yStart := notSet
	yEnd := notSet
	peekfirst := ""

	for i, line := range lines {
		if yStart == notSet {
			if isEditTargetLine(line, editTarget) {
				yStart = i + 1
			}
			continue
		}

		peekfirst = string([]rune(strings.TrimSpace(line))[:1])
		//fmt.Printf("%02d [%s] (listPrefix=[%s])\n", i, peekfirst, listPrefix)
		if peekfirst == listPrefix {
			yEnd = i
			continue
		}
		if yEnd != notSet {
			yEnd += 1
		}
		break
	}

	return yStart, yEnd
}

func lines2toclines(lines []string, args Args, duplicator *Duplicator) ([]string,int) {
	tocLines := []string{}

	listPrefix := determinListPrefix(args)

	isInHilight := false
	const NoEditTargetPosition = -1
	editTargetPos := NoEditTargetPosition

	for i, line := range lines {
		// skip parsing blank line
		if len(line) == 0 {
			continue
		}

		// skip parsing in the range of hilight syntax.
		peekForHilightStart := string([]rune(line)[:3])
		if peekForHilightStart == "```" {
			if isInHilight == false {
				isInHilight = true
				continue
			}
			isInHilight = false
			continue
		}
		if isInHilight == true {
			continue
		}

		// skip parsing non section line.
		// but editTarget line is also non section, must judge whether editTarget line or not firstly.
		if *args.useEdit && editTargetPos == NoEditTargetPosition && isEditTargetLine(line, *args.editTarget) {
			//  if edit target is found, memo the line number for tocline direct insertion.
			editTargetPos = i
			continue
		}
		sectionlevel, sectionbody := line2sectioninfo(line)
		if sectionlevel == 0 {
			continue
		}

		// skip parsing line unsatisfied about parseDepth
		// @todo from ">= *args.parseDepth+1" to "> *args.parseDepth" ok?
		if *args.parseDepth >= 0 && sectionlevel >= *args.parseDepth+1 {
			continue
		}

		sl := NewSectionLine(sectionlevel, sectionbody, duplicator, listPrefix)
		sl.SetIndentDepth(*args.indentDepth)
		appendeeLine := sl.GetTocLine()
		if *args.noLinkformat {
			appendeeLine = sl.GetTocLineWithoutLinkformat()
		}
		if *args.usePlainEnum {
			appendeeLine = sl.GetPlainLine()
		}
		tocLines = append(tocLines, appendeeLine)
	}
	return tocLines, editTargetPos
}

func createOutLines(lines []string, tocLines []string, args Args, editTargetPos int) ([]string) {
	// [How to construct]
	//                       - normal line
	//                       * toc line
	//
	// infile:    -----*****-------------
	//            11111     22222222222222
	//
	// toclines:   ********
	//             33333333
	//
	//            ||
	//            VV
	//
	// outlines: 111113333333322222222222222
	//           head toc     tail
	//
	// In other words, use the contents of infile basically.
	// but about TOC, use toclines generated by intoc.

	listPrefix := determinListPrefix(args)

	outLines := []string{}

	haedLines := lines[:editTargetPos+1]

	// if old TOC exists, skip parsing it.
	const NotFoundTocRange = -1
	startPos, endPos := GetTocRange(lines, *args.editTarget, listPrefix)
	skipPos := 0
	if startPos != NotFoundTocRange && endPos != NotFoundTocRange{
		skipPos = endPos - startPos
	}
	tailLines := lines[editTargetPos+1+skipPos:]

	if *args.debugPrintAll {
		fmt.Println("==== [headLines] ====")
		fmt.Println(strings.Join(haedLines, "\n"))
		fmt.Println("==== [tocLines] ====")
		fmt.Println(strings.Join(tocLines, "\n"))
		fmt.Println("==== [tailLines] ====")
		fmt.Println(strings.Join(tailLines, "\n"))
	}

	outLines = append(outLines, haedLines...)
	outLines = append(outLines, tocLines...)
	outLines = append(outLines, tailLines...)

	if *args.debugPrintAll {
		fmt.Println("==== [outLines] ====")
		fmt.Println(strings.Join(outLines, "\n"))
	}

	return outLines
}

func ____argument____() {
}

type Args struct {
	debugPrintAll *bool

	input        *string
	parseDepth   *int
	indentDepth  *int
	useAsterisk  *bool
	usePlainEnum *bool
	noLinkformat *bool

	useEdit    *bool
	editTarget *string
}

func argparse(useDefaultGetting bool) Args {
	args := Args{}

	args.input = flag.String("input", "", "A input filename.")

	args.parseDepth = flag.Int("parse-depth", -1, "The depth of the TOC list nesting. If minus then no limit depth.")
	args.indentDepth = flag.Int("indent-depth", 2, "The number of spaces per a nest in TOC.")

	args.useAsterisk = flag.Bool("use-asterisk", false, "Use an asterisk '*' as a list grammer.")
	args.usePlainEnum = flag.Bool("use-plain-enum", false, "Not use Markdown grammer, but use simple plain section name listing.")
	args.noLinkformat = flag.Bool("no-linkformat", false, "Not use '[text](#anochor)', but use 'text'.")

	args.useEdit = flag.Bool("edit", false, "If given then insert TOC to the file from '-input'.")
	args.editTarget = flag.String("edit-target", "<!-- TOC", "A insertion destination label when '-edit' given. NOT CASE-SENSITIVE.")

	args.debugPrintAll = flag.Bool("debug-print-all", false, "[DEBUG] print all options with name and value.")
	isShowingVersion := flag.Bool("version", false, "Show this intoc version.")

	flag.Parse()

	if useDefaultGetting == true{
		// For test
		return args
	}

	if *isShowingVersion{
		fmt.Printf("%s %s\n", ProductName, ProductVersion)
		success()
	}

	// Preprocess
	// ----------

	printOption := func(flg *flag.Flag) {
		fmt.Printf("%s=%s\n", flg.Name, flg.Value)
	}
	if *args.debugPrintAll {
		fmt.Println("==== Options ====")
		flag.VisitAll(printOption)
	}

	// Required
	// --------

	if *args.input == "" {
		fmt.Println("-input required.")
		flag.PrintDefaults()
		os.Exit(2)
	}

	return args
}

func main() {
	useDefaultGetting := false
	args := argparse(useDefaultGetting)

	infile, err := filepath.Abs(*args.input)
	if err != nil {
		abort("Invalid input file.")
	}
	if !isExist(infile) {
		abort("Input file does not exists.")
	}

	lines := file2list(infile)
	duplicator := NewDuplicator()

	// tocLines.
	// if not -edit then print the tocLines simply.
	const NoEditTargetPosition = -1
	tocLines, editTargetPos := lines2toclines(lines, args, &duplicator)
	if editTargetPos == NoEditTargetPosition {
		for _, line := range tocLines {
			fmt.Println(line)
		}
		success()
	}

	// outLines for re-writing infile.
	outLines := createOutLines(lines, tocLines, args, editTargetPos)
	writeeFilename := infile
	list2file(writeeFilename, outLines)
}
