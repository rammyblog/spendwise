/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ['./templates/**/*.{html,gohtml,js}'],
  theme: {
    extend: {
      fontFamily: {
        sans: ['Inter', 'ui-sans-serif', 'system-ui', 'sans-serif'],
      },
    },
  },
  plugins: [],
};
