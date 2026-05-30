import js from '@eslint/js'
import ts from 'typescript-eslint'
import pluginVue from 'eslint-plugin-vue'

const browserGlobals = {
  window: 'readonly',
  document: 'readonly',
  localStorage: 'readonly',
  sessionStorage: 'readonly',
  navigator: 'readonly',
  location: 'readonly',
  history: 'readonly',
  setTimeout: 'readonly',
  clearTimeout: 'readonly',
  setInterval: 'readonly',
  clearInterval: 'readonly',
  requestAnimationFrame: 'readonly',
  cancelAnimationFrame: 'readonly',
  URL: 'readonly',
  URLSearchParams: 'readonly',
  Event: 'readonly',
  KeyboardEvent: 'readonly',
  PointerEvent: 'readonly',
  MouseEvent: 'readonly',
  HTMLElement: 'readonly',
  HTMLInputElement: 'readonly',
  HTMLDivElement: 'readonly',
  HTMLSpanElement: 'readonly',
  SVGElement: 'readonly',
  FileReader: 'readonly',
  IntersectionObserver: 'readonly',
  MutationObserver: 'readonly',
  ResizeObserver: 'readonly',
  performance: 'readonly',
  fetch: 'readonly',
  console: 'readonly',
  globalThis: 'readonly',
}

export default ts.config(
  { ignores: ['dist/', 'node_modules/'] },
  js.configs.recommended,
  ...ts.configs.recommended,
  ...pluginVue.configs['flat/essential'],
  {
    languageOptions: {
      globals: browserGlobals,
    },
  },
  {
    files: ['**/*.vue'],
    languageOptions: {
      globals: browserGlobals,
      parserOptions: { parser: ts.parser },
    },
  },
)
