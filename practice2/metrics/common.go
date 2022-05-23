package metrics

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

const (
	GO_SRC_EXTENSION = ".go"
	CUR_DIR_NAME     = "."
	PARENT_DIR_NAME  = ".."
	NEW_LINE_CHAR    = '\n'
)

type ModuleMetrics struct {
	Path                   string
	CountSourceFiles       int
	CountDirectories       int
	LinesOfCode            int
	BlankLines             int
	PhysicalLOC            int
	LogicalLOC             int
	CommentsLOC            int
	CommentSaturationLevel float64
	CyclomaticComplexity   int
}

type FileLogicalMetrics struct {
	CountSelectStmt  int
	CountIterStmt    int
	CountJumpStmt    int
	CountExprStmt    int
	CountGeneralStmt int
	CountBlockStmt   int
	CountDataDecl    int
	CountIfElse      int
	CountCases       int
}

type SourceFileInfo struct {
	TokensPos *token.FileSet
	AST       *ast.File
}

func (m *ModuleMetrics) CalculateCommentSaturationLevel() {
	m.CommentSaturationLevel = float64(m.CommentsLOC) / float64(m.LinesOfCode)
}

func (m *ModuleMetrics) String() string {
	sb := strings.Builder{}
	sb.WriteString("Path: \t\t\t\t")
	sb.WriteString(m.Path)
	sb.WriteString("\n")
	sb.WriteString("Count source files: \t\t")
	sb.WriteString(strconv.Itoa(m.CountSourceFiles))
	sb.WriteString("\n")
	sb.WriteString("Count directories: \t\t")
	sb.WriteString(strconv.Itoa(m.CountDirectories))
	sb.WriteString("\n")
	sb.WriteString("Lines of code: \t\t\t")
	sb.WriteString(strconv.Itoa(m.LinesOfCode))
	sb.WriteString("\n")
	sb.WriteString("Count blank lines: \t\t")
	sb.WriteString(strconv.Itoa(m.BlankLines))
	sb.WriteString("\n")
	sb.WriteString("Physical LOC: \t\t\t")
	sb.WriteString(strconv.Itoa(m.PhysicalLOC))
	sb.WriteString("\n")
	sb.WriteString("Logical LOC: \t\t\t")
	sb.WriteString(strconv.Itoa(m.LogicalLOC))
	sb.WriteString("\n")
	sb.WriteString("Comments LOC: \t\t\t")
	sb.WriteString(strconv.Itoa(m.CommentsLOC))
	sb.WriteString("\n")
	sb.WriteString("Comment saturation level: \t")
	sb.WriteString(fmt.Sprintf("%f", m.CommentSaturationLevel))
	sb.WriteString("\n")
	sb.WriteString("Cyclomatic complexity: \t\t")
	sb.WriteString(strconv.Itoa(m.CyclomaticComplexity))
	sb.WriteString("\n")
	return sb.String()
}

func (m *FileLogicalMetrics) CalculateLogicalLOC() int {
	return m.CountSelectStmt + m.CountIterStmt + m.CountJumpStmt + m.CountExprStmt + m.CountGeneralStmt +
		m.CountBlockStmt + m.CountDataDecl
}

func (m *FileLogicalMetrics) CalculateCyclomaticComplexity() int {
	return 1 + m.CountIfElse + m.CountCases
}

func GetSourceFileExtensions() []string {
	return []string{GO_SRC_EXTENSION}
}

func IsSourceFile(path string) bool {
	for _, extension := range GetSourceFileExtensions() {
		if strings.HasSuffix(path, extension) {
			return true
		}
	}
	return false
}

func IsHiddenDirectory(dirName string) bool {
	return strings.HasPrefix(dirName, CUR_DIR_NAME) && dirName != CUR_DIR_NAME && dirName != PARENT_DIR_NAME
}

func IsEmptyLine(line string) bool {
	return len(line) == 0 || line == string(NEW_LINE_CHAR)
}
