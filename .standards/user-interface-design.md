# User-interface design standard

## Purpose and scope

This standard defines how projectious.work products look and behave for the
people who use them. It applies to graphical user interfaces, terminal user
interfaces, human-readable CLI output, documentation sites, and websites.
Brand identity (logo, colors, typefaces, voice, imagery) comes from the
[brand repository](https://github.com/projectious-work/brand), which is
authoritative: its released brand contract (documentation, structured data,
generated tokens, and brand MCP server) wins over any other source. Standards
and skills reference it and never restate brand values. This standard governs
only how the brand is applied in interfaces.

Terms in *italics* are [Semantic Anchors](https://github.com/llm-coding/Semantic-Anchors/tree/6554641a11fdf99fdbdf09cc0602268e3db88966); see
[README](README.md#conventions).

## Principles

- Design reviews use *Dieter Rams' 10 Principles of Good Design* as
  criteria. "As little design as possible" decides between otherwise equal
  options.
- Interfaces show the primary task and its state first; advanced options
  unfold on demand (*Progressive Disclosure*). Destructive consequences,
  security state, and errors are never hidden behind disclosure.
- Interface text follows the
  [documentation quality rules](open-source-documentation-strategy.md#quality-accessibility-and-maintenance).

## Interaction and components

- Native Apple-platform interfaces follow the *Apple Human Interface
  Guidelines*.
- Other graphical and web interfaces use *Adobe Spectrum* as the reference for
  component behavior, states, and accessibility where the brand system is
  silent. The reference covers behavior, not Spectrum's visual identity or
  libraries.
- Terminal interfaces follow their platform's keyboard conventions, work at
  80 columns, and never convey meaning by color alone.

## Typography

Typography follows *The Elements of Typographic Style (Bringhurst)*: body text
measure of 45 to 75 characters, a consistent type scale, steady vertical
rhythm, and correct typographic characters (quotes, dashes, ellipses) in
prose.

## Accessibility

Every graphical and web interface provides keyboard operation, visible focus,
sufficient contrast, text alternatives, and reduced-motion behavior, as the
[documentation-site shell](open-source-documentation-strategy.md#documentation-site-shell)
requires for documentation sites. Terminal output provides keyboard
operation, meaning independent of color, and `NO_COLOR` support.

## Agent briefing

> The brand repository's released contract wins; never restate it. Review with
> *Dieter Rams' 10 Principles of Good Design*; primary task first (*Progressive
> Disclosure*), never hide risk or errors. Apple: *Apple Human Interface
> Guidelines*; other GUIs: *Adobe Spectrum* behavior. Long-form text per
> *Bringhurst*. GUIs: keyboard, focus, contrast, text alternatives, reduced
> motion; terminals: 80 columns, never color alone, `NO_COLOR`.
