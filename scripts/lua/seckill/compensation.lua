-- KEY[1] seckill:stock:
-- KEY[2] seckill:purchasedusers:
-- ARGV[1] userid
redis.call('srem', KEY[2], ARGV[1])
redis.call('incr', KEY[1])
return 1
