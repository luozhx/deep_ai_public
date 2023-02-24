module.exports = {
    transpileDependencies: [
        'vue-echarts',
        'resize-detector'
    ],
    devServer: {
        proxy: {
            '^/api/v1/*': {
                target: 'http://172.18.167.37:30000',
                // target: 'http://127.0.0.1:1323',
                // target: 'http://172.18.167.37:1323',
                ws: true,
                changeOrigin: true
            }
        }
    }
}
