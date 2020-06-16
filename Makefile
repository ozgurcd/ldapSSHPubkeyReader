
TARGET_ARCH=amd64
TARGET_OS=darwin
#TARGET_OS=linux
OUTPUT=ldapPubKeyReader
LDFLAGS=-ldflags '-w -s'


all: build

build:
	GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build $(LDFLAGS) -o $(OUTPUT)


static:
	env CGO_ENABLED=0 GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build $(LDFLAGS) -o $(OUTPUT)

