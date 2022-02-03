build_linux:
	CGO_ENABLED=1 GOOS=linux go build -buildmode=c-shared -o libflagship.so

build_mac:
	CGO_ENABLED=1 GOOS=darwin go build -buildmode=c-shared -o libflagship.dylib

build_windows:
	CGO_ENABLED=1 GOOS=windows go build -buildmode=c-shared -o libflagship.dylib