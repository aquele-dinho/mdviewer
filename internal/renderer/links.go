package renderer

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// Obsidian wiki link: [[target]] or [[target|Label]]
	wikiLinkRegexp = regexp.MustCompile(`\[\[([^\]|]+)(\|([^\]]+))?\]\]`)
	// Obsidian image embed: ![[path]] or ![[path|size]]
	imageEmbedRegexp = regexp.MustCompile(`!\[\[([^\]|]+)(\|([0-9]+))?\]\]`)
)

// PreprocessLinks rewrites Obsidian-style links and embeds into standard
// Markdown links/images so Glamour can style them normally.
func PreprocessLinks(content string) string {
	lines := strings.Split(content, "\n")
	var out []string
	inCodeFence := false

	for _, line := range lines {
		trim := strings.TrimSpace(line)

		// Track fenced code blocks; do not rewrite inside them.
		if strings.HasPrefix(trim, "```") {
			inCodeFence = !inCodeFence
			out = append(out, line)
			continue
		}

		if inCodeFence {
			out = append(out, line)
			continue
		}

		processed := line

		// Rewrite Obsidian image embeds first (they start with ![[).
		processed = imageEmbedRegexp.ReplaceAllStringFunc(processed, func(match string) string {
			m := imageEmbedRegexp.FindStringSubmatch(match)
			if len(m) < 2 {
				return match
			}
			path := strings.TrimSpace(m[1])
			width := strings.TrimSpace(m[3]) // Optional width from |width syntax

			// Normalize path: if it's bare, assume ./path
			if !strings.HasPrefix(path, "./") && !strings.HasPrefix(path, "/") {
				path = "./" + path
			}

			// Alt text: base name without extension, optionally with width hint
			base := filepath.Base(path)
			if dot := strings.LastIndex(base, "."); dot != -1 {
				base = base[:dot]
			}
			
			// If width specified, encode it in the alt text for our detector to find
			if width != "" {
				base = base + "|width=" + width
			}

			return "![" + base + "](" + path + ")"
		})

		// Rewrite wiki links [[target]] / [[target|Label]].
		processed = wikiLinkRegexp.ReplaceAllStringFunc(processed, func(match string) string {
			m := wikiLinkRegexp.FindStringSubmatch(match)
			if len(m) < 2 {
				return match
			}
			target := strings.TrimSpace(m[1])
			label := strings.TrimSpace(m[3])
			if label == "" {
				label = target
			}

			// If target already looks like a path or has an extension, keep it.
			href := target
			if !strings.Contains(href, "/") && !strings.Contains(href, ".") {
				// Map note name to ./name.md
				href = "./" + href + ".md"
			} else if !strings.HasPrefix(href, "./") && !strings.HasPrefix(href, "/") {
				href = "./" + href
			}

			return "[" + label + "](" + href + ")"
		})

		out = append(out, processed)
	}

	return strings.Join(out, "\n")
}