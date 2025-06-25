/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}", // Isso é crucial para o Tailwind escanear seus arquivos React
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}