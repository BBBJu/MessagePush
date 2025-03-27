-- KEYS[1]: 有序集合的key
-- ARGV[1]: 当前时间戳
local elements = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', ARGV[1])
if #elements > 0 then
    redis.call('ZREM', KEYS[1], table.unpack(elements))
end
return elements