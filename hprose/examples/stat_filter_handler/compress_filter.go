package main

import (
	"compress/gzip"

	"io/ioutil"

	"github.com/JexLib/golang/hprose/io"
	"github.com/JexLib/golang/hprose/rpc"
)

// CompressFilter ...
type CompressFilter struct{}

// InputFilter ...
func (CompressFilter) InputFilter(data []byte, context rpc.Context) []byte {
	b := io.NewByteReader(data)
	reader, _ := gzip.NewReader(b)
	defer reader.Close()
	data, _ = ioutil.ReadAll(reader)
	return data
}

// OutputFilter ...
func (CompressFilter) OutputFilter(data []byte, context rpc.Context) []byte {
	b := &io.ByteWriter{}
	writer := gzip.NewWriter(b)
	writer.Write(data)
	writer.Flush()
	return b.Bytes()
}
