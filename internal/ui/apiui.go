package ui

import "github.com/rancher/apiserver/pkg/writer"

type APIUI struct {
	offline StringSetting
	release BoolSetting
	embed   bool
}

func apiUI(opt *Options) APIUI {
	var rtn = APIUI{
		offline: opt.Offline,
		release: opt.ReleaseSetting,
		embed:   true,
	}
	if rtn.offline == nil {
		rtn.offline = StaticSetting("dynamic")
	}
	if rtn.release == nil {
		rtn.release = StaticSetting(false)
	}
	for _, file := range []string{
		"ui/api-ui/ui.min.css",
		"ui/api-ui/ui.min.js",
	} {
		if _, err := staticContent.Open(file); err != nil {
			rtn.embed = false
			break
		}
	}
	return rtn
}

func (a APIUI) content(name string) writer.StringGetter {
	return func() (rtn string) {
		switch a.offline() {
		case "dynamic":
			if !a.release() && !a.embed {
				return ""
			}
		case "false":
			return ""
		}
		return name
	}
}

func (a APIUI) CSS() writer.StringGetter {
	return a.content("/api-ui/ui.min.css")
}

func (a APIUI) JS() writer.StringGetter {
	return a.content("/api-ui/ui.min.js")
}
