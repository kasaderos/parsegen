package main

import (
	"io"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	f, err := bnfFunction(nil)
	assert(t, err == nil, err)
	printTree(f)
}

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
	pd.Data().Print()

	pd, err = parser.Parse([]byte(input))
	assert(t, err == nil, err)
	pd.Data().Print()
}
