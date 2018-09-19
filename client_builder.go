package zabbix

// ClientBuilder is Zabbix API client builder
type ClientBuilder struct {
	cache         SessionAbstractCache
	hasCache      bool
	url           string
	credentials   map[string]string
	skipTLSVerify bool
}

// WithCache sets cache for Zabbix sessions
func (builder *ClientBuilder) WithCache(cache SessionAbstractCache) *ClientBuilder {
	builder.cache = cache
	builder.hasCache = true

	return builder
}

// WithCredentials sets auth credentials for Zabbix API
func (builder *ClientBuilder) WithCredentials(username string, password string) *ClientBuilder {
	builder.credentials["username"] = username
	builder.credentials["password"] = password

	return builder
}

// SkipTLSVerify sets TLS verification for Zabbix API URL
func (builder *ClientBuilder) SkipTLSVerify() *ClientBuilder {
	builder.skipTLSVerify = true

	return builder
}

// Connect creates Zabbix API client and connects to the API server
// or provides a cached server if any cache was specified
func (builder *ClientBuilder) Connect() (session *Session, err error) {
	// Check if any cache was defined and if it has a valid cached session
	if builder.hasCache && builder.cache.HasSession() {
		if session, err = builder.cache.GetSession(); err == nil {
			return session, nil
		}
	}

	// Otherwise - login to a Zabbix server
	session, err = NewSession(builder.url, builder.credentials["username"], builder.credentials["password"], builder.skipTLSVerify)

	if err != nil {
		return nil, err
	}

	// Try to cache session if any cache used
	if builder.hasCache {
		return session, builder.cache.SaveSession(session)
	}

	return session, err
}

// CreateClient creates a Zabbix API client builder
func CreateClient(apiEndpoint string) *ClientBuilder {
	return &ClientBuilder{
		url:           apiEndpoint,
		credentials:   make(map[string]string),
		skipTLSVerify: false,
	}
}
