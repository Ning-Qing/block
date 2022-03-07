.PHONY: run
run:
	go run .

.PHONY: build
build:
	go build -o ./build/block .

.PHONY: clear
clear:
	rm -rf wallet.dat block.db

.PHONY: test
createwallet:
	build/block createwallet
createblockchain:
	build/block createblockchain -address $(address)
	build/block getbalance -address $(address)
send:
	build/block send -from $(from) -to $(to) -amount $(amount)
getbalance:
	build/block getbalance -address $(address)