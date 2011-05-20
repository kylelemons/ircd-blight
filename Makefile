# Explicitly make targets phony, just in case
.PHONY : all pkgs cmds install clean nuke test bench

MAKE += -s

PKGS = parser user conn core
CMDS = ircd

# By default, build everything
all : pkgs cmds
	@

define recurse
@echo "$(3) $(1) $(2)" | sed -e 's/^all/build/'
@$(MAKE) -C src/$(1)/$(2) $(3)

endef

define all_packages
$(foreach pkg,$(PKGS),$(call recurse,pkg,$(pkg),$(1)))
endef

define all_commands
$(foreach cmd,$(CMDS),$(call recurse,cmd,$(cmd),$(1)))
endef

pkgs :
	$(call all_packages,install)

cmds :
	$(call all_commands,all)

install clean nuke :
	$(call all_packages,$@)
	$(call all_commands,$@)

test bench :
	$(call all_packages,$@)
