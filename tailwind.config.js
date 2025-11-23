/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./components/**/*.templ"
  ],
  theme: {
    extend: {
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideInFromRight: {
          '0%': { transform: 'translateX(100%)' },
          '100%': { transform: 'translateX(0)' },
        },
        slideOutToRight: {
          '0%': { transform: 'translateX(0)' },
          '100%': { transform: 'translateX(100%)' },
        }
      },
      animation: {
        fadeIn: 'fadeIn 0.3s ease-in-out',
        slideInFromRight: 'slideInFromRight 0.2s ease-out forwards',
        slideOutToRight: 'slideOutToRight 0.2s ease-in forwards',
      }
    },
  },
  plugins: [],
}

