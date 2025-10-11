local id = redis.call('HGET', KEYS[1], 'id')
local status = redis.call('HGET', KEYS[1], 'status')
local start_time = redis.call('HGET', KEYS[1], 'start_time')
local end_time = redis.call('HGET', KEYS[1], 'end_time')

if not id or not status or not start_time or not end_time then
    return 0
end

local current_time = tonumber(ARGV[1])
local start_time_num = tonumber(start_time)
local end_time_num = tonumber(end_time)
local status_num = tonumber(status)

if not start_time_num or not end_time_num or not status_num then
    return 0
end

-- 如果活动已结束，更新状态为3
if current_time > end_time_num and status_num ~= 3 then
    redis.call('HSET', KEYS[1], 'status', '3')
    return 3
    -- 如果活动正在进行，更新状态为2
elseif current_time >= start_time_num and current_time <= end_time_num and status_num ~= 2 then
    redis.call('HSET', KEYS[1], 'status', '2')
    return 2
    -- 如果活动未开始，确保状态为1
elseif current_time < start_time_num and status_num ~= 1 then
    redis.call('HSET', KEYS[1], 'status', '1')
    return 1
end
return 4

