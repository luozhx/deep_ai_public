module.exports = {
    root: true,
    env: {
        node: true
    },
    'extends': [
        'plugin:vue/essential',
        '@vue/standard',
        '@vue/typescript'
    ],
    rules: {
        'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
        'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
        'object-curly-spacing': ['error', 'never'],
        'space-before-function-paren': ['error', 'never'],
        '@typescript-eslint/indent': 'off',
        'indent': 'off',
        'vue/html-indent': 'off',
        'no-trailing-spaces': ["error", { "skipBlankLines": true }]
    },
    parserOptions: {
        parser: '@typescript-eslint/parser'
    }
}
