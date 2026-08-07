package logo

// Source is the raw Nix snowflake ASCII art (fastfetch/neofetch style).
// Two glyph classes are used for coloring: '█'-family cells are "light",
// and the rest ('▄','▖', etc. used only as connective glyphs) share the
// same light color in this simple two-tone scheme, split by column half.
var Source = []string{
	`          ▗▄▄▄       ▗▄▄▄▄    ▄▄▄▖`,
	`          ▜███▙       ▜███▙  ▟███▛`,
	`           ▜███▙       ▜███▙▟███▛`,
	`            ▜███▙       ▜██████▛`,
	`     ▟█████████████████▙ ▜████▛     ▟▙`,
	`    ▟███████████████████▙ ▜███▙    ▟██▙`,
	`           ▄▄▄▄▖           ▜███▙  ▟███▛`,
	`          ▟███▛             ▜██▛ ▟███▛`,
	`         ▟███▛               ▜▛ ▟███▛`,
	`▟███████████▛                  ▟██████████▙`,
	`▜██████████▛                  ▟███████████▛`,
	`      ▟███▛ ▟▙               ▟███▛`,
	`     ▟███▛ ▟██▙             ▟███▛`,
	`    ▟███▛  ▜███▙           ▝▀▀▀▀`,
	`    ▜██▛    ▜███▙ ▜██████████████████▛`,
	`     ▜▛     ▟████▙ ▜████████████████▛`,
	`           ▟██████▙         ▜███▙`,
	`          ▟███▛▜███▙         ▜███▙`,
	`         ▟███▛  ▜███▙         ▜███▙`,
	`         ▝▀▀▀    ▀▀▀▀▘         ▀▀▀▘`,
}

// Default colors for the two-tone snowflake (left/top half vs right/bottom
// half); overridable via Parse's colors argument.
const (
	ColorA = "#7ebae4" // lighter blue
	ColorB = "#4a90d9" // darker blue
)
