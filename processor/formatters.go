package processor

import (
	"fmt"
	"os"
	"time"
)

// Get the time as standard UTC/Zulu format
func getFormattedTime() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// Prints a message to stdout if flag to enable warning output is set
func printVerbose(msg string) {
	if Verbose {
		fmt.Println(fmt.Sprintf("VERBOSE %s: %s", getFormattedTime(), msg))
	}
}

// Prints a message to stdout if flag to enable debug output is set
func printDebug(msg string) {
	if Debug {
		fmt.Println(fmt.Sprintf("DEBUG %s: %s", getFormattedTime(), msg))
	}
}

// Used when explicitly for os.exit output when crashing out
func printError(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, "ERROR %s: %s", getFormattedTime(), msg)
}



func fileSummarize(input chan Result) string {
	//switch {
	//case More || strings.ToLower(Format) == "wide":
	//	return fileSummarizeLong(input)
	//case strings.ToLower(Format) == "json":
	//	return toJSON(input)
	//case strings.ToLower(Format) == "csv":
	//	return toCSV(input)
	//}

	return ""
}