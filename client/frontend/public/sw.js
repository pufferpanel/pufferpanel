self.addEventListener('install', () => {
  self.skipWaiting()
})

self.addEventListener('activate', () => {
  self.clients.matchAll({
    type: 'window'
  }).then(clients => {
    for (let client of clients) {
      client.navigate(client.url)
    }
  })
})
