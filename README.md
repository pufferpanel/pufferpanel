[![Circle CI](https://img.shields.io/circleci/project/PufferPanel/PufferPanel/kraken.svg?style=flat-square)](https://circleci.com/gh/PufferPanel/PufferPanel/tree/kraken)
[![](https://img.shields.io/david/PufferPanel/PufferPanel/kraken.svg?style=flat-square)](https://david-dm.org/PufferPanel/PufferPanel/kraken) [![](https://img.shields.io/david/dev/PufferPanel/PufferPanel/kraken.svg?style=flat-square)](https://david-dm.org/PufferPanel/PufferPanel/kraken#info=devDependencies&view=table)
[![](https://img.shields.io/coveralls/PufferPanel/PufferPanel/kraken.svg?style=flat-square)](https://coveralls.io/github/PufferPanel/PufferPanel?branch=kraken)
[![](https://img.shields.io/codacy/5181b766fb7d49a6bf47b3dabc93686c.svg?style=flat-square)](https://www.codacy.com/app/dane_2/PufferPanel)


[![](https://img.shields.io/github/license/PufferPanel/PufferPanel.svg?style=flat-square)](https://github.com/PufferPanel/PufferPanel/blob/kraken/LICENSE)

# PufferPanel — Written in Hapi
This is a development branch of PufferPanel written in Hapi, a Node.js server framework.

# Installing
To begin working with this development branch you will need to have `node` and `npm` installed on your system, as well as a working copy of the PufferPanel database. After that, simply clone this repository to a folder (making sure to switch to this branch). You will then need to create a `config.json` file in the root directory of this branch.

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
