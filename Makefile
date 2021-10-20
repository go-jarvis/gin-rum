
up:
	cd cmd/demo/ && go run .
tidy:
	go mod tidy

rename:
	find ./ -type f -name "*.go" |xargs sed -i -e 's/goft/rum/g'  -e 's/Goft/Rum/g'
	find ./ -type f -name "*.go" |xargs sed -i -e 's/goft/rum/g'  -e 's/Goft/Rum/g'
	sed -i -e 's/goft/rum/g'  -e 's/Goft/Rum/g' go.mod go.sum
	find ./ -type f -name "*.go" |xargs sed -i -e 's/tangx-labs/go-jarvis/g'  
	sed -i -e 's/tangx-labs/go-jarvis/g' go.mod go.sum
