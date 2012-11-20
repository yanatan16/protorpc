MAKE=make
TARGS=all install
DIRS=proto go

$(TARGS):
	for b in $(DIRS) ; do $(MAKE) -C $$b $@ ; done
