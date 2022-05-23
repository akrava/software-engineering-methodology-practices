package metrics

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"strings"
)

func (m *ModuleMetrics) CountAllMetrics() error {
	return m.traverseModulePath()
}

func (m *ModuleMetrics) traverseModulePath() error {
	if fileInfo, err := os.Stat(m.Path); os.IsNotExist(err) {
		log.Printf("Path to the module %s should exists, please pass the correct path", m.Path)
		return err
	} else if !fileInfo.IsDir() && !fileInfo.Mode().IsRegular() {
		log.Printf("Path to the module %s should be a directory or regular file, please pass the correct path", m.Path)
		return fmt.Errorf("module path %s is not a directory or a regular file", m.Path)
	} else if err != nil {
		log.Printf("Error while checking stat of the module path %s", m.Path)
		return err
	}
	return m.traverseModulePathHelper(m.Path)
}

func (m *ModuleMetrics) traverseModulePathHelper(path string) error {
	if fileInfo, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Error ocurred: path to the entry %s should exists", path)
		return err
	} else if err != nil {
		log.Printf("Error ocurred: error while checking stat of the entry %s", path)
		return err
	} else if fileInfo.IsDir() {
		if IsHiddenDirectory(fileInfo.Name()) {
			return nil
		}
		m.CountDirectories++
		files, err := os.ReadDir(path)
		if err != nil {
			log.Printf("Error ocurred during reading directory %s", path)
			return err
		}
		for _, file := range files {
			if err := m.traverseModulePathHelper(path + string(os.PathSeparator) + file.Name()); err != nil {
				log.Printf("Error ocurred during traversing directory %s", path)
				return err
			}
		}
	} else if fileInfo.Mode().IsRegular() && IsSourceFile(path) {
		if err := m.countMetricsInRegularFile(path); err != nil {
			log.Printf("Error ocurred during counting metrics in file %s", path)
			return err
		}
	}
	return nil
}

func (m *ModuleMetrics) countMetricsInRegularFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("Error ocurred during opening file %s", path)
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	isEmptyFile := true
	countLines := 0
	countBlankLines := 0
	for {
		line, err := reader.ReadString(NEW_LINE_CHAR)
		if err != nil && err != io.EOF {
			log.Printf("Error ocurred during reading line from file %s", path)
			return err
		}
		if isEmptyFile && IsEmptyLine(line) && err == io.EOF {
			// empty file, do nothing
			break
		}
		isEmptyFile = false
		countLines++
		if IsEmptyLine(line) {
			countBlankLines++
		}
		// end of file, stop the loop
		if err == io.EOF {
			break
		}
	}
	if !isEmptyFile {
		m.CountSourceFiles++
	}
	m.LinesOfCode += countLines
	m.BlankLines += countBlankLines
	m.PhysicalLOC += calculatePhysicalLOC(countLines, countBlankLines)
	parsedInfo, err := parseSourceFile(path)
	if err != nil {
		log.Printf("Error ocurred during parsing file %s", path)
		return err
	}
	m.CommentsLOC += calculateCommentsLOC(parsedInfo)
	logicalMetrics := calculateLogicalMetrics(parsedInfo)
	// log.Printf("%s - %#v\n", path, logicalMetrics)
	m.LogicalLOC += logicalMetrics.CalculateLogicalLOC()
	m.CyclomaticComplexity += logicalMetrics.CalculateCyclomaticComplexity()
	return nil
}

func calculatePhysicalLOC(linesOfCode, blankLines int) int {
	blankLinesThreshold := linesOfCode / 4
	if blankLines > blankLinesThreshold {
		return linesOfCode - blankLines + blankLinesThreshold
	}
	return linesOfCode
}

func parseSourceFile(filePath string) (*SourceFileInfo, error) {
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		log.Printf("Error ocurred during parsing file %s", filePath)
		return nil, err
	}
	return &SourceFileInfo{
		TokensPos: fileSet,
		AST:       f,
	}, nil
}

func calculateCommentsLOC(s *SourceFileInfo) int {
	commentsLOC := 0
	for _, commentsGroup := range s.AST.Comments {
		for _, comment := range commentsGroup.List {
			commentsLOC += strings.Count(comment.Text, string(NEW_LINE_CHAR)) + 1
		}
	}
	return commentsLOC
}

