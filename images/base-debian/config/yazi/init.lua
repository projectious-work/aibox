-- =============================================================================
-- Yazi init.lua — aibox defaults
-- Runs on every Yazi startup. Register plugins that need setup here.
-- =============================================================================

-- git.yazi: show git status (modified/untracked/ignored/staged) in file list.
-- Fetcher registration is in yazi.toml [plugin.prepend_fetchers].
require("git"):setup {}
