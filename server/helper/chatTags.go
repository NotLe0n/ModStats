package helper

import (
	"html"
	"regexp"
	"strings"
)

func Unescape(str string) string {
	str = html.UnescapeString(str)

	replaceMap := map[string]string{
		"\r\n": "<br>",
		"\n":   "<br>",
		"\t":   "    ",
		"\\'":  "'",
		"\\\\": "\\",
	}

	for oldStr, newStr := range replaceMap {
		str = strings.ReplaceAll(str, oldStr, newStr)
	}

	return str
}

func ParseChatTags(str string) string {
	itemTagRegex := regexp.MustCompile(`\[i(.*?):(\w+)\]`)
	str = itemTagRegex.ReplaceAllString(str, "<img src=\"https://tmlapis.tomat.dev/img/Item_$2.png\" id=\"item-icon\">")
	colorTagRegex := regexp.MustCompile(`\[c\/(\w+):([\s\S]+?)\]`)
	str = colorTagRegex.ReplaceAllString(str, "<span style=\"color: #$1;\">$2</span>")
	return str
}
