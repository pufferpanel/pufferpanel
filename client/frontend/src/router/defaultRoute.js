export default function(api) {
  if (!api.auth.isLoggedIn()) {
    return { name: 'Login' }
  }

  if (api.auth.hasScope('servers.view')) {
    return { name: 'ServerList' }
  } else if (api.auth.hasScope('templates.view')) {
    return { name: 'TemplateList' }
  } else if (api.auth.hasScope('users.view')) {
    return { name: 'UserList' }
  } else if (api.auth.hasScope('nodes.view')) {
    return { name: 'NodeList' }
  }

  return { name: 'Self' }
}
