const errors = {
  notlinked: 'users.OidcErrorNotLinked',
  taken: 'users.OidcErrorTaken',
  denied: 'users.OidcErrorDenied',
  state: 'users.OidcErrorState',
  unavailable: 'users.OidcErrorUnavailable',
  nologin: 'users.OidcErrorNoLogin'
}

// the reasons the OIDC callback sends to the error page and back to the security settings
export function oidcErrorMessage(t, reason) {
  return t(errors[reason] || 'users.OidcErrorGeneric')
}
