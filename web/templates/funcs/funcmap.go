package funcs

import (
	"html/template"

	common "github.com/buro9/funcs"
)

// Namespace is a prefix applied to a func in the `template.FuncMap`. The
// purpose is to differentiate application specific funcs from the common funcs
// so that it is easier to find and debug them, and where funcs exist in both
// the common and application we can prevent overwriting and also be more explicit
// about which we are calling from templates.
const Namespace = "microcosm_"

// FuncMap is the template.FuncMap used by this application and includes both
// application specific funcs as well as common funcs.
var FuncMap template.FuncMap

// buildFuncMap is called by init, and takes the common funcs provided by
// github.com/buro9/funcs and then combines them with application specific funcs
// to produce the template.FuncMap that our application will use.
func buildFuncMap() template.FuncMap {
	funcMap := common.Map

	funcMap[Namespace+"api2ui"] = Api2ui
	funcMap[Namespace+"url"] = Url
	funcMap[Namespace+"avatar_url"] = AvatarURL
	funcMap[Namespace+"isImage"] = IsImage
	funcMap[Namespace+"link"] = LinkFromLinks
	funcMap[Namespace+"reverseLinks"] = ReverseLinks
	funcMap[Namespace+"stat"] = Stat

	return funcMap
}
