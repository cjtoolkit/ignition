package constant

const (
	// Environment
	ConfigDirectoryEnvVar  = "SITE_CONFIG_PATH"
	DefaultConfigDirectory = "/etc/site/admin"

	// Cache
	CachePrefix         = "prefix-%s"
	CachePrefixModified = "prefix-%s-modified"
	CacheFileFolderName = "site"

	// Session
	SessionCookie      = "session"
	SessionCachePrefix = "session-%s-%s"
	FlashBagSession    = "flashBag"
)
