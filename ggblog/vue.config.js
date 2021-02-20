module.exports = {
 
  devServer: {
    port: 8888,
    open: true,
    overlay: {
      warnings: false,
      errors: true
    },
    proxy: {
      // change xxx-api/login => mock/login
      // detail: https://cli.vuejs.org/config/#devserver-proxy
  //  [process.env.VUE_APP_BASE_API]: {
    '/blog/api': {
      //  target: `http://127.0.0.1:${port}/mock`,
      target: `http://5.123.com`,
        changeOrigin: true,
        pathRewrite: {
         // ['^' + process.env.VUE_APP_BASE_API]: ''
         '^/blog/api':'/blog/api'
        },
         secure: false 

      }
    } 




       
    },
    lintOnSave: false
}
