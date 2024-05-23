const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  publicPath: './', // 配置项目根路径
  outputDir: '../static', // 配置打包文件输出路径
  transpileDependencies: true
})
