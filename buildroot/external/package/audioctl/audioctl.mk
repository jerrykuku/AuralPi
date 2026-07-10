################################################################################
#
# audioctl
#
################################################################################

AUDIOCTL_VERSION = 1.0.0
AUDIOCTL_SITE = $(BR2_EXTERNAL_RPI4_USBDAC_PATH)/package/audioctl/src
AUDIOCTL_SITE_METHOD = local
AUDIOCTL_LICENSE = MIT

define AUDIOCTL_BUILD_CMDS
	$(HOST_GO_TARGET_ENV) GOPROXY=off $(HOST_DIR)/bin/go build -trimpath -ldflags "-s -w" -o $(@D)/audioctl .
endef

define AUDIOCTL_INSTALL_TARGET_CMDS
	$(INSTALL) -D -m 0755 $(@D)/audioctl $(TARGET_DIR)/usr/sbin/audioctl
endef

$(eval $(generic-package))
