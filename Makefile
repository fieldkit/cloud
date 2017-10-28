BUILD=$(shell pwd)/build
dirs=src sqs-tester sqs-worker

all:
	mkdir -p $(BUILD)
	for dir in $(dirs); do       \
		(cd $$dir && BUILD=$(BUILD) make);  \
	done

clean:
	rm -rf $(BUILD)
