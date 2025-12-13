local M = {}

local config = require('mdviewer.config')
local preview = require('mdviewer.preview')

function M.setup(opts)
  config.setup(opts)
  preview.setup_autorefresh()
end

-- Commands
vim.api.nvim_create_user_command('MDViewerPreview', function()
  preview.open()
end, {})

vim.api.nvim_create_user_command('MDViewerClose', function()
  preview.close()
end, {})

vim.api.nvim_create_user_command('MDViewerToggle', function()
  if preview.is_active() then
    preview.close()
  else
    preview.open()
  end
end, {})

vim.api.nvim_create_user_command('MDViewerRefresh', function()
  preview.refresh()
end, {})

return M
