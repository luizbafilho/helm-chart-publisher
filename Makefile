extension = $(patsubst windows,.exe,$(filter windows,$(1)))
PKG_NAME := helm-chart-publisher
define gocross
	GOOS=$(1) GOARCH=$(2) go build -o ./dist/$(PKG_NAME)_$(1)-$(2)$(call extension,$(1));
endef

release:
	ghr -draft -u luizbafilho $$TAG ./dist

clean:
	rm -Rf ./dist/*

build-all: clean
	$(call gocross,linux,amd64)
	$(call gocross,linux,arm)
	$(call gocross,darwin,amd64)
	$(call gocross,windows,amd64)

test:
	go test -v -race `glide novendor`

ghr:
	go get -u github.com/tcnksm/ghr

docker: build-all
	docker build -t $(PKG_NAME) .
