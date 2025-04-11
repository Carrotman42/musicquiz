"""Reverse-engineering the `ytmusicapi` python library so we can use it in Go.

Usage:
    python ytm.py --client_id="(...).apps.googleusercontent.com" --client_secret="(...)"

Debug output and conclusions documented at the bottom of the file.

The first time running this, it will prompt you for an OAuth Client ID /
Secret; IT MUST BE A TV AND LIMITED INPUT DEVICES CLIENT. Follow instructions
at https://ytmusicapi.readthedocs.io/en/stable/setup/oauth.html and pay
attention to the end of that yellow box at the top!
"""

from ytmusicapi import OAuthCredentials, YTMusic

from http.client import HTTPConnection
import urllib.parse
import logging
import argparse

flags = argparse.ArgumentParser()
flags.add_argument(
        '-c', '--client_id', type=str, required=True,
        help='Google OAuth Client ID (for "TV and Limited Devices"')
flags.add_argument(
        '-s', '--client_secret', type=str, required=True,
        help='Google OAuth Client Secret for the above client')



def main(args, print_api_response=True):
    # Log HTTP requests?!
    HTTPConnection.debuglevel = 1
    logging.basicConfig()
    logging.getLogger().setLevel(logging.DEBUG)
    requests_log = logging.getLogger("requests.packages.urllib3")
    requests_log.setLevel(logging.DEBUG)
    requests_log.propagate = True

    ytm = YTMusic('oauth.json', oauth_credentials=OAuthCredentials(
        client_id=args.client_id, client_secret=args.client_secret))

    playlist_id = 'PLBeRCPzch5o5oiiLzTOUI2f6X70pWo5vZ'
    resp = ytm.get_playlist(playlistId=playlist_id)
    if print_api_response:
        print([track['videoId'] for track in resp['tracks']])

    track_id = 'gthpodDZJeQ'
    resp = ytm.get_song(videoId=track_id)
    if print_api_response:
        # print(resp)
        print([findurl(vid, print_api_response) for vid in resp['streamingData']['formats']])
        # print([findurl(vid) for vid in resp['streamingData']['adaptiveFormats']])
        print(resp['videoDetails']['title'])
        print(resp['videoDetails']['author'])

def findurl(vid, print_api_response):
    sig : str = vid['signatureCipher']
    if print_api_response:
        print(f'found URL in signature {sig} at {sig.index("&url=")}')
        print(f'found keys {vid.keys()}')
    return vid['mimeType'] + '///' + urllib.parse.unquote(sig[sig.index('&url=') + len('&url='):])


if __name__ == '__main__':
    main(flags.parse_args(), print_api_response=False)



