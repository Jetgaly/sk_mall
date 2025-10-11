-- KEYS[1] seckill:stock:
-- KEYS[2] seckill:purchasedusers:
-- ARGV[1] userid
local stock = redis.call('GET', KEYS[1])
if stock == false then
    return 0
    -- 不存在key
end

-- 将字符串转换为数字进行比较
local stock_num = tonumber(stock)

-- 库存不足
if stock_num <= 0 then
    return 1
end

-- 用户集
local ok = redis.call('sadd', KEYS[2], ARGV[1])
if ok == 0 then
    -- 用户已存在
    return 2
end

-- stock - 1
redis.call('decr', KEYS[1])
return 3
-- success