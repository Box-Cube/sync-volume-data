APP_DIR = ./bin

all:
	if [ ! -d ${APP_DIR} ];then mkdir -p bin; fi
	CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-s -w -extldflags "-static"' -o bin/sync-pod-data && upx -5 bin/sync-pod-data