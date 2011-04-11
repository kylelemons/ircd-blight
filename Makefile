# Explicitly make targets phony, just in case
.PHONY : all deppkg pkgs cmds install clean nuke test bench

# By default, build everything
all : pkgs ircd

deppkg :
	@echo "Installing dependencies..."
	#goinstall -u goconf.googlecode.com/hg
	#goinstall -u log4go.googlecode.com/hg

pkgs :
	@echo "Building packages..."
	@$(MAKE) -C src/pkg/ircd -f Makefile.sub install
	@$(MAKE) -C src/pkg/ircd install

ircd : pkgs
	@echo "Building ircd..."
	@$(MAKE) -C src/cmd/ircd

cmds :
	@echo "Building helpers..."

test bench clean nuke install :
	@echo "Performing $@..."
	@$(MAKE) -C src/pkg/ircd -f Makefile.sub $@
	@$(MAKE) -C src/pkg/ircd $@
	@$(MAKE) -C src/cmd/ircd $@
