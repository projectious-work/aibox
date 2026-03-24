" projectious.work color scheme for Vim
" Midnight navy + accent orange palette

set background=dark
hi clear
if exists("syntax_on")
  syntax reset
endif
let g:colors_name = "projectious"

" UI elements
hi Normal       guifg=#C5DAF0 guibg=#0E1720
hi CursorLine                 guibg=#131E2B
hi CursorLineNr guifg=#E05232               gui=bold
hi LineNr       guifg=#546A82
hi Visual                     guibg=#1d3352
hi StatusLine   guifg=#C5DAF0 guibg=#1d3352 gui=bold
hi StatusLineNC guifg=#546A82 guibg=#131E2B
hi VertSplit    guifg=#1d3352 guibg=#0E1720
hi TabLine      guifg=#546A82 guibg=#131E2B
hi TabLineSel   guifg=#C5DAF0 guibg=#1d3352 gui=bold
hi TabLineFill  guifg=#131E2B guibg=#131E2B
hi Pmenu        guifg=#C5DAF0 guibg=#131E2B
hi PmenuSel     guifg=#0E1720 guibg=#8AACC8
hi Search       guifg=#0E1720 guibg=#8B6508
hi IncSearch    guifg=#0E1720 guibg=#E05232
hi MatchParen   guifg=#E05232               gui=bold
hi Folded       guifg=#546A82 guibg=#131E2B
hi FoldColumn   guifg=#546A82 guibg=#0E1720
hi SignColumn                 guibg=#0E1720
hi NonText      guifg=#1d3352
hi SpecialKey   guifg=#1d3352
hi Directory    guifg=#8AACC8

" Syntax groups
hi Comment      guifg=#546A82               gui=italic
hi Constant     guifg=#E05232
hi String       guifg=#2D6A4F
hi Character    guifg=#2D6A4F
hi Number       guifg=#EA7558
hi Float        guifg=#EA7558
hi Boolean      guifg=#E05232
hi Identifier   guifg=#C5DAF0
hi Function     guifg=#8AACC8               gui=bold
hi Statement    guifg=#2B4D78               gui=bold
hi Conditional  guifg=#2B4D78               gui=bold
hi Repeat       guifg=#2B4D78               gui=bold
hi Operator     guifg=#C5DAF0
hi Keyword      guifg=#2B4D78               gui=bold
hi PreProc      guifg=#8B6508
hi Include      guifg=#8B6508
hi Define       guifg=#8B6508
hi Type         guifg=#8AACC8
hi StorageClass guifg=#8AACC8
hi Structure    guifg=#8AACC8
hi Special      guifg=#EA7558
hi Delimiter    guifg=#546A82
hi Underlined   guifg=#8AACC8               gui=underline
hi Error        guifg=#C5DAF0 guibg=#A32D2D
hi Todo         guifg=#8B6508 guibg=#0E1720 gui=bold

" Diff
hi DiffAdd      guifg=#2D6A4F guibg=#131E2B
hi DiffChange   guifg=#8B6508 guibg=#131E2B
hi DiffDelete   guifg=#A32D2D guibg=#131E2B
hi DiffText     guifg=#C5DAF0 guibg=#1d3352 gui=bold

" Terminal colors (for :terminal in vim 8+)
let g:terminal_ansi_colors = [
  \ '#131E2B', '#A32D2D', '#2D6A4F', '#8B6508',
  \ '#2B4D78', '#546A82', '#8AACC8', '#C5DAF0',
  \ '#1d3352', '#E05232', '#3D5C4A', '#EA7558',
  \ '#8AACC8', '#7E8F9E', '#C5DAF0', '#E2E9F2'
  \ ]
