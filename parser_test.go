package parsegen

import (
	"io"
	"os"
	"testing"
)

// TestSDPParser tests generated parser based on ./sdp.bnf.
// Our goal is to pull out ports and ip addresses from SIP-SDP call between
// two callers.
// More about SDP see in rfc4566
func TestSDPParser(t *testing.T) {
	f, err := os.Open("sdp.bnf")
	assert(t, err == nil)
	bnf, err := io.ReadAll(f)
	assert(t, err == nil)

	parser, err := Generate(bnf)
	assert(t, err == nil, err)

	input := "v=0\r\n" +
		"o=- 345678 345979 IN IP4 10.0.1.2\r\n" +
		"s=My sample redundant flow\r\n" +
		"i=2 channels: c6, c7\r\n" +
		"t=0 0\r\n" +
		"a=recvonly\r\n" +
		"a=group:DUP prim sec\r\n" +
		"m=audio 5004 RTP/AVP 98\r\n" +
		"c=IN IP4 239.69.22.33/32\r\n" +
		"a=rtpmap:98 L24/48000/2\r\n" +
		"a=ptime:1\r\n" +
		"a=ts-refclk:ptp=IEEE1588-2008:00-11-22-FF-FE-33-44-55:0\r\n" +
		"a=mediaclk:direct=0\r\n" +
		"a=mid:prim\r\n" +
		"m=audio 5004 RTP/AVP 98\r\n" +
		"c=IN IP4 239.69.22.33/32\r\n" +
		"a=rtpmap:98 L24/48000/2\r\n" +
		"a=ptime:1\r\n" +
		"a=ts-refclk:ptp=IEEE1588-2008:00-11-22-FF-FE-33-44-55:0\r\n" +
		"a=mediaclk:direct=0\r\n" +
		"a=mid:prim\r\n" +
		"m=audio 5004 RTP/AVP 98\r\n" +
		"c=IN IP4 239.69.44.55/32\r\n" +
		"a=rtpmap:98 L24/48000/2\r\n" +
		"a=ptime:1\r\n" +
		"a=ts-refclk:ptp=IEEE1588-2008:00-11-22-FF-FE-33-44-55:0\r\n" +
		"a=mediaclk:direct=0\r\n" +
		"a=mid:sec\r\n"

	pd, err := parser.Parse([]byte(input))
	assert(t, err == nil, err)

	// check first port
	// Get returns first instance of Port
	assert(t, string(pd.Get("Port")) == "5004")
	assert(t, string(pd.Get("Connection-address")) == "239.69.22.33/32")
}
