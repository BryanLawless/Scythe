package utility

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	ColorReset = "\033[0m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Purple     = "\033[35m"
	Cyan       = "\033[36m"
	White      = "\033[37m"
	Black      = "\u001b[30;1m"
	BlackBg    = "\u001b[40;1m"
)

func GetTimestamp() string {
	return time.Now().Format("15:04:05")
}

func GetQuote() string {
	quotes := [3]string{
		"For data hoarders <3",
		"Pretty darn fast, if I do say so myself",
		"GAS GAS GAS, Im gonna step on the GAS",
	}

	return quotes[rand.Intn(len(quotes))]
}

func GetBanner() string {
	banner := `
███████╗ ██████╗██╗   ██╗████████╗██╗  ██╗███████╗  ┏━━━━━━━━━━━━━━━━━━ Info ━━━━━━━━━━━━━━━━┓
██╔════╝██╔════╝╚██╗ ██╔╝╚══██╔══╝██║  ██║██╔════╝    ` + Cyan + `@ Package:` + ColorReset + ` Scythe
███████╗██║      ╚████╔╝    ██║   ███████║█████╗      ` + Cyan + `@ Author:` + ColorReset + ` TheLawlessDev
╚════██║██║       ╚██╔╝     ██║   ██╔══██║██╔══╝      ` + Cyan + `@ License:` + ColorReset + ` APL-GPL 3.0
███████║╚██████╗   ██║      ██║   ██║  ██║███████╗    ` + Cyan + `@ Github:` + ColorReset + ` /TheLawlessDev/Scythe  ` + Green + `
╚══════╝ ╚═════╝   ╚═╝      ╚═╝   ╚═╝  ╚═╝╚══════╝  ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
	`

	return strings.ReplaceAll(banner, "█", White+"█"+Green)
}

func WriteToConsole(status string, mode string) {
	switch mode {
	case "info":
		fmt.Fprintf(color.Output, "%s%s[INFO]%s %s%s%s\n", BlackBg, Cyan, ColorReset, White, status, ColorReset)
	case "warning":
		fmt.Fprintf(color.Output, "%s%s[WARNING]%s %s%s%s\n", BlackBg, Yellow, ColorReset, White, status, ColorReset)
	case "success":
		fmt.Fprintf(color.Output, "%s%s[SUCCESS]%s %s%s%s\n", BlackBg, Green, ColorReset, White, status, ColorReset)
	case "error":
		fmt.Fprintf(color.Output, "%s%s[ERROR]%s %s%s%s\n", BlackBg, Red, ColorReset, White, status, ColorReset)
	}
}
