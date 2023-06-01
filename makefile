.PHONY: default fmt
.DEFAULT_GOAL := default

default:
	@$(eval VER := $(shell go run *.go -v |  awk '{print $$3}'))
	@$(eval DIRLIN := dist/linux/amd64/${VER})
	@$(eval DIRWIN := dist/windows/amd64/${VER})
	@rm -rf dist
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ${DIRLIN}/logsync .
	@(cd ${DIRLIN} && md5sum logsync > logsync.md5)
	@(cd ${DIRLIN} && sha256sum logsync > logsync.sha256)
	@(cd ${DIRLIN} && zip -qr linux_amd64_logsync-${VER}.zip logsync logsync.md5 logsync.sha256)
	@(cd ${DIRLIN} && md5sum linux_amd64_logsync-${VER}.zip > linux_amd64_logsync-${VER}.zip.md5)
	@(cd ${DIRLIN} && sha256sum linux_amd64_logsync-${VER}.zip > linux_amd64_logsync-${VER}.zip.sha256)
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o ${DIRWIN}/logsync.exe .
	@(cd ${DIRWIN} && md5sum logsync.exe > logsync.exe.md5)
	@(cd ${DIRWIN} && sha256sum logsync.exe > logsync.exe.sha256)
	@(cd ${DIRWIN} && zip -qr windows_amd64_logsync-${VER}.zip logsync.exe logsync.exe.md5 logsync.exe.sha256)
	@(cd ${DIRWIN} && md5sum windows_amd64_logsync-${VER}.zip > windows_amd64_logsync-${VER}.zip.md5)
	@(cd ${DIRWIN} && sha256sum windows_amd64_logsync-${VER}.zip > windows_amd64_logsync-${VER}.zip.sha256)

fmt:
	@golangci-lint run
