const { defineConfig } = require('@vue/cli-service');
module.exports = defineConfig({
  transpileDependencies: true,
  lintOnSave: false,
  devServer: {
    proxy: {
      '/api': {
        target: 'http://172.17.77.197:8090',
        changeOrigin: true,
        pathRewrite: {
          '^/api': '/api', // 如果接口路径有/api前缀，可以使用pathRewrite将其去掉
        },
      },
    },
  },
});
