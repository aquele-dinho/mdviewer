# mdviewer.nvim

In-editor Markdown preview for Neovim using the mdviewer CLI. Renders inside a split or floating terminal window — no browser required.

## Requirements
- Neovim 0.8+
- mdviewer binary in your PATH (or set `mdviewer_path` in setup)

---

## Quick start (most users)
1) Ensure `mdviewer` runs in your shell: `mdviewer --version`.
2) Download or clone this plugin directory (`mdviewer.nvim`).
3) Install it using one of the methods below (plugin manager or manual copy).
4) Add a minimal config (shown below).
5) Open a Markdown file and run `:MDViewerToggle`.

---

## Install methods

### A) lazy.nvim (recommended)
If you have the mdviewer repo locally, point lazy.nvim to this subdirectory:

```lua
{
  dir = '/path/to/mdviewer/contrib/neovim/mdviewer.nvim',
  ft = 'markdown',
  opts = {
    style = 'clean',
    window_type = 'vsplit', -- 'vsplit' | 'hsplit' | 'float'
    auto_refresh = true,
    -- mdviewer_path = '/usr/local/bin/mdviewer',
    -- mermaid_mode = 'terminal', -- 'terminal' | 'svg' | 'png' | 'url'
    -- width = 0, -- mdviewer auto-detects terminal width
    -- debounce_ms = 500,
  },
}
```

**Where to add this configuration:**
Lazy.nvim configuration is typically in one of these locations:
- `~/.config/nvim/lua/plugins.lua` (single plugin file)
- `~/.config/nvim/lua/config/lazy.lua` (dedicated lazy config)
- `~/.config/nvim/lua/plugins/` (directory with multiple plugin files)
- `~/.config/nvim/init.lua` (if configured inline)

To find your lazy.nvim configuration:
```sh
find ~/.config/nvim -name "*.lua" -type f -exec grep -l "lazy" {} \;
```

To keep the plugin separate from the mdviewer repo, clone it into any folder and set `dir` to that path.

### B) Manual install via native packages (no plugin manager)
Neovim can load plugins from the "packages" directory. Copy or clone `mdviewer.nvim` into a `pack/.../start/` directory. Use the OS-specific path below:

- macOS & Linux (XDG): `~/.config/nvim/pack/mdviewer/start/mdviewer.nvim`
- Legacy Vim path (works in Neovim too): `~/.vim/pack/mdviewer/start/mdviewer.nvim`
- Windows: `%USERPROFILE%\AppData\Local\nvim\pack\mdviewer\start\mdviewer.nvim`

Example (macOS/Linux):
```sh
mkdir -p ~/.config/nvim/pack/mdviewer/start
cp -R /path/to/mdviewer/contrib/neovim/mdviewer.nvim ~/.config/nvim/pack/mdviewer/start/
# or clone directly
# git clone https://example.com/mdviewer.nvim ~/.config/nvim/pack/mdviewer/start/mdviewer.nvim
```

Example (Windows PowerShell):
```powershell
New-Item -ItemType Directory -Force "$env:USERPROFILE\AppData\Local\nvim\pack\mdviewer\start" | Out-Null
Copy-Item -Recurse "C:\\path\\to\\mdviewer\\contrib\\neovim\\mdviewer.nvim" "$env:USERPROFILE\AppData\Local\nvim\pack\mdviewer\start\mdviewer.nvim"
```

Restart Neovim and run `:echo &rtp` to confirm the plugin directory is on the runtimepath. See `:h packages` for details.

### C) Other plugin managers
- packer.nvim: `use { dir = '/path/to/mdviewer/contrib/neovim/mdviewer.nvim', ft = 'markdown', config = function() require('mdviewer').setup({}) end }`
- vim-plug: add the directory to `&rtp` or place under `pack/.../start/` (vim-plug is not required when using native packages).

---

## Configure
Create or edit your Neovim config to call `setup`. For init.lua users, add:

```lua
require('mdviewer').setup({
  style = 'clean',
  window_type = 'vsplit',
  auto_refresh = true,
  -- mdviewer_path = '/usr/local/bin/mdviewer',
  -- mermaid_mode = 'terminal', -- 'terminal' | 'svg' | 'png' | 'url'
  -- width = 0,               -- 0 = auto
  -- debounce_ms = 500,
})
```

Available options (defaults shown):
```lua
{
  mdviewer_path = nil,        -- autodetect via exepath('mdviewer')
  style = 'clean',
  mermaid_mode = 'terminal',  -- 'terminal' | 'svg' | 'png' | 'url'
  width = 0,                  -- 0 = auto
  window_type = 'vsplit',     -- 'vsplit' | 'hsplit' | 'float'
  auto_refresh = true,
  debounce_ms = 500,
}
```

---

## Usage
- `:MDViewerToggle`  Toggle preview on/off
- `:MDViewerPreview` Open preview window
- `:MDViewerClose`   Close preview
- `:MDViewerRefresh` Force re-render

Keymap example:
```lua
vim.keymap.set('n', '<leader>mp', '<cmd>MDViewerToggle<cr>', { desc = 'mdviewer: toggle preview' })
```

Typical workflow:
1) Open a Markdown file (`:edit README.md`).
2) Run `:MDViewerToggle` to open a preview split (or set `window_type = 'float'`).
3) Edit your file — the preview auto-refreshes after a short debounce. Use `:MDViewerRefresh` at any time.
4) Close the preview with `:MDViewerClose` or toggle again.

Notes:
- Relative image paths resolve because the mdviewer job runs with `cwd` set to the file’s directory.
- Inline images/diagrams show in terminals that support inline image protocols; otherwise, mdviewer displays a text/ASCII fallback.

---

## Update / Remove
- Update (packages/manual): replace the `mdviewer.nvim` directory with a newer copy or `git pull` if it’s a git clone, then restart Neovim.
- Remove: delete the `mdviewer.nvim` directory from `pack/.../start/` (or remove the lazy.nvim/picker entry) and restart.

---

## Troubleshooting
- "mdviewer binary not found": ensure `mdviewer` is on your PATH or set `mdviewer_path` in setup; verify with `:echo exepath('mdviewer')`.
- Preview doesn’t update: try increasing `debounce_ms` or set `auto_refresh = false` and use `:MDViewerRefresh`.
- No inline images: your terminal may not support inline image protocols; use a compatible terminal or rely on the fallback.
- Verify plugin loaded: run `:scriptnames` and search for `mdviewer.nvim`, or check `:echo &runtimepath`.
