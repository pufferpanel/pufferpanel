[![Circle CI](https://circleci.com/gh/PufferPanel/PufferPanel/tree/hapi.js.svg?style=svg)](https://circleci.com/gh/PufferPanel/PufferPanel/tree/hapi.js)

# PufferPanel — Written in Hapi
This is a development branch of PufferPanel written in Hapi, a Node.js server framework.

# Installing
To begin working with this development branch you will need to have `node` and `npm` installed on your system, as well as a working copy of the PufferPanel database. After that, simply clone this repository to a folder (making sure to switch to this branch). You will then need to create a `configuration.json` file in the root directory of this branch.

You must have `rethinkdb` installed on your system in order to run this version.

```json
{
  "server": {
    "port": 3000
  },
  "rethink": {
    "host": "localhost",
    "port": 28015,
    "database": "pufferpanel"
  },
  "yarPassword": "<some_password>"
}
```

After that, execute the commands below to get the server up and running.
```
npm install
node index.js
```

*Please do not report bugs in this branch to the issue tracker!*
