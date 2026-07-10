################################################################################
#
# roonbridge
#
################################################################################

ROONBRIDGE_VERSION = 1.0
ROONBRIDGE_SITE = $(call qstrip,$(BR2_PACKAGE_ROONBRIDGE_SITE))
ROONBRIDGE_SOURCE = RoonBridge_linuxarmv8.tar.bz2
ROONBRIDGE_LICENSE = Roon proprietary
ROONBRIDGE_REDISTRIBUTE = NO

define ROONBRIDGE_INSTALL_TARGET_CMDS
	$(INSTALL) -d -m 0755 $(TARGET_DIR)/opt/RoonBridge
	cp -dpfr $(@D)/* $(TARGET_DIR)/opt/RoonBridge/
endef

$(eval $(generic-package))
