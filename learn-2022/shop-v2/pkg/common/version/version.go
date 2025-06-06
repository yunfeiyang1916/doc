// Package version supplies version information collected at build time to
// apimachinery components.
package version

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/gosuri/uitable"
)

var (
	// GitVersion is semantic version.
	GitVersion = "v0.0.0-master+$Format:%h$"
	// 构建日期 in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ').
	BuildDate = "1970-01-01T00:00:00Z"
	// 最后一次提交的 Git SHA1 值, output of $(git rev-parse HEAD).
	GitCommit = "$Format:%H$"
	// Git 工作树状态，默认为空，可为 "clean" 或 "dirty"
	GitTreeState = ""
)

// 版本信息
type Info struct {
	// Git版本号
	GitVersion string `json:"gitVersion"`
	// 最后一次提交的 Git SHA1 值,
	GitCommit string `json:"gitCommit"`
	// Git 工作树状态
	GitTreeState string `json:"gitTreeState"`
	// 构建日期
	BuildDate string `json:"buildDate"`
	// Go 版本
	GoVersion string `json:"goVersion"`
	// 编译器信息
	Compiler string `json:"compiler"`
	// 目标平台信息
	Platform string `json:"platform"`
}

// String returns info as a human-friendly version string.
func (info Info) String() string {
	if s, err := info.Text(); err == nil {
		return string(s)
	}

	return info.GitVersion
}

// ToJSON returns the JSON string of version information.
func (info Info) ToJSON() string {
	s, _ := json.Marshal(info)

	return string(s)
}

// Text encodes the version information into UTF-8-encoded text and
// returns the result.
func (info Info) Text() ([]byte, error) {
	table := uitable.New()
	table.RightAlign(0)
	table.MaxColWidth = 80
	table.Separator = " "
	table.AddRow("gitVersion:", info.GitVersion)
	table.AddRow("gitCommit:", info.GitCommit)
	table.AddRow("gitTreeState:", info.GitTreeState)
	table.AddRow("buildDate:", info.BuildDate)
	table.AddRow("goVersion:", info.GoVersion)
	table.AddRow("compiler:", info.Compiler)
	table.AddRow("platform:", info.Platform)

	return table.Bytes(), nil
}

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() Info {
	// These variables typically come from -ldflags settings and in
	// their absence fallback to the settings in pkg/version/base.go
	return Info{
		GitVersion:   GitVersion,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
