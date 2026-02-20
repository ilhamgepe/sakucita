package utils

import (
	"fmt"
	"strings"

	"sakucita/internal/server/security"

	"github.com/gofiber/fiber/v3"
)

func ExtractClientInfo(c fiber.Ctx) security.ClientInfo {
	ip := c.IP()
	userAgent := c.Get("User-Agent")
	deviceName := parseDeviceName(userAgent)
	return security.ClientInfo{
		IP:         ip,
		UserAgent:  userAgent,
		DeviceName: deviceName,
	}
}

func parseDeviceName(ua string) string {
	browser := "Unknown Browser"
	os := "Unknown OS"

	switch {
	case strings.Contains(ua, "Chrome"):
		browser = "Chrome"
	case strings.Contains(ua, "Safari") && !strings.Contains(ua, "Chrome"):
		browser = "Safari"
	case strings.Contains(ua, "Firefox"):
		browser = "Firefox"
	case strings.Contains(ua, "Edge"):
		browser = "Edge"
	}

	switch {
	case strings.Contains(ua, "Windows"):
		os = "Windows"
	case strings.Contains(ua, "Mac OS X"):
		os = "macOS"
	case strings.Contains(ua, "Android"):
		os = "Android"
	case strings.Contains(ua, "iPhone"):
		os = "iOS"
	case strings.Contains(ua, "Linux"):
		os = "Linux"
	}

	return fmt.Sprintf("%s on %s", browser, os)
}
