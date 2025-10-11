local start_time = redis.call('HGET', KEYS[1], 'start_time')
local end_time = redis.call('HGET', KEYS[1], 'end_time')
local status = redis.call('HGET', KEYS[1], 'status')
if not start_time or not end_time or not status then
    return 0
end

local current_time = tonumber(ARGV[1])
local start_time_num = tonumber(start_time)
local end_time_num = tonumber(end_time)
local status_num = tonumber(status)

if not start_time_num or not end_time_num or not status_num then
    return 0
end

if current_time >= start_time_num and current_time <= end_time_num and status_num == 1 then
    return end_time_num-current_time
end
return 0