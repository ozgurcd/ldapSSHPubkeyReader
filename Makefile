
TARGET_ARCH=amd64
#TARGET_OS=darwin
TARGET_OS=linux
OUTPUT=ldapPubkeyReader
LDFLAGS=-ldflags '-w'


all: build

build:
	GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build $(LDFLAGS) -o $(OUTPUT)


static:
	env CGO_ENABLED=0 GOOS=$(TARGET_OS) GOARCH=$(TARGET_ARCH) go build $(LDFLAGS) -o $(OUTPUT)

