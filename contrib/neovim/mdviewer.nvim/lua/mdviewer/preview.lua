local M = {}

local api = vim.api
local config = require('mdviewer.config')
local window = require('mdviewer.window')
local utils = require('mdviewer.utils')

local state = {
  active = false,
  source_buf = nil,
  term_buf = nil,
  term_win = nil,
  term_chan = nil,
  temp_path = nil,
}

-- Build mdviewer CLI args for a given input path (file)
local function mdviewer_cmd(input_path)
  local bin = config.get().mdviewer_path or 'mdviewer'
  local args = { bin, input_path }
  local o = config.get()
  if o.style and o.style ~= '' then
    table.insert(args, '--style')
    table.insert(args, o.style)
  end
  if o.mermaid_mode and o.mermaid_mode ~= '' then
    table.insert(args, '--mermaid-mode')
    table.insert(args, o.mermaid_mode)
  end
  if o.width and o.width > 0 then
    table.insert(args, '--width')
    table.insert(args, tostring(o.width))
  end
  return args
end

-- Write current buffer to a temp markdown file and return its path
-- Prefer placing the temp file in the source buffer's directory so mdviewer
-- resolves relative image paths correctly.
local function write_temp_from_buffer(bufnr)
  local dir = utils.buf_dir(bufnr)
  local fname = string.format('.mdviewer_preview_%s_%d.md', vim.fn.strftime('%Y%m%d%H%M%S'), math.random(1000, 9999))
  local tmp = dir .. '/' .. fname
  local lines = api.nvim_buf_get_lines(bufnr, 0, -1, false)
  -- Ensure newline at EOF
  if #lines == 0 or lines[#lines] ~= '' then
    table.insert(lines, '')
  end
  local ok = pcall(vim.fn.writefile, lines, tmp)
  if not ok then
    -- Fallback to system temp if directory is not writable
    tmp = vim.fn.tempname() .. '.md'
    vim.fn.writefile(lines, tmp)
  end
  return tmp
end

local function start_term(buf, win, cwd, args)
  api.nvim_set_current_win(win)
  api.nvim_win_set_buf(win, buf)
  local chan = vim.fn.termopen(args, {
    cwd = cwd,
    env = {
      TERM = 'xterm-256color',
      COLORTERM = 'truecolor',
    },
    on_exit = function()
      -- Leave window/buffer for user to inspect; close handled by Close command
    end,
  })
  return chan
end

function M.is_active()
  return state.active and state.term_buf and api.nvim_buf_is_valid(state.term_buf)
end

function M.open()
  local o = config.get()
  if not o.mdviewer_path or o.mdviewer_path == '' then
    vim.notify('[mdviewer.nvim] mdviewer binary not found in PATH', vim.log.levels.ERROR)
    return
  end

  local src = api.nvim_get_current_buf()
  if not utils.is_markdown(src) then
    vim.notify('[mdviewer.nvim] Current buffer is not markdown', vim.log.levels.WARN)
  end

  local buf, win = window.open_terminal({ window_type = o.window_type })
  local cwd = utils.buf_dir(src)

  -- Write to temp file and run mdviewer on it (prevents PTY echo of stdin)
  local tmp = write_temp_from_buffer(src)
  local args = mdviewer_cmd(tmp)
  local chan = start_term(buf, win, cwd, args)

  state.active = true
  state.source_buf = src
  state.term_buf = buf
  state.term_win = win
  state.term_chan = chan
  state.temp_path = tmp
end

function M.close()
  if state.term_chan then
    pcall(vim.fn.chanclose, state.term_chan)
  end
  window.close(state.term_buf, state.term_win)
  if state.temp_path then
    pcall(vim.fn.delete, state.temp_path)
  end
  state.active = false
  state.source_buf = nil
  state.term_buf = nil
  state.term_win = nil
  state.term_chan = nil
  state.temp_path = nil
end

function M.refresh()
  if not M.is_active() then return end
  if not state.source_buf or not api.nvim_buf_is_valid(state.source_buf) then return end
  local o = config.get()
  -- Recreate terminal to ensure a clean render
  window.close(state.term_buf, state.term_win)
  local buf, win = window.open_terminal({ window_type = o.window_type })
  local cwd = utils.buf_dir(state.source_buf)

  -- Update temp file content
  if state.temp_path then
    local lines = api.nvim_buf_get_lines(state.source_buf, 0, -1, false)
    if #lines == 0 or lines[#lines] ~= '' then table.insert(lines, '') end
    vim.fn.writefile(lines, state.temp_path)
  else
    state.temp_path = write_temp_from_buffer(state.source_buf)
  end

  local args = mdviewer_cmd(state.temp_path)
  local chan = start_term(buf, win, cwd, args)

  state.term_buf = buf
  state.term_win = win
  state.term_chan = chan
end

function M.setup_autorefresh()
  local o = config.get()
  if not o.auto_refresh then return end

  local debounced = utils.debounce(o.debounce_ms, function()
    if M.is_active() then
      M.refresh()
    end
  end)

  local augroup = api.nvim_create_augroup('MDViewerPreview', { clear = true })
  api.nvim_create_autocmd({ 'TextChanged', 'TextChangedI' }, {
    group = augroup,
    callback = function(ev)
      if state.source_buf and ev.buf == state.source_buf then
        debounced()
      end
    end,
  })
end

return M
