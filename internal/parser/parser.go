package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// Plan represents the parsed structure of plan.md
type Plan struct {
	Goal        string
	Constraints []string
	OutOfScope  []string
}

// PlanItem represents a parsed list item with line information.
type PlanItem struct {
	Text string
	Line int
}

// PlanWithLines includes line numbers for explain and lint workflows.
type PlanWithLines struct {
	Plan        *Plan
	GoalLine    int
	Constraints []PlanItem
	OutOfScope  []PlanItem
}

// ParsePlan extracts structured data from plan.md content.
// Returns an error if the Goal section is missing or if the content is invalid.
func ParsePlan(content []byte) (*Plan, error) {
	if len(content) == 0 {
		return nil, fmt.Errorf("plan content is empty")
	}

	contentStr := string(content)

	// Normalize line endings (handle Windows \r\n)
	contentStr = strings.ReplaceAll(contentStr, "\r\n", "\n")
	contentStr = strings.ReplaceAll(contentStr, "\r", "\n")

	// Extract Goal section (required)
	goalLine, goalText := extractSectionWithLine(contentStr, "Goal")
	if goalText == "" {
		return nil, fmt.Errorf("Goal section is required but not found. Expected: ## Goal")
	}

	// Extract Constraints section (optional, but warn if completely missing)
	constraintsText := extractSection(contentStr, "Constraints")
	constraints := parseList(constraintsText)

	// Extract Out of Scope section (optional)
	outOfScopeText := extractSection(contentStr, "Out of Scope")
	outOfScope := parseList(outOfScopeText)

	// Trim and validate Goal
	goal := strings.TrimSpace(goalText)
	if goal == "" {
		return nil, fmt.Errorf("Goal section is empty at line %d. Please provide a goal description", goalLine)
	}

	return &Plan{
		Goal:        goal,
		Constraints: constraints,
		OutOfScope:  outOfScope,
	}, nil
}

// ParsePlanWithLines extracts structured data from plan.md content with line numbers.
func ParsePlanWithLines(content []byte) (*PlanWithLines, error) {
	if len(content) == 0 {
		return nil, fmt.Errorf("plan content is empty")
	}

	contentStr := normalizeNewlines(string(content))
	lines := strings.Split(contentStr, "\n")
	strippedLines := stripComments(lines)

	startLine, endLine, found := sectionRange(strippedLines, "Goal")
	if !found {
		return nil, fmt.Errorf("Goal section is required but not found. Expected: ## Goal")
	}

	goalLine := -1
	var goalLines []string
	for idx := startLine; idx <= endLine && idx <= len(strippedLines); idx++ {
		line := strings.TrimSpace(strippedLines[idx-1])
		if line == "" {
			continue
		}
		if goalLine == -1 {
			goalLine = idx
		}
		goalLines = append(goalLines, line)
	}

	goal := strings.TrimSpace(strings.Join(goalLines, "\n"))
	if goal == "" {
		return nil, fmt.Errorf("Goal section is empty at line %d. Please provide a goal description", startLine)
	}

	constraintsStart, constraintsEnd, constraintsFound := sectionRange(strippedLines, "Constraints")
	var constraints []PlanItem
	if constraintsFound {
		constraints = parseListWithLines(strippedLines, constraintsStart, constraintsEnd)
	}

	outStart, outEnd, outFound := sectionRange(strippedLines, "Out of Scope")
	var outOfScope []PlanItem
	if outFound {
		outOfScope = parseListWithLines(strippedLines, outStart, outEnd)
	}

	plan := &Plan{
		Goal:        goal,
		Constraints: toTextList(constraints),
		OutOfScope:  toTextList(outOfScope),
	}

	return &PlanWithLines{
		Plan:        plan,
		GoalLine:    goalLine,
		Constraints: constraints,
		OutOfScope:  outOfScope,
	}, nil
}

// extractSection extracts content from a markdown section.
// Returns empty string if section is not found.
func extractSection(content, sectionName string) string {
	_, text := extractSectionWithLine(content, sectionName)
	return text
}

