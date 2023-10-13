export default function(api) {
  if (!api.auth.isLoggedIn()) {
    return { name: 'Login' }
  }

  return { name: 'ServerList' }
}
