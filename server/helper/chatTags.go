package helper

import (
	"regexp"
	"strings"
)

func ParseChatTags(str string) string {
	replaceMap := map[string]string{
		"<":    "&lt;",
		">":    "&gt;",
		"\r\n": "<br>",
		"\n":   "<br>",
		"\t":   "    ",
		"\\'":  "'",
		"\\\\": "\\",
		"\\\"": "&quot",
	}

	for oldStr, newStr := range replaceMap {
		str = strings.ReplaceAll(str, oldStr, newStr)
	}

	itemTagRegex := regexp.MustCompile(`\[i(.*?):(\w+)\]`)
	str = itemTagRegex.ReplaceAllString(str, "<img src=\"https://tmlapis.repl.co/img/Item_$2.png\" id=\"item-icon\">")
	colorTagRegex := regexp.MustCompile(`\[c\/(\w+):([\s\S]+?)\]`)
	str = colorTagRegex.ReplaceAllString(str, "<span style=\"color: #$1;\">$2</span>")
	return str
}
