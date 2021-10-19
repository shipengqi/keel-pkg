package fmtutil

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

// Color defines a single SGR Code
type Color int

// Foreground text colors
const (
	FgBlack Color = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

func PrettyPrefix(msg string) string {
	return prettyPrefix(msg).String()
}

func Pretty(msg, status string) string {
	b := prettyPrefix(msg)
	b.WriteString(status)
	b.WriteString(" ]")
	return b.String()
}

func prettyPrefix(msg string) *strings.Builder {
	var b strings.Builder
	ml := len(msg)
	pl := 64 - ml - 1
	padding := strings.Repeat(".", pl)
	b.WriteString("    ")
	b.WriteString(msg)
	b.WriteString(" ")
	b.WriteString(padding)
	b.WriteString(" ")
	b.WriteString("[ ")
	return &b
}

func PrettyWithColor(msg, status string, c Color) string {
	b := prettyPrefix(msg)
	b.WriteString(Colorize(status, c))
	b.WriteString(" ]")
	return b.String()
}

// Colorize a string based on given color.
func Colorize(s string, c Color) string {
	return fmt.Sprintf("\033[1;%s;40m%s\033[0m", strconv.Itoa(int(c)), s)
}

func VersionTmpl(cmd, version, commit, time string) string {
	return fmt.Sprintf("%s:\n%s\n%s\n%s\n%s\n%s\n%s\n",
		cmd,
		prettyVersionLine("Version", version),
		prettyVersionLine("CommitID", commit),
		prettyVersionLine("BuildTime", time),
		prettyVersionLine("GoVersion", runtime.Version()),
		prettyVersionLine("Compiler", runtime.Compiler),
		prettyVersionLine("Platform", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)))
}

func prettyVersionLine(msg, value string) string {
	var b strings.Builder
	ml := len(msg)
	pl := 12 - ml - 1
	padding := strings.Repeat(" ", pl)
	b.WriteString("  ")
	b.WriteString(msg)
	b.WriteString(": ")
	b.WriteString(padding)
	b.WriteString(value)
	return b.String()
}
