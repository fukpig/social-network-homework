function getAlphaChar()
    return string.char(math.random(65, 90))
end

request = function()
  local path = wrk.path .. "?name=" .. getAlphaChar() .. "&surname=" .. getAlphaChar()
  return wrk.format("GET", path, wrk.headers, wrk.body)
end
