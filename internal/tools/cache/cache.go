package tools

// Remember 从缓存中获取数据，如果不存在则从数据库中获取
func Remember[T any](key string, ttl int, callback func() (*T, error)) (*T, error) {
    // TODO 实现缓存逻辑
    return callback()
}