#####################
### DEBUG OUTPUT: ###
#####################
#
# Step 1:
#  - Hit https://music.youtube.com with default headers:
#    return {
#        "user-agent": USER_AGENT, # "Mozilla/5.0 (Windows NT 10.0; Win64; rv:88.0) Gecko/20100101 Firefox/88.0"
#        "accept": "*/*",
#        "accept-encoding": "gzip, deflate",
#        "content-type": "application/json",
#        "content-encoding": "gzip",
#        "origin": YTM_DOMAIN, # "https://music.youtube.com"
#    }
#  - Obey the `Set-Cookie` directives in the response, plus adding 'SOCS=CAI'
#  - From the response, grab regex r`ytcfg\.set\s*\(\s*({.+?})\s*\)\s*;`
#    * Parse the match as JSON
#    * Take json['VISITOR_DATA' as the 'X-Goog-Visitor-Id' header
#
# Step 2:
#  - In every post request, include additional 'context' header:
#    return {
#        "context": {
#            "client": {
#                "clientName": "WEB_REMIX",
#                "clientVersion": "1." + time.strftime("%Y%m%d", time.gmtime()) + ".01.00", # e.g. '1.20250411.01.00'
#            },
#            "user": {},
#        }
#    }
#
# DEBUG:urllib3.connectionpool:Starting new HTTPS connection (1): music.youtube.com:443
# send: b'GET / HTTP/1.1\r\nHost: music.youtube.com\r\nuser-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0\r\naccept-encoding: gzip, deflate\r\naccept: */*\r\nConnection: keep-alive\r\ncontent-type: application/json\r\ncontent-encoding: gzip\r\norigin: https://music.youtube.com\r\nCookie: SOCS=CAI\r\n\r\n'
# reply: 'HTTP/1.1 200 OK\r\n'
# header: Content-Type: text/html; charset=utf-8
# header: X-Content-Type-Options: nosniff
# header: Cache-Control: no-cache, no-store, max-age=0, must-revalidate
# header: Pragma: no-cache
# header: Expires: Mon, 01 Jan 1990 00:00:00 GMT
# header: Date: Thu, 10 Apr 2025 22:27:54 GMT
# header: X-Frame-Options: SAMEORIGIN
# header: Strict-Transport-Security: max-age=31536000
# header: Origin-Trial: AmhMBR6zCLzDDxpW+HfpP67BqwIknWnyMOXOQGfzYswFmJe+fgaI6XZgAzcxOrzNtP7hEDsOo1jdjFnVr2IdxQ4AAAB4eyJvcmlnaW4iOiJodHRwczovL3lvdXR1YmUuY29tOjQ0MyIsImZlYXR1cmUiOiJXZWJWaWV3WFJlcXVlc3RlZFdpdGhEZXByZWNhdGlvbiIsImV4cGlyeSI6MTc1ODA2NzE5OSwiaXNTdWJkb21haW4iOnRydWV9
# header: Report-To: {"group":"youtube_music","max_age":2592000,"endpoints":[{"url":"https://csp.withgoogle.com/csp/report-to/youtube_music"}]}
# header: Content-Security-Policy: require-trusted-types-for 'script';report-uri /cspreport
# header: Permissions-Policy: ch-ua-arch=*, ch-ua-bitness=*, ch-ua-full-version=*, ch-ua-full-version-list=*, ch-ua-model=*, ch-ua-wow64=*, ch-ua-form-factors=*, ch-ua-platform=*, ch-ua-platform-version=*
# header: Cross-Origin-Opener-Policy: same-origin-allow-popups; report-to="youtube_music"
# header: Content-Security-Policy-Report-Only: base-uri 'self';default-src 'self' https: blob:;font-src https: data:;img-src https: data: android-webview-video-poster:;media-src blob: https:;object-src 'none';report-uri /cspreport/common;script-src 'nonce-wuuNMf-ib2OLKUb2dTJnJA' 'unsafe-inline' 'strict-dynamic' https: http: 'unsafe-eval';style-src https: 'unsafe-inline'
# header: P3P: CP="This is not a P3P policy! See http://support.google.com/accounts/answer/151657?hl=en for more info."
# header: Content-Encoding: gzip
# header: Server: ESF
# header: X-XSS-Protection: 0
# header: Set-Cookie: YSC=hMJwexSC9Ms; Domain=.youtube.com; Path=/; Secure; HttpOnly; SameSite=none
# header: Set-Cookie: __Secure-YEC=; Domain=.youtube.com; Expires=Fri, 15-Jul-2022 22:27:54 GMT; Path=/; Secure; HttpOnly; SameSite=lax
# header: Set-Cookie: VISITOR_INFO1_LIVE=VZGt1Ud3UGI; Domain=.youtube.com; Expires=Tue, 07-Oct-2025 22:27:54 GMT; Path=/; Secure; HttpOnly; SameSite=none
# header: Set-Cookie: VISITOR_PRIVACY_METADATA=CgJVUxIEGgAgGg%3D%3D; Domain=.youtube.com; Expires=Tue, 07-Oct-2025 22:27:54 GMT; Path=/; Secure; HttpOnly; SameSite=none
# header: Set-Cookie: __Secure-ROLLOUT_TOKEN=CIP2rs_O--2S1QEQ1LzkxsHOjAMY1LzkxsHOjAM%3D; Domain=youtube.com; Expires=Tue, 07-Oct-2025 22:27:54 GMT; Path=/; Secure; HttpOnly; SameSite=none; Partitioned
# header: Alt-Svc: h3=":443"; ma=2592000,h3-29=":443"; ma=2592000
# header: Transfer-Encoding: chunked
# DEBUG:urllib3.connectionpool:https://music.youtube.com:443 "GET / HTTP/1.1" 200 None
# send: b'POST /youtubei/v1/browse?alt=json HTTP/1.1\r\nHost: music.youtube.com\r\nuser-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0\r\naccept-encoding: gzip, deflate\r\naccept: */*\r\nConnection: keep-alive\r\ncontent-type: application/json\r\ncontent-encoding: gzip\r\norigin: https://music.youtube.com\r\nX-Goog-Visitor-Id: CgtWWkd0MVVkM1VHSSjqi-G
# _BjIKCgJVUxIEGgAgGg%3D%3D\r\nauthorization: Bearer (REDACTED)\r\nX-Goog-Request-Time: 1744324075\r\nCookie: YSC=hMJwexSC9Ms; VISITOR_INFO1_LIVE=VZGt1Ud3UGI; VISITOR_PRIVACY_META
# DATA=CgJVUxIEGgAgGg%3D%3D; __Secure-ROLLOUT_TOKEN=CIP2rs_O--2S1QEQ1LzkxsHOjAMY1LzkxsHOjAM%3D; SOCS=CAI\r\nContent-Length: 165\r\n\r\n'
# send: b'{"browseId": "VLPLBeRCPzch5o5oiiLzTOUI2f6X70pWo5vZ", "context": {"client": {"clientName": "WEB_REMIX", "clientVersion": "1.20250410.01.00", "hl": "en"}, "user": {}}}'
# reply: 'HTTP/1.1 200 OK\r\n'
# header: Content-Type: application/json; charset=UTF-8
# header: Vary: Origin
# header: Vary: X-Origin
# header: Vary: Referer
# header: Content-Encoding: gzip
# header: Date: Thu, 10 Apr 2025 22:27:55 GMT
# header: Server: scaffolding on HTTPServer2
# header: X-XSS-Protection: 0
# header: X-Frame-Options: SAMEORIGIN
# header: X-Content-Type-Options: nosniff
# header: Alt-Svc: h3=":443"; ma=2592000,h3-29=":443"; ma=2592000
# header: Transfer-Encoding: chunked
# DEBUG:urllib3.connectionpool:https://music.youtube.com:443 "POST /youtubei/v1/browse?alt=json HTTP/1.1" 200 None
# send: b'POST /youtubei/v1/player?alt=json HTTP/1.1\r\nHost: music.youtube.com\r\nuser-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:88.0) Gecko/20100101 Firefox/88.0\r\naccept-encoding: gzip, deflate\r\naccept: */*\r\nConnection: keep-alive\r\ncontent-type: application/json\r\ncontent-encoding: gzip\r\norigin: https://music.youtube.com\r\nX-Goog-Visitor-Id: CgtWWkd0MVVkM1VHSSjqi-G
# _BjIKCgJVUxIEGgAgGg%3D%3D\r\nauthorization: Bearer (REDACTED)\r\nX-Goog-Request-Time: 1744324075\r\nCookie: YSC=hMJwexSC9Ms; VISITOR_INFO1_LIVE=VZGt1Ud3UGI; VISITOR_PRIVACY_META
# DATA=CgJVUxIEGgAgGg%3D%3D; __Secure-ROLLOUT_TOKEN=CIP2rs_O--2S1QEQ1LzkxsHOjAMY1LzkxsHOjAM%3D; SOCS=CAI\r\nContent-Length: 218\r\n\r\n'
# send: b'{"playbackContext": {"contentPlaybackContext": {"signatureTimestamp": 20188}}, "video_id": "gthpodDZJeQ", "context": {"client": {"clientName": "WEB_REMIX", "clientVersion": "1.20250410.01.00", "hl": "en"}, "user": {}}}'
# reply: 'HTTP/1.1 200 OK\r\n'
# header: Content-Type: application/json; charset=UTF-8
# header: Vary: Origin
# header: Vary: X-Origin
# header: Vary: Referer
# header: Content-Encoding: gzip
# header: Date: Thu, 10 Apr 2025 22:27:55 GMT
# header: Server: scaffolding on HTTPServer2
# header: X-XSS-Protection: 0
# header: X-Frame-Options: SAMEORIGIN
# header: X-Content-Type-Options: nosniff
# header: Alt-Svc: h3=":443"; ma=2592000,h3-29=":443"; ma=2592000
# header: Transfer-Encoding: chunked
# DEBUG:urllib3.connectionpool:https://music.youtube.com:443 "POST /youtubei/v1/player?alt=json HTTP/1.1" 200 None
