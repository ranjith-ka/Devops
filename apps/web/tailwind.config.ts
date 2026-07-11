import type { Config } from 'tailwindcss';
const config: Config = {
  darkMode: ['class'],
  content: ['./app/**/*.{ts,tsx}', './components/**/*.{ts,tsx}', './lib/**/*.{ts,tsx}'],
  theme: { extend: { boxShadow: { glow: '0 24px 80px rgba(0,0,0,0.32)' } } },
  plugins: []
};
export default config;
