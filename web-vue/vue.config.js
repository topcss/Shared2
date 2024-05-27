const { defineConfig } = require('@vue/cli-service')

let plugins = []
const MonacoWebpackPlugin = require('monaco-editor-webpack-plugin')
plugins.push(new MonacoWebpackPlugin({ languages: ['javascript', 'java', 'css', 'html'] }));

module.exports = defineConfig({
  publicPath: './', // 配置项目根路径
  outputDir: '../static', // 配置打包文件输出路径
  transpileDependencies: true,
  configureWebpack: {
    plugins: plugins
  }
})
