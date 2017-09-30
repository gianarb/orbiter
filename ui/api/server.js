const jsonServer = require('json-server')
const database = require('./db.json')

const server = jsonServer.create()
const router = jsonServer.router(database)
const middleware = jsonServer.defaults()
const rewriter = jsonServer.rewriter({
  '/handle/:autoscaler_name/:service_name/inspect-service': '/services-inspect/:service_name',
  '/handle/:autoscaler_name/:service_name/up': '/services-scale-up',
  '/handle/:autoscaler_name/:service_name/down': '/services-scale-down'
})

server.use(rewriter)
server.use(middleware)
server.use(router)
server.listen(3001, () => {
  console.log('JSON Server is running...', '\n')
})
