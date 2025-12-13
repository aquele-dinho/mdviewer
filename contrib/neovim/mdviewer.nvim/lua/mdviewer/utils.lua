local M = {}

local api = vim.api

function M.is_markdown(bufnr)
  bufnr = bufnr or api.nvim_get_current_buf()
  local ft = api.nvim_buf_get_option(bufnr, 'filetype')
  return ft == 'markdown' or ft == 'md'
end

function M.buf_dir(bufnr)
  bufnr = bufnr or api.nvim_get_current_buf()
  local name = api.nvim_buf_get_name(bufnr)
  if name == '' then
    return vim.loop.cwd()
  end
  return vim.fn.fnamemodify(name, ':p:h')
end

-- Simple debounce wrapper using luv timers
function M.debounce(ms, fn)
  local timer = vim.loop.new_timer()
  local running = false
  return function(...)
    local argv = { ... }
    if running then
      timer:stop()
    end
    running = true
    timer:start(ms, 0, function()
      timer:stop()
      running = false
      vim.schedule(function()
        fn(unpack(argv))
      end)
    end)
  end
end

return M
