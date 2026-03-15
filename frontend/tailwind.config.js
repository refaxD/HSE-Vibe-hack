/** @type {import('tailwindcss').Config} */

module.exports = {
  content: [
    "./src/**/*.{html,ts}",
  ],
  theme: {
    extend: {
      colors: {
        'primary-color': '#ffffff',
        'secondary-color': '#1d2025',
        'accent-light': '#799496',
        'accent-color-1': '#a30b37',
        'accent-color-2': '#d78521',
        'accent-color-3': '#A259FF',
        'text-primary': '#ffffff',
        'text-secondary': '#a2a5aa',
        'background-primary': '#000000',
        'background-secondary': '#1d2025',
        'background-neutral': '#303135',
      },
    },
    fontFamily: {
        'roboto': ["Roboto Mono", 'monospace']
    },
    borderRadius: {
      'none': '0',
      'sm': '4px',
      'md': '8px',
      'lg': '16px',
      'full': '9999px',
    },
    fontSize: {
      'sm': '14px',
      'lg': '16px',
      'xl': 'clamp(18px, 6vw, 20px)',
      '2xl': 'clamp(20px, 4vw, 25px)',
      '3xl': 'clamp(30px, 6vw, 40px)',
    }

  },
  plugins: [],
}