// extractSectionWithLine extracts content from a markdown section and returns the line number.
// Returns -1 for line number and empty string if section is not found.
func extractSectionWithLine(content, sectionName string) (int, string) {
	// Split content into lines for easier processing
	lines := strings.Split(content, "\n")

	// Find the section header (case-insensitive, handles extra whitespace)
	var sectionStart = -1
	headerPattern := regexp.MustCompile(`(?i)^##\s+` + regexp.QuoteMeta(sectionName) + `\s*$`)

	for i, line := range lines {
		if headerPattern.MatchString(line) {
			sectionStart = i + 1
			break
		}
	}

	if sectionStart == -1 {
		return -1, ""
	}

	// Collect lines until next section or separator
	var sectionLines []string
	nextHeaderPattern := regexp.MustCompile(`^##\s+`)

	for i := sectionStart; i < len(lines); i++ {
		line := lines[i]

		// Stop if we hit another section header
		if nextHeaderPattern.MatchString(line) {
			break
		}
		// Stop if we hit the separator
		if strings.TrimSpace(line) == "---" {
			break
		}

		sectionLines = append(sectionLines, line)
	}

	text := strings.Join(sectionLines, "\n")

	// Remove HTML comments (handles multi-line comments)
	commentRegex := regexp.MustCompile(`<!--[\s\S]*?-->`)
	text = commentRegex.ReplaceAllString(text, "")

	// Remove leading/trailing whitespace but preserve internal structure
	text = strings.TrimSpace(text)

	return sectionStart, text
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

func sectionRange(lines []string, sectionName string) (int, int, bool) {
	headerPattern := regexp.MustCompile(`(?i)^##\s+` + regexp.QuoteMeta(sectionName) + `\s*$`)
	nextHeaderPattern := regexp.MustCompile(`^##\s+`)

	sectionStart := -1
	for i, line := range lines {
		if headerPattern.MatchString(strings.TrimSpace(line)) {
			sectionStart = i + 1
			break
		}
	}
	if sectionStart == -1 {
		return -1, -1, false
	}

	sectionEnd := len(lines)
	for i := sectionStart; i < len(lines); i++ {
		line := lines[i]
		if nextHeaderPattern.MatchString(strings.TrimSpace(line)) {
			sectionEnd = i
			break
		}
		if strings.TrimSpace(line) == "---" {
			sectionEnd = i
			break
		}
	}

	return sectionStart, sectionEnd, true
}

func parseListWithLines(lines []string, startLine, endLine int) []PlanItem {
	if startLine < 1 || endLine < startLine {
		return []PlanItem{}
	}

	var items []PlanItem
	seen := make(map[string]bool)

	for idx := startLine; idx <= endLine && idx <= len(lines); idx++ {
		line := strings.TrimSpace(lines[idx-1])
		if line == "" {
			continue
		}

		line = regexp.MustCompile(`^[\s]*[-*ƒ?›]\s+`).ReplaceAllString(line, "")
		line = regexp.MustCompile(`^\d+\.\s+`).ReplaceAllString(line, "")
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		lowerLine := strings.ToLower(line)
		if seen[lowerLine] {
			continue
		}
		seen[lowerLine] = true

		items = append(items, PlanItem{
			Text: line,
			Line: idx,
		})
	}

	return items
}

func toTextList(items []PlanItem) []string {
	if len(items) == 0 {
		return []string{}
	}
	texts := make([]string, 0, len(items))
	for _, item := range items {
		text := strings.TrimSpace(item.Text)
		if text == "" {
			continue
		}
		texts = append(texts, text)
	}
	return texts
}

// parseList converts text into a list of items.
// Handles bullet points, numbered lists, and plain text lines.
// Filters out empty items and HTML comments.
func parseList(text string) []string {
	if text == "" {
		return []string{}
	}

	lines := strings.Split(text, "\n")
	var items []string
	seen := make(map[string]bool) // Track duplicates

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		// Skip HTML comments (including multi-line)
		if strings.HasPrefix(line, "<!--") {
			continue
		}

		// Remove markdown list markers (-, *, •, 1., etc.)
		line = regexp.MustCompile(`^[\s]*[-*•]\s+`).ReplaceAllString(line, "")
		line = regexp.MustCompile(`^\d+\.\s+`).ReplaceAllString(line, "")
		line = strings.TrimSpace(line)

		// Skip if empty after processing
		if line == "" {
			continue
		}

		// Skip duplicates (case-insensitive)
		lowerLine := strings.ToLower(line)
		if seen[lowerLine] {
			continue
		}
		seen[lowerLine] = true

		items = append(items, line)
	}

	// If no items found but text exists (and wasn't just comments), treat entire text as one item
	if len(items) == 0 && strings.TrimSpace(text) != "" {
		cleaned := strings.TrimSpace(text)
		// Only add if it's not just a comment
		if !strings.HasPrefix(cleaned, "<!--") {
			items = append(items, cleaned)
		}
	}

	return items
}
