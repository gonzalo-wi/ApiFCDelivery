package service

import (
	_ "embed"
	"encoding/base64"
	"fmt"
)

//go:embed fonts/montserrat-latin.woff2
var montserratLatinWoff2 []byte

// montserratFontCSS devuelve el bloque @font-face con Montserrat embebida como
// data URI — funciona en todos los clientes sin depender de URLs externas.
var montserratFontCSS = func() string {
	b64 := base64.StdEncoding.EncodeToString(montserratLatinWoff2)
	src := fmt.Sprintf("url(data:font/woff2;base64,%s) format('woff2')", b64)
	return fmt.Sprintf(`<style>
@font-face{font-family:'Montserrat';font-style:normal;font-weight:400;font-display:swap;src:%s;}
@font-face{font-family:'Montserrat';font-style:normal;font-weight:600;font-display:swap;src:%s;}
@font-face{font-family:'Montserrat';font-style:normal;font-weight:700;font-display:swap;src:%s;}
</style>`, src, src, src)
}()
