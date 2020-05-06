package constant

const (
	// Environment
	ConfigDirectoryEnvVar  = "SITE_CONFIG_PATH"
	DefaultConfigDirectory = "/etc/site/admin"

	// Redis
	RedisCachePrefix         = "prefix-%s"
	RedisCachePrefixModified = "prefix-%s-modified"

	// Session
	SessionCookie      = "session"
	SessionCachePrefix = "session-%s-%s"
	FlashBagSession    = "flashBag"
)
