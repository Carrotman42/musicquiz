
function queryArgString(args) {
	let params = ''
	for (const [key, val] of Object.entries(args)) {
		if (params != '') {
			params += '&';
		}
		params += key+'='+encodeURIComponent(val);
	}
	return params;
}

// Returns a promise which gets resolved after a load event.
// Errors are sent to the reject of the Promise.
function xhrDo(method, path, args) {
        return new Promise((resolve, reject) => {
                const req = new XMLHttpRequest();
                req.addEventListener("error", () => { reject(req); });
                req.addEventListener("abort", () => { req.chowski_aborted = true; reject(req); });
                req.addEventListener("load", () => {
                        if (req.status != 200) {
                                //alert("Failed to contact backend (status: "+req.statusText+")\nError was: "+req.responseText);
                                reject(req);
                        } else {
                                resolve(req);
                        }
                });
		const queryArgs = queryArgString(args);
		if (method == "GET" && queryArgs != "") {
			path += '?' + queryArgs;
		}
                req.open(method, path, true);
		let sendParams = null;
		if (method == "POST") {
			req.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
			sendParams = queryArgs;
		}
                req.send(sendParams);
        });
}
