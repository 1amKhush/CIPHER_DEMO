package controllers

import (
	"bytes"
	"errors"
)

func AnalyseRequest(data []byte) (method string, fname string, payload []byte, err error) {
	// split the data 
	parts := bytes.SplitN(data, []byte{0x00}, 3)

	// requred params available ? 
	if len(parts) < 3 {
		return "", "", nil, errors.New("malformed request protocol: missing fields or delimiters")
	}

	method = string(parts[0])
	fname = string(parts[1])
	payload = parts[2]

	// 4. Basic sanity validation
	if method == "" {
		return "", "", nil, errors.New("protocol error: empty method")
	}
	return method, fname, payload, nil
}

func FormatRequest(method, fname string, payload []byte) []byte {
	return bytes.Join([][]byte{[]byte(method), []byte(fname), payload}, []byte{0x00})
}


// what next 
