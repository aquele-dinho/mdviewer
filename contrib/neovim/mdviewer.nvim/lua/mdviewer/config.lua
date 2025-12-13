local M = {}

local defaults = {
  mdviewer_path = nil, -- auto-detect via exepath('mdviewer') if nil
  style = 'clean',
  mermaid_mode = 'terminal', -- 'terminal' | 'svg' | 'png' | 'url'
  width = 0, -- 0 = auto (mdviewer will detect)
  window_type = 'vsplit', -- 'vsplit' | 'hsplit' | 'float'
  auto_refresh = true,
  debounce_ms = 500,
}

M.options = vim.deepcopy(defaults)

function M.setup(opts)
  opts = opts or {}
  M.options = vim.tbl_deep_extend('force', defaults, opts)
end

function M.get()
  -- Resolve mdviewer binary if not provided
  if not M.options.mdviewer_path or M.options.mdviewer_path == '' then
    local exe = vim.fn.exepath('mdviewer')
    if exe ~= '' then
      M.options.mdviewer_path = exe
    end
  end
  return M.options
end

return M
