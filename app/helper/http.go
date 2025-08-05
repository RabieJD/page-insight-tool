package helper

var statusExplanations = map[int]string{
	// 1xx - Informational
	100: "The link host server has received the request headers and the client should proceed to send the request body.",
	101: "The link host server is switching protocols as requested by the client.",
	102: "The link host server has received and is processing the request, but no response is available yet.",
	103: "The link host server is sending some headers before the final response.",

	// 2xx - Success
	200: "The request succeeded.",
	201: "The request succeeded and a new resource was created.",
	202: "The request has been accepted for processing, but the processing is not complete.",
	203: "The request was successful, but the enclosed metadata may come from a third party.",
	204: "The link host server successfully processed the request, but is not returning any content.",
	205: "The link host server successfully processed the request and asks the client to reset the document view.",
	206: "The link host server is delivering only part of the resource due to a range header sent by the client.",

	// 3xx - Redirection
	300: "The requested resource has multiple choices (e.g., different formats).",
	301: "The resource has been moved permanently to a new URL.",
	302: "The resource has been temporarily moved to a different URL.",
	303: "The link host server is redirecting to a different resource, typically after a POST.",
	304: "The resource has not been modified since the last request.",
	305: "The requested resource must be accessed through the proxy given by The link host server.",
	307: "The request should be repeated with another URL, but future requests should still use the original URL.",
	308: "The resource is permanently redirected to a new URL.",

	// 4xx - Client Errors
	400: "The link host server could not understand the request due to invalid syntax.",
	401: "Authentication is required and has failed or has not been provided.",
	403: "The link host server understood the request but refuses to authorize it.",
	404: "The link host server can't find the requested resource.",
	405: "The request method is not allowed for the requested resource.",
	406: "The link host server cannot produce a response matching the list of acceptable values defined in the request's headers.",
	407: "Authentication is required via a proxy.",
	408: "The link host server timed out waiting for the request.",
	409: "The request conflicts with the current state of The link host server.",
	410: "The requested resource is no longer available and will not be available again.",
	411: "The request did not specify the length, which is required by the resource.",
	412: "The link host server does not meet one of the preconditions specified by the client.",
	413: "The request is larger than The link host server is willing or able to process.",
	414: "The URI provided was too long for The link host server to process.",
	415: "The media format of the request is not supported by The link host server.",
	416: "The requested range is not satisfiable.",
	417: "The link host server cannot meet the expectations specified in the request headers.",
	418: "I'm a teapot.",
	422: "The request was well-formed but could not be followed due to semantic errors.",
	426: "The client should switch to a different protocol.",
	429: "Too many requests have been sent in a given amount of time.",
	431: "The link host server is unwilling to process the request because its header fields are too large.",

	// 5xx - Server Errors
	500: "The link host server encountered an unexpected condition that prevented it from fulfilling the request.",
	501: "The link host server does not support the functionality required to fulfill the request.",
	502: "The link host server received an invalid response from an upstream server.",
	503: "The link host server is not ready to handle the request (e.g., down for maintenance).",
	504: "The link host server did not receive a timely response from an upstream server.",
	505: "The HTTP version used in the request is not supported by The link host server.",
	507: "The link host server is unable to store the representation needed to complete the request.",
	511: "Network authentication is required to access the resource.",
}

func GetExplanation(code int) string {
	explanation, _ := statusExplanations[code]
	return explanation
}
