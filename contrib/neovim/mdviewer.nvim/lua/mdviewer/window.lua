local M = {}

local api = vim.api

-- Create a terminal buffer in a split or floating window
function M.open_terminal(opts)
  opts = opts or {}
  local kind = opts.window_type or 'vsplit'

  -- Create buffer
  local buf = api.nvim_create_buf(false, true)
  api.nvim_buf_set_option(buf, 'bufhidden', 'wipe')
  api.nvim_buf_set_option(buf, 'filetype', 'mdviewer')

  local win
  if kind == 'float' then
    -- Centered floating window (80% size)
    local columns = api.nvim_get_option_value('columns', {})
    local lines = api.nvim_get_option_value('lines', {})
    local width = math.floor(columns * 0.8)
    local height = math.floor(lines * 0.8)
    local row = math.floor((lines - height) / 2)
    local col = math.floor((columns - width) / 2)

    win = api.nvim_open_win(buf, true, {
      relative = 'editor',
      row = row,
      col = col,
      width = width,
      height = height,
      border = 'rounded',
    })
  elseif kind == 'hsplit' then
    vim.cmd('split')
    win = api.nvim_get_current_win()
    api.nvim_win_set_buf(win, buf)
  else
    vim.cmd('vsplit')
    win = api.nvim_get_current_win()
    api.nvim_win_set_buf(win, buf)
  end

  -- Terminal-style options
  api.nvim_set_option_value('number', false, { win = win })
  api.nvim_set_option_value('relativenumber', false, { win = win })
  api.nvim_set_option_value('signcolumn', 'no', { win = win })
  api.nvim_buf_set_option(buf, 'scrollback', 100000)

  return buf, win
end

function M.close(buf, win)
  if win and api.nvim_win_is_valid(win) then
    pcall(api.nvim_win_close, win, true)
  end
  if buf and api.nvim_buf_is_valid(buf) then
    pcall(api.nvim_buf_delete, buf, { force = true })
  end
end

return M
