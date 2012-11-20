MAKE=make
TARGS=all install
ALL_DIRS=src
TEST_DIRS=test

$(TARGS): 
	for b in $(ALL_DIRS) ; do $(MAKE) -C $$b $@ ; done

.PHONY: test
test: 
	for b in $(TEST_DIRS) ; do $(MAKE) -C $$b $@ ; done

