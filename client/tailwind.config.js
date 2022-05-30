const colors = require("tailwindcss/colors");

module.exports = {
  mode: "jit",
  content: ["./src/**/*.{js,jsx,ts,tsx}"],
  theme: {
    extend: {
      fontFamily: {
        Ropa: ["Ropa", "sans-serif"],
        RopaItalic: ["RopaItalic", "sans-serif"],
      },
      listStyleType: {
        none: "none",
        disc: "disc",
        decimal: "decimal",
        square: "square",
        roman: "upper-roman",
      },
    },
    colors: {
      transparent: "transparent",
      current: "currentColor",
      black: colors.black,
      white: colors.white,
      gray: colors.gray,
      red: colors.red,
      yellow: colors.amber,
      sky: colors.sky,
      purple: colors.purple,
      ...colors,
      ross2: "#6a1b9a",
      lwhite: {
        1: "#F4F4F4",
        2: "#E4E4E4",
      },
      indigo: "#425CBA",
    },
  },
  plugins: [],
  corePlugins: {
    fontFamily: true,
  },
};
