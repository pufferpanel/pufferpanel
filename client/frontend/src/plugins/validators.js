const validate = {
  username(username) {
    return username.length >= 5
  },
  email(email) {
    return /^.+@.+\..{2,}$/.test(email)
  },
  password(password) {
    return password.length >= 8
  }
}

export default {
  install: (app) => {
    app.config.globalProperties.$validate = validate
    app.provide('validate', validate)
  }
}
