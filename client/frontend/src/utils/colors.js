export function hex2rgba(hex) {
  if (hex.charAt(0) === '#') hex = hex.substring(1)
  const parsed = hex.length === 6 || hex.length === 8
    ? /^([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})?$/i.exec(hex)
    : /^([a-f\d])([a-f\d])([a-f\d])([a-f\d])?$/i.exec(hex)
  if (hex.length < 6) {
    parsed[1] += parsed[1]
    parsed[2] += parsed[2]
    parsed[3] += parsed[3]
    if (parsed[4]) {
      parsed[4] += parsed[4]
    }
  }
  if (!parsed[4]) {
    parsed[4] = "ff"
  }
  return parsed ? [
    parseInt(parsed[1], 16),
    parseInt(parsed[2], 16),
    parseInt(parsed[3], 16),
    parseInt(parsed[4], 16) / 255
  ] : null
}

export function rgba2hex(rgba) {
  const r = (rgba[0] | 1 << 8).toString(16).slice(1)
  const g = (rgba[1] | 1 << 8).toString(16).slice(1)
  const b = (rgba[2] | 1 << 8).toString(16).slice(1)
  const rgb = r + g + b

  if (!rgba[3]) {
    return '#' + rgb
  } else {
    const a = ((rgba[3] * 255) | 1 << 8).toString(16).slice(1)
    return '#' + rgb + a
  }
}

export function rgb2hsl(rgb) {
  rgb[0] /= 255
  rgb[1] /= 255
  rgb[2] /= 255

  const max = Math.max(rgb[0], rgb[1], rgb[2])
  const min = Math.min(rgb[0], rgb[1], rgb[2])
  let h = (max + min) / 2
  let s = h
  let l = h

  if (max == min) {
    h = s = 0
  } else {
    const d = max - min
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min)

    if (max === rgb[0]) {
      h = (rgb[1] - rgb[2]) / d + (rgb[1] < rgb[2] ? 6 : 0)
    } else if (max === rgb[1]) {
      h = (rgb[2] - rgb[0]) / d + 2
    } else if (max === rgb[2]) {
      h = (rgb[0] - rgb[1]) / d + 4
    }

    h /= 6
  }

  return [h, s, l]
}

export function hsl2rgb(hsl) {
  if (hsl[1] == 0) {
    return [hsl[2] * 255, hsl[2] * 255, hsl[2] * 255]
  } else {
    function hue2rgb(p, q, t) {
      if (t < 0) t += 1;
      if (t > 1) t -= 1;
      if (t < 1/6) return p + (q - p) * 6 * t;
      if (t < 1/2) return q;
      if (t < 2/3) return p + (q - p) * (2/3 - t) * 6;
      return p;
    }

    const q = hsl[2] < 0.5 ? hsl[2] * (1 + hsl[1]) : hsl[2] + hsl[1] - hsl[2] * hsl[1];
    const p = 2 * hsl[2] - q;

    const r = hue2rgb(p, q, hsl[0] + 1/3)
    const g = hue2rgb(p, q, hsl[0])
    const b = hue2rgb(p, q, hsl[0] - 1/3)

    return [r * 255, g * 255, b * 255]
  }
}
