package formatter

// Portions of this files were taken from https://github.com/sirupsen/logrus
// Copyrighted by Simon Eskildsen under the MIT license

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"math"
	"os"
)

// Defines a log format type that wil output line separated JSON objects
// in the GELF format.
type GelfFormatter struct {
	AppName string
}


func NewGelfFormatter(appName string) *GelfFormatter {
    return  &GelfFormatter{
        AppName: appName,
    }
}

type fields map[string]interface{}

// Format formats the log entry to GELF JSON
func (f *GelfFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(fields, len(entry.Data)+6)
	blacklist := []string{"_id", "id", "timestamp", "version", "level", "application"}

	for k, v := range entry.Data {

		if contains(k, blacklist) {
			continue
		}

		switch v := v.(type) {
		case error:
			// Otherwise errors are ignored by `encoding/json`
			data["_"+k] = v.Error()
		default:
			data["_"+k] = v
		}
	}

	data["version"] = "1.1"
	data["message"] = entry.Message
	data["timestamp"] = round((float64(entry.Time.UnixNano())/float64(1000000))/float64(1000), 4)
	data["level"] = entry.Level
	data["level_name"] = entry.Level.String()
	data["_pid"] = os.Getpid()
	data["application"] = f.AppName

	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}

	return append(serialized, '\n'), nil
}

func contains(needle string, haystack []string) bool {
	for _, a := range haystack {
		if needle == a {
			return true
		}
	}
	return false
}

func round(val float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor((val*shift)+.5) / shift
}