func calculateLogicalMetrics(s *SourceFileInfo) *FileLogicalMetrics {
	logicalMetrics := &FileLogicalMetrics{}
	// node with position which present in map as key is not parsed. If value of this key is true, children are parsed,
	// otherwise they are not parsed
	skipNodesPosMap := map[token.Pos]bool{}
	ast.Inspect(s.AST, func(n ast.Node) bool {
		return logicalMetrics.parseNode(n, skipNodesPosMap)
	})
	return logicalMetrics
}

func (m *FileLogicalMetrics) parseNode(n ast.Node, skipNodesPosMap map[token.Pos]bool) bool {
	if n == nil {
		return true
	} else if parseChildren, ok := skipNodesPosMap[n.Pos()]; ok {
		delete(skipNodesPosMap, n.Pos())
		return parseChildren
	}
	switch cur := n.(type) {
	// selection statements
	case *ast.IfStmt:
		{
			m.CountSelectStmt++
			m.CountIfElse++
			if cur.Else != nil {
				if _, ok := cur.Else.(*ast.IfStmt); !ok {
					m.CountSelectStmt++
					// m.CountIfElse++
					skipNodesPosMap[cur.Else.Pos()] = true
				}
			}
			skipNodesPosMap[cur.Body.Pos()] = true
			skipNodesPosMap[cur.Cond.Pos()] = false
			if cur.Init != nil {
				skipNodesPosMap[cur.Init.Pos()] = false
			}
		}
	case *ast.SwitchStmt:
		{
			m.CountSelectStmt++
			skipNodesPosMap[cur.Body.Pos()] = true
			if cur.Init != nil {
				skipNodesPosMap[cur.Init.Pos()] = false
			}
			if cur.Tag != nil {
				skipNodesPosMap[cur.Tag.Pos()] = false
			}
		}
	case *ast.TypeSwitchStmt:
		{
			m.CountSelectStmt++
			skipNodesPosMap[cur.Body.Pos()] = true
			skipNodesPosMap[cur.Assign.Pos()] = false
			if cur.Init != nil {
				skipNodesPosMap[cur.Init.Pos()] = false
			}
		}
	// iteration statements
	case *ast.ForStmt:
		{
			m.CountIterStmt++
			skipNodesPosMap[cur.Body.Pos()] = true
			if cur.Init != nil {
				skipNodesPosMap[cur.Init.Pos()] = false
			}
			if cur.Cond != nil {
				skipNodesPosMap[cur.Cond.Pos()] = false
			}
			if cur.Post != nil {
				skipNodesPosMap[cur.Post.Pos()] = false
			}
		}
	case *ast.RangeStmt:
		{
			m.CountIterStmt++
			skipNodesPosMap[cur.Body.Pos()] = true
			if cur.Key != nil {
				skipNodesPosMap[cur.Key.Pos()] = false
			}
			if cur.Value != nil {
				skipNodesPosMap[cur.Value.Pos()] = false
			}
			skipNodesPosMap[cur.X.Pos()] = false
		}
	// jump statements
	case *ast.ReturnStmt, *ast.BranchStmt:
		{
			m.CountJumpStmt++
		}
	// expression statements
	case *ast.EmptyStmt, *ast.AssignStmt, *ast.CallExpr:
		{
			m.CountExprStmt++
		}
	// block statements
	case *ast.BlockStmt:
		{
			m.CountBlockStmt++
		}
	// data declaration statements
	case *ast.DeclStmt:
		{
			m.CountDataDecl++
		}
	// general statements
	case *ast.DeferStmt, *ast.GoStmt, *ast.SendStmt, *ast.IncDecStmt, *ast.SelectStmt:
		{
			m.CountGeneralStmt++
		}
	case *ast.ExprStmt:
		{
			if _, ok := cur.X.(*ast.CallExpr); !ok {
				m.CountGeneralStmt++
			}
		}
	// cyclomatic complexity
	case *ast.CaseClause:
		{
			if cur.List != nil {
				m.CountCases++
			}
		}
	}
	return true
}
