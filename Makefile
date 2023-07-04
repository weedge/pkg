# This repo's root import path (under $GOPATH)
KITEX_MODULE := github.com/weedge/pkg

.PHONY: kitex
kitex:
	@rm -rf ./kitex_gen/
	@kitex -module $(KITEX_MODULE) ./idl/thrift/base.thrift
