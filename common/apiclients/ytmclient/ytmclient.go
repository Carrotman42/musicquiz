package ytmclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type GetSongRequest struct {
	PlaybackContext PlaybackContext `json:"playbackContext"`
	VideoID         string `json:"video_id"`
}

type PlaybackContext struct {
	ContentPlaybackContext struct {
		SignatureTimestamp int `json:"signatureTimestamp"`
	} `json:"contentPlaybackContext"`
}

const auth = ""

type GetSongResponse map[any]any

const apiBase = "https://music.youtube.com/youtubei/v1/"

func New(client *http.Client) *Client {
	return &Client{client}
}

type Client struct {
	http *http.Client
}

// the http response will probably be cut off if you cancel the context.
func (c *Client) post(ctx context.Context, endpoint string, urlParams url.Values, body io.Reader) (*http.Response, error) {
	u := apiBase + endpoint + "?alt=json"
	if len(urlParams) > 0 {
		u += "&" + urlParams.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, "POST", u, body)
	if err != nil {
		panic(err)
	}
		req.Header = http.Header{
			//"Accept-Encoding": {"gzip, deflate, br, zstd"},
			//"Accept-Language":               {"en-US,en;q=0.5"},
			"Alt-Used":                      {"music.youtube.com"},
			//"Authorization":                 {"SAPISIDHASH 1744313888_0da62b9a005187d5c0e6189188c4203daa7e6084_u SAPISID1PHASH 1744313888_0da62b9a005187d5c0e6189188c4203daa7e6084_u SAPISID3PHASH 1744313888_0da62b9a005187d5c0e6189188c4203daa7e6084_u"},
			//"Connection":                    {"keep-alive"},
			"Content-Type":                  {"application/json"},
			"Cookie": {"SOCS=CAI"},
			//"Cookie":                        {"VISITOR_INFO1_LIVE=9Z9eMpl1rq4; VISITOR_PRIVACY_METADATA=CgJVUxIEGgAgQQ%3D%3D; PREF=f6=80&f7=100&tz=America.Denver&repeat=NONE&autoplay=false&has_user_changed_default_autoplay_mode=true&f5=30000&guide_collapsed=false&volume=65; __Secure-1PSIDTS=sidts-CjIB7pHptdCxZYllHUGOT_YLkyUX8LxzJPhJIBfyI6w7bKzEsQK4MssQG-BIWAhSFP2nlxAA; __Secure-3PSIDTS=sidts-CjIB7pHptdCxZYllHUGOT_YLkyUX8LxzJPhJIBfyI6w7bKzEsQK4MssQG-BIWAhSFP2nlxAA; HSID=AoJgTOgyCjM7MgIJC; SSID=ANeb9JfhIKoiV-zVt; APISID=t1HoWEFIMWV8rqAX/ASShvoPdT9DYxffmp; SAPISID=bTyLQGLj3gCTsZfZ/AOuTuULXniyPyJRUG; __Secure-1PAPISID=bTyLQGLj3gCTsZfZ/AOuTuULXniyPyJRUG; __Secure-3PAPISID=bTyLQGLj3gCTsZfZ/AOuTuULXniyPyJRUG; SID=g.a000ugjEwy54zm1h7PJTtYutB3nRRK5WXkeB4nBGcZmku1SxbFp3_4t9nRmc8wRWIpXBPONvGwACgYKAc0SARMSFQHGX2MiMTzfSIjJzoQwEgfKigVyQhoVAUF8yKqdfOuv1C3xar0HHuNBdBDP0076; __Secure-1PSID=g.a000ugjEwy54zm1h7PJTtYutB3nRRK5WXkeB4nBGcZmku1SxbFp3q3hrocKkrvglsCaJRRQtnQACgYKAfsSARMSFQHGX2MiTaZhKeEaRV31UMUpqIRZ2hoVAUF8yKp5fcVEMafb5opipes1mymP0076; __Secure-3PSID=g.a000ugjEwy54zm1h7PJTtYutB3nRRK5WXkeB4nBGcZmku1SxbFp3XvHCxGnRMJydt7qV7SSVPQACgYKAekSARMSFQHGX2MiLgR8qYVhXgTBlKeXEG5PFxoVAUF8yKq73hFuc6LkKnAL0GpliuY80076; LOGIN_INFO=AFmmF2swRQIhAIUBcDl8QqjwwAYTwIgU2db4b2SRZu10w_Jy_h5oBbFOAiB13J50H-MVY57RcVRpDS_M68UeZ7pV21HsiVf2T6jDjA:QUQ3MjNmeDdCWFZqQ1FTLW5keHFCeGlMcWhQbTg2aUlBRWdKalRvWERyNENTOUtaREZ4OHBkZ0pJcEJOQ1FyZ24zSlc3QzRiUWRpOEZZdEdGY1NEbzJrRWlnX29OcDBENy01Tl9Yb0hqNVRmakVvczRkTDJFclJNSnQzbDM3QzFTY3JKam1qak44SEFxUHNSTThPVmc0WXVGRU94Mzk5QXRB; SIDCC=AKEyXzWMWtcK6wiUtBXPZLlWYL0v_qR3om7EHQULTDf4nUa2v1uO6TQm2L6FPMMMY8cnuaueqnRK; __Secure-1PSIDCC=AKEyXzUpXJZoRGP7rwlvvpLUiP7Z2HMAKhxYiRKm4xhmBhYmAApOA12rSJjgq9qdjynfu5OPf1s; __Secure-3PSIDCC=AKEyXzWJRUUGTh5YJ5I76mDLlRer7R8qCEsGvfY-kPjVMM7AzMowwYtSldMQrucerLQe4P7z1Cd1; __Secure-ROLLOUT_TOKEN=CJaE-PrlreCsygEQ-Zmv1YXsigMY7bDnwprOjAM%3D; YSC=080Kr5Oi5B0"},
			"Host":                          {"music.youtube.com"},
			"Origin":                        {"https://music.youtube.com"},
			"Referer":                       {"https://music.youtube.com/watch?v=h_r1CR6Q8z0"},
			"Sec-Fetch-Dest":                {"empty"},
			"Sec-Fetch-Mode":                {"cors"},
			"Sec-Fetch-Site":                {"same-origin"},
			"TE":                            {"trailers"},
			"User-Agent":                    {"Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0"},
			"X-Goog-AuthUser":               {"0"},
			"X-Goog-Visitor-Id":             {"Cgs5WjllTXBsMXJxNCiVvOC_BjIKCgJVUxIEGgAgQQ%3D%3D"},
			"X-Origin":                      {"https://music.youtube.com"},
			"X-Youtube-Bootstrap-Logged-In": {"true"},
			"X-Youtube-Client-Name":         {"67"},
			"X-Youtube-Client-Version":      {"1.20250407.01.00"},
		}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		bs, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("%v: %v; body: %s\nHeaders: %v", u, resp.Status, bs, resp.Header)
	}
	return resp, nil
}

func (c *Client) GetSong(ctx context.Context, id string) (GetSongResponse, error) {
	req := GetSongRequest{
		VideoID: id,
	}
	req.PlaybackContext.ContentPlaybackContext.SignatureTimestamp = int(time.Since(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))/(24*time.Hour))

	body, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	resp, err := c.post(ctx, "player", nil, bytes.NewReader(body))
	if err != nil {
		return GetSongResponse{}, fmt.Errorf("post failed: %v", err)
	}
	defer resp.Body.Close()
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetSongResponse{}, fmt.Errorf("read post result failed: %v", err)
	}
	var ret GetSongResponse
	if err := json.Unmarshal(bs, &ret); err != nil {
		return GetSongResponse{}, fmt.Errorf("unmarshal: %v", err)
	}
	return ret, nil
}
