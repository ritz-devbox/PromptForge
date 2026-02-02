package linter

import (
	"fmt"
	"regexp"
	"strings"
)

type Severity string

const (
	SeverityError Severity = "error"
	SeverityWarn  Severity = "warn"
)

type Diagnostic struct {
	Severity Severity
	Code     string
	Message  string
	Line     int
	Column   int
}

type sectionInfo struct {
	name      string
	startLine int
	endLine   int
}

var (
	headingPattern = regexp.MustCompile(`^##\s+(.+?)\s*$`)
	listBulletRe   = regexp.MustCompile(`^[\s]*[-*•]\s+`)
	listNumberRe   = regexp.MustCompile(`^\d+\.\s+`)
	vagueTermsRe   = regexp.MustCompile(`\b(etc|misc|various|stuff|things)\b`)
)

// LintPlan analyzes plan.md content and returns diagnostics.
func LintPlan(content []byte) []Diagnostic {
	if len(content) == 0 {
		return []Diagnostic{
			{
				Severity: SeverityError,
				Code:     "PF100",
				Message:  "plan content is empty",
				Line:     1,
				Column:   1,
			},
		}
	}

	contentStr := normalizeNewlines(string(content))
	lines := strings.Split(contentStr, "\n")
	strippedLines := stripComments(lines)

	sections := make(map[string][]sectionInfo)
	sectionBounds := make(map[string]sectionInfo)
	var order []sectionInfo
	var diagnostics []Diagnostic

	for i, line := range strippedLines {
		matches := headingPattern.FindStringSubmatch(strings.TrimSpace(line))
		if len(matches) == 0 {
			continue
		}

		name := matches[1]
		info := sectionInfo{name: name, startLine: i + 1}
		key := strings.ToLower(strings.TrimSpace(name))

		sections[key] = append(sections[key], info)
		order = append(order, info)

		if !isKnownSection(name) {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: SeverityError,
				Code:     "PF103",
				Message:  fmt.Sprintf("unknown section heading: %s", name),
				Line:     i + 1,
				Column:   1,
			})
		}
	}

	for name, infos := range sections {
		if isKnownSection(name) && len(infos) > 1 {
			for i := 1; i < len(infos); i++ {
				diagnostics = append(diagnostics, Diagnostic{
					Severity: SeverityError,
					Code:     "PF102",
					Message:  fmt.Sprintf("duplicate section heading: %s", infos[i].name),
					Line:     infos[i].startLine,
					Column:   1,
				})
			}
		}
	}

	for i, info := range order {
		info.endLine = len(lines)
		for j := i + 1; j < len(order); j++ {
			if order[j].startLine > info.startLine {
				info.endLine = order[j].startLine - 1
				break
			}
		}

		for idx := info.startLine; idx <= info.endLine && idx <= len(lines); idx++ {
			if strings.TrimSpace(lines[idx-1]) == "---" {
				info.endLine = idx - 1
				break
			}
		}

		key := strings.ToLower(strings.TrimSpace(info.name))
		if _, ok := sectionBounds[key]; !ok {
			sectionBounds[key] = info
		}
	}

	goalInfo := firstSectionInfo(sectionBounds, "goal")
	if goalInfo == nil {
		diagnostics = append(diagnostics, Diagnostic{
			Severity: SeverityError,
			Code:     "PF100",
			Message:  "missing required section: Goal",
			Line:     1,
			Column:   1,
		})
	} else {
		goalText := strings.TrimSpace(strings.Join(strippedLines[goalInfo.startLine:goalInfo.endLine], "\n"))
		if goalText == "" {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: SeverityError,
				Code:     "PF101",
				Message:  "Goal section is empty",
				Line:     goalInfo.startLine,
				Column:   1,
			})
		} else if isGoalTooShort(goalText) {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: SeverityWarn,
				Code:     "PF202",
				Message:  "Goal looks too short; add more detail",
				Line:     goalInfo.startLine,
				Column:   1,
			})
		}
	}

	constraintsInfo := firstSectionInfo(sectionBounds, "constraints")
	if constraintsInfo == nil {
		diagnostics = append(diagnostics, Diagnostic{
			Severity: SeverityWarn,
			Code:     "PF200",
			Message:  "Constraints section is missing",
			Line:     1,
			Column:   1,
		})
	} else {
		items := parseListWithLines(strippedLines, constraintsInfo.startLine, constraintsInfo.endLine)
		if len(items) == 0 {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: SeverityWarn,
				Code:     "PF200",
				Message:  "Constraints section is empty",
				Line:     constraintsInfo.startLine,
				Column:   1,
			})
		}
		for _, item := range items {
			if vagueTermsRe.MatchString(strings.ToLower(item.text)) {
				diagnostics = append(diagnostics, Diagnostic{
					Severity: SeverityWarn,
					Code:     "PF203",
					Message:  "Constraint looks vague; avoid words like 'etc' or 'various'",
					Line:     item.line,
					Column:   1,
				})
			}
		}
	}

	outOfScopeInfo := firstSectionInfo(sectionBounds, "out of scope")
	if outOfScopeInfo == nil {
		diagnostics = append(diagnostics, Diagnostic{
			Severity: SeverityWarn,
			Code:     "PF201",
			Message:  "Out of Scope section is missing",
			Line:     1,
			Column:   1,
		})
	} else {
		items := parseListWithLines(strippedLines, outOfScopeInfo.startLine, outOfScopeInfo.endLine)
		if len(items) == 0 {
			diagnostics = append(diagnostics, Diagnostic{
				Severity: SeverityWarn,
				Code:     "PF201",
				Message:  "Out of Scope section is empty",
				Line:     outOfScopeInfo.startLine,
				Column:   1,
			})
		}
	}

	return diagnostics
}

