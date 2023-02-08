import { APCAcontrast, sRGBtoY } from 'apca-w3'
import { hex2rgba, rgba2hex, rgb2hsl, hsl2rgb } from './colors'

// taken from apca-w3 sources and slightly modified to make it work in this context
function alphaBlend (rgbaFG = [0, 0, 0, 1.0], rgbBG = [0, 0, 0], isInt = true) {
  rgbaFG[3] = Math.max(Math.min(rgbaFG[3], 1.0), 0.0)
  const compBlend = 1.0 - rgbaFG[3]
  const rgbOut = [0, 0, 0]

  for (let i = 0; i < 3; i++) {
    rgbOut[i] = rgbBG[i] * compBlend + rgbaFG[i] * rgbaFG[3]
    if (isInt) rgbOut[i] = Math.min(Math.round(rgbOut[i]), 255)
  }
  return rgbOut
}

export function contrast(color1, color2) {
  if (typeof color1 === 'string') color1 = hex2rgba(color1)
  if (typeof color2 === 'string') color2 = hex2rgba(color2)
  const y1 = sRGBtoY(alphaBlend(color1, color2))
  const y2 = sRGBtoY(color2)
  const c = APCAcontrast(y1, y2)
  return c < 0 ? c * -1 : c
}

function deriveContrast(base, config) {
  if (!Array.isArray(config.options)) {
    console.error('contrast derivation options missing or invalid')
    return ''
  }

  const selected = config.options.map(opt => {
    // APCA contrast calculation is polar, therefore conditional parameter order
    if (config.foreground) {
      return [opt, contrast(opt, base)]
    } else {
      return [opt, contrast(base, opt)]
    }
  }).reduce((a, b) => {
    if (a[1] >= b[1]) {
      return a
    }
    return b
  }, [undefined, 0])

  return `${config.var}: ${selected[0]};`
}

function deriveOpacity(base, config) {
  const opacity = Number(config.opacity)
  const rgba = hex2rgba(base)
  rgba[3] = opacity
  return `${config.var}: ${rgba2hex(rgba)};`
}

function deriveHueShift(base, config) {
  const hs = Number(config.hueShift) || 0
  const hsl = rgb2hsl(hex2rgba(base))
  hsl[0] = (((hsl[0] * 360) + hs) % 360) / 360
  return `${config.var}: ${rgba2hex(hsl2rgb(hsl))};`
}

export function deriveColor(base, config) {
  if (!config.var) console.error('color derivation missing target variable')

  switch (config.type) {
    case 'contrast':
      return deriveContrast(base, config)
    case 'opacity':
      return deriveOpacity(base, config)
    case 'hueShift':
      return deriveHueShift(base, config)
    default:
      console.error(`unknown color derivation: ${config.type}`)
  }
}
