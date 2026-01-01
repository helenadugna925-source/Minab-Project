/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./node_modules/flowbite/**/*.{js,ts}",
    "./components/**/*.{js,vue,ts}",
    "./layouts/**/*.vue",
    "./pages/**/*.vue",
    "./plugins/**/*.{js,ts}",
    "./app.vue",
    "./error.vue",
  ],
  theme: {
    extend: {
      colors: {
        primary: 'var(--color-primary)',
        accent: 'var(--color-accent)',
        'bg-default': 'var(--color-background)',
        'text-default': 'var(--color-text)',
        muted: 'var(--color-muted)'
      }
    },
  },
  plugins: [
    require('flowbite/plugin')
  ],
}

