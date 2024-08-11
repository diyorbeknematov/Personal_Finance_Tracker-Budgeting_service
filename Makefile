include .env
export $(shell sed 's/=.*//' .env)

#CURRENT_DIR ni yaratish
CURRENT_DIR := $(shell pwd)

gen-proto:
	./scripts/gen-proto.sh $(CURRENT_DIR)


#dasturni run qilish
run:
	go run cmd/main.go
	