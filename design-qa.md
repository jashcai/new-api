**Source Visual Truth**
- Source reference: `/var/folders/3m/0sr6zxjx6wxbw8tfnb4kfvpc0000gn/T/codex-clipboard-1d96a2d7-cf86-4a82-ad52-3ed24fe6242a.png`
- Implementation screenshot: `/tmp/openbridge-enterprise-cards.png`
- Full-view comparison evidence: `/tmp/openbridge-cards-comparison.png`
- Viewport: desktop, current in-app browser viewport, homepage scrolled/cropped to enterprise card section
- State: public homepage, light theme

**Focused Region Comparison**
- Focused region used: enterprise operation card grid.
- Reason: the requested change is scoped to the four-card homepage section shown in the supplied screenshots.

**Findings**
- No actionable P0/P1/P2 issues found.
- Typography: the implementation keeps the product's existing font stack and uses stronger card titles, readable Chinese body copy, and stable line heights.
- Spacing and layout rhythm: the section now follows the reference's top-visual/bottom-content card rhythm, with a divider and consistent four-column desktop grid.
- Colors and visual tokens: the cards preserve OpenBridge's teal/blue theme while using a calmer white surface, subtle borders, and light shadows similar to the reference.
- Image quality and asset fidelity: the section uses icon-library visuals and compact interface-style illustrations, matching the reference pattern without introducing raster placeholders.
- Copy and content: all descriptions are localized Chinese and aligned with private enterprise AI gateway use cases.

**Patches Made**
- Replaced the original simple enterprise operation cards with reference-style capability cards.
- Added four visual treatments: model entry, availability routing, performance observation, and policy/audit.
- Rewrote card descriptions in Chinese.
- Added responsive card styling and visual-section styling in the classic theme stylesheet.

**Implementation Checklist**
- Build frontend bundle.
- Restart embedded Go service.
- Verify homepage renders four capability cards with visual areas and Chinese descriptions.

**Follow-up Polish**
- P3: if desired, the model-icon cloud can be swapped for real provider brand icons after final legal/brand review.

final result: passed
