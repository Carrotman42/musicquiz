package kjs

import (
	_ "embed"

	"chowski3/common/ksuite/khttp"
)

//go:embed xhr.js
var XHR []byte

var XHRHandler = khttp.ServeStaticBytes(XHR)