func normalizeNewlines(input string) string {
	input = strings.ReplaceAll(input, "\r\n", "\n")
	return strings.ReplaceAll(input, "\r", "\n")
}

func stripComments(lines []string) []string {
	result := make([]string, len(lines))
	inComment := false

	for i, line := range lines {
		cleaned := line
		for {
			if inComment {
				end := strings.Index(cleaned, "-->")
				if end == -1 {
					cleaned = ""
					break
				}
				cleaned = cleaned[end+3:]
				inComment = false
				continue
			}

			start := strings.Index(cleaned, "<!--")
			if start == -1 {
				break
			}
			end := strings.Index(cleaned[start+4:], "-->")
			if end == -1 {
				cleaned = cleaned[:start]
				inComment = true
				break
			}
			cleaned = cleaned[:start] + cleaned[start+4+end+3:]
		}
		result[i] = cleaned
	}

	return result
}

func isKnownSection(name string) bool {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "goal", "constraints", "out of scope":
		return true
	default:
		return false
	}
}

func firstSectionInfo(sections map[string]sectionInfo, name string) *sectionInfo {
	key := strings.ToLower(strings.TrimSpace(name))
	if info, ok := sections[key]; ok {
		return &info
	}
	return nil
}

type listItem struct {
	text string
	line int
}

func parseListWithLines(lines []string, startLine int, endLine int) []listItem {
	var items []listItem

	for idx := startLine; idx <= endLine && idx <= len(lines); idx++ {
		line := strings.TrimSpace(lines[idx-1])
		if line == "" {
			continue
		}

		trimmed := listBulletRe.ReplaceAllString(line, "")
		trimmed = listNumberRe.ReplaceAllString(trimmed, "")
		trimmed = strings.TrimSpace(trimmed)

		if trimmed == "" {
			continue
		}

		items = append(items, listItem{
			text: trimmed,
			line: idx,
		})
	}

	return items
}

func isGoalTooShort(goal string) bool {
	if len(strings.TrimSpace(goal)) < 15 {
		return true
	}
	words := strings.Fields(goal)
	return len(words) < 3
}
