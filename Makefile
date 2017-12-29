RACEPATH = "log_path=./report strip_path_prefix=pkg/"

all: httpserver 

httpserver:
	cd bin && GOPATH=$(PWD) go build -gcflags "-N -l" httpserver 

