package bag

type key struct {
	key string
}

// context* contains the list of keys that are used for putting items into and
// fetching from the context.
var (
	contextIP          = key{key: "IP"}
	contextAccessToken = key{key: "accessToken"}
	contextAPIRoot     = key{key: "apiRoot"}
	contextSite        = key{key: "site"}
	contextProfile     = key{key: "profile"}
)
