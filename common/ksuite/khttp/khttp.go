package khttp

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"net/http"
	"time"
)

func ServeStaticBytes(bs []byte) http.HandlerFunc {
	hasher := fnv.New64a()
	hasher.Write(bs)
	etag := fmt.Sprintf(`"ssb.fnv64a-%08X"`, hasher.Sum64())
	startTime := time.Now()

	return func(wr http.ResponseWriter, req *http.Request) {
		wr.Header().Set("Cache-Control", "must-revalidate, max-age=604800") // 7 days
		wr.Header().Set("Etag", etag)
		http.ServeContent(wr, req, "", startTime, bytes.NewReader(bs))
	}
}
