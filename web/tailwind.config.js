/** @type {import('tailwindcss').Config} */
export default {
    content: [
        './templates/**/*.templ',
        './**/*.html',
        './**/*.go',
    ],
    darkMode: 'class',
  theme: {
    extend: {
      colors: {
        indigo: { 500: '#6366F1' },
        purple: { 600: '#7C3AED' },
      },
      fontSize: {
        h1: '3rem',
        h2: '2.5rem',
      },
    },
  },
  plugins: [
    require('@tailwindcss/typography') // optional, if you want typography plugin
  ],
};

