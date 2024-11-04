package interfaces

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/theme"
    "image/color"
)

type CustomTheme struct{}

func (CustomTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
    return theme.DefaultTheme().Color(n, v)
}

func (CustomTheme) Font(s fyne.TextStyle) fyne.Resource {
    return theme.DefaultTheme().Font(s)
}

func (CustomTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
    return theme.DefaultTheme().Icon(n)
}

func (CustomTheme) Size(n fyne.ThemeSizeName) float32 {
    switch n {
    case theme.SizeNameText:
        return theme.DefaultTheme().Size(n) * 1.5 // Increase text size by 50%
    default:
        return theme.DefaultTheme().Size(n)
    }
}
