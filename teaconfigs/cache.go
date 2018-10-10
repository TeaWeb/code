package teaconfigs

// 缓存配置
// 参考：http://nginx.org/en/docs/http/ngx_http_proxy_module.html#proxy_cache_path
type CacheConfig struct {
	Key     string `yaml:"key"`     //@TODO
	Path    string `yaml:"path"`    //@TODO
	MaxSize string `yaml:"maxSize"` //@TODO
	MaxLife string `yaml:"maxLife"` //@TODO
	Memory  bool   `yaml:"memory"`  // @TODO 是否为内存缓存，默认为文件缓存
}

func (this *CacheConfig) Validate() error {
	if len(this.Key) == 0 {
		this.Key = "${host}${requestURI}"
	}
	return nil
}
