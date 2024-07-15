package api

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mileusna/useragent"
)

// Predefined dictionary with browser versions and their corresponding utls values.
// Updated as of utls v1.6.6.
// Extracted from here: https://github.com/refraction-networking/utls/blob/master/u_common.go#L573
var utlsDict = map[string]map[int]string{
	"Firefox": {
		-1:  "55",
		56:  "56",
		63:  "63",
		65:  "65",
		99:  "99",
		102: "102",
		105: "105",
		120: "120",
	},
	"Chrome": {
		-1:  "58",
		62:  "62",
		70:  "70",
		72:  "72",
		83:  "83",
		87:  "87",
		96:  "96",
		100: "100",
		102: "102",
		106: "106",
		112: "112_PSK",
		114: "114_PSK",
		120: "120",
	},
	"iOS": {
		-1: "111",
		12: "12.1",
		13: "13",
		14: "14",
	},
	"Android": {
		-1: "11",
	},
	"Edge": {
		-1: "85",
		// 106: "106", incompatible with utls
	},
	"Safari": {
		-1: "16.0",
	},
	"360Browser": {
		-1: "7.5",
		// 11: "11.0", incompatible with utls
	},
	"QQBrowser": {
		-1: "11.1",
	},
}

func uagentToUtls(uagent string) (string, string, error) {
	ua := useragent.Parse(uagent)
	utlsVersion, err := utlsVersion(ua.Name, ua.Version)
	if err != nil {
		return "", "", err
	}
	return ua.Name, utlsVersion, nil
}

func utlsVersion(browserName, browserVersion string) (string, error) {
	if versions, ok := utlsDict[browserName]; ok {
		// Extract the major version number from the browser version string
		majorVersionStr := strings.Split(browserVersion, ".")[0]
		majorVersion, err := strconv.Atoi(majorVersionStr)
		if err != nil {
			return "", fmt.Errorf("error parsing major version number from browser version: %v", err)
		}

		// Find the highest version that is less than or equal to the browser version
		var selectedVersion int = -1
		for version := range versions {
			if version <= majorVersion && version > selectedVersion {
				selectedVersion = version
			}
		}

		if utls, ok := versions[selectedVersion]; ok {
			return utls, nil
		} else {
			return "", fmt.Errorf("no UTLS value found for browser '%s' with version '%s'", browserName, browserVersion)
		}
	}
	return "", fmt.Errorf("browser '%s' not found in UTLS dictionary", browserName)
}
