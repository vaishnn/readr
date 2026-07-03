/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./src/**/*.{html,ts}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        brand: {
          50:  '#f0f4ff',
          100: '#e0e9ff',
          400: '#6b8cff',
          500: '#4f6ef7',
          600: '#3b55e0',
          700: '#2d42c2',
        },
        // Warm golden-yellow accent — #FDBA31 per design spec
        accent: {
          300: '#FFD166',
          400: '#FDBA31',
          500: '#E8A820',
          600: '#CC9210',
        },
        // Override the default slate ramp with the dark teal-gray palette from the design spec.
        // #1A2026 (panels/sidebar), #232B32 (cards), #9CA3AF (secondary text).
        // Every existing slate-* class picks this up automatically — no template changes needed.
        slate: {
          50:  '#F6F7F9',
          100: '#EEF0F3',
          200: '#DDE1E6',
          300: '#C4CBD4',
          400: '#9CA3AF',  // secondary text — matches spec exactly
          500: '#5A6A7A',
          600: '#3D4F5E',
          700: '#2D3744',  // borders, hover states
          800: '#232B32',  // card / input backgrounds — matches spec
          900: '#1A2026',  // sidebar / header panels — matches spec
          950: '#12181E',  // page background
        },
      },
      fontFamily: {
        sans: ['Inter', 'ui-sans-serif', 'system-ui'],
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography'),
  ],
};
