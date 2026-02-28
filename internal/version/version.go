package version

import "runtime/debug"

// These variables are set via ldflags during CI/goreleaser builds.
var (
	version = ""
	commit  = ""
	date    = ""
)

// GetVersion returns the version string.
// Priority: ldflags > go install build info > "dev"
func GetVersion() string {
	if version != "" {
		return version
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		if v := info.Main.Version; v != "" && v != "(devel)" {
			return v
		}
	}

	return "dev"
}

// GetCommit returns the short commit hash.
func GetCommit() string {
	if commit != "" {
		return commit
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				if len(setting.Value) > 7 {
					return setting.Value[:7]
				}

				return setting.Value
			}
		}
	}

	return "unknown"
}

// GetDate returns the build date.
func GetDate() string {
	if date != "" {
		return date
	}

	return "unknown"
}
