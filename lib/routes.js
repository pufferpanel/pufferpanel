var Hapi = require('hapi')
var Nunjucks = require('nunjucks')
var server = new Hapi.Server()
var Path = require('path')

server.connection({ port: 3000 })

server.register(require('inert'), function (err) {
  if (err) {
    console.log('Failed to load inert.')
  }
})

server.register(require('vision'), function (err) {

  if (err) {
    console.log('Failed to load vision.')
  }

  server.views({
    engines: {
      html: {
        compile: function (src, options) {
          var template = Nunjucks.compile(src, options.environment)
          return function (context) {
            return template.render(context)
          }
        },
        prepare: function (options, next) {
          options.compileOptions.environment = Nunjucks.configure(options.path, { watch: false })
          return next()
        }
      }
    },
    path: Path.join(__dirname, '../app/views')
  })

})

// Handle public folder
server.route({
  method: 'GET',
  path: '/public/{param*}',
  handler: {
    directory: {
      path: 'public'
    }
  }
})

server.route({
  method: 'GET',
  path: '/auth/login',
  handler: function (request, reply) {
    reply.view('auth/login')
  }
})

server.route({
  method: 'GET',
  path: '/two',
  handler: function (request, reply) {
    reply('Hello World 2')
  }
})

server.start(function () {
  console.log('Server running at: ', server.info.uri)
})
