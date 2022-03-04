.PHONY: run
run:
	go run .

.PHONY: build
build:
	go build -o ./build/block .

.PHONY: test
createblockchain:
	build/block createblockchain -address ningqing
	build/block getbalance -address ningqing
send:
	build/block send -from ningqing -to vonevone -amount 6
getbalance:
	build/block getbalance -address ningqing
	build/block getbalance -address vonevone