# makefile

mod:
	go mod tidy -v                   
	go mod download
	go mod vendor


run:
	@AppBin="${PWD}/vendor/bin"; \
	echo "Creating bin dir at $$AppBin"; \
	mkdir -p "$$AppBin"; \
	\
	if [ ! -f "$$AppBin/air" ]; then \
		echo "Installing air to $$AppBin"; \
		GOBIN="$$AppBin" go install github.com/air-verse/air@latest; \
	fi; \
	\
	$$AppBin/air


build:
	go mod vendor
	go build -mod=vendor -o main


exe:
	go build -o main && ./main


clean:
	rm -vrf "${PWD}/vendor/bin"
	rm -vf "${PWD}"/main
