-- 检查是否已存在
if redis.call('exists', KEYS[1]) == 0 then
    -- 创建活动数据
    redis.call('hmset', KEYS[1], 'id', ARGV[1], 'status', ARGV[2], 'start_time', ARGV[3], 'end_time', ARGV[4])
    -- 设置过期时间
    redis.call('expireat', KEYS[1], ARGV[5])
    return 1
else
    return 0
end
