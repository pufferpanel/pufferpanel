/*
 * PufferPanel — Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Rfr = require('rfr');
var Logger = Rfr('lib/logger.js');
var Hapi = require('hapi');
var Nunjucks = require('nunjucks');
var server = new Hapi.Server();
var Path = require('path');
var Configuration = Rfr('lib/config.js');
var Routes = Rfr('lib/routes/core.js');
var I18n = require('i18n');
var SettingsController = Rfr('lib/controllers/settings.js');
var UserController = Rfr('lib/controllers/user.js');
var Fs = require('fs-extra');

var serverConfig = Configuration.server || { port: 3000 };
var yarConfig = Configuration.yar || { password: 'default' };

server.connection({ port: serverConfig.port });

server.register(require('inert'), function (err) {

    if (err) {
        Logger.error('Failed to register \'inert\' plugin in server.');
    }
});

server.register({
    register: require('hapi-locals'),
    options: {
        isMergeArrays: false
    }
}, function (err) {

    if (err) {
        Logger.error('Failed to register \'hapi-locals\' plugin in server.');
    }
});

server.register({
    register: require('h2o2')
}, function (err) {

    if (err) {
        console.log('Failed to load h2o2');
    }

});

server.register({
    register: require('crumb'),
    options: {
        cookieOptions: {
            isSecure: false
        }
    }
}, function (err) {

    if (err) {
        Logger.error('Failed to register \'crumb\' plugin in server.');
    }
});

server.register({
    register: require('yar'),
    options: {
        cookieOptions: {
            password: yarConfig.password,
            isSecure: false,
            isHttpOnly: true
        }
    }
}, function (err) {

    if (err) {
        Logger.error('Failed to register \'yar\' plugin in server.');
    }
});

server.register(require('hapi-auth-cookie'), function (err) {

    if (err) {
        Logger.error('Failed to register \'hapi-auth-cookie\' plugin in server.');
    }

    server.auth.strategy('session', 'cookie', {
        password: yarConfig.password,
        cookie: 'pp_session',
        redirectTo: '/auth/login',
        isSecure: false,
        validateFunc: function (request, session, next) {

            UserController.getData(session.id, function (err, data) {

                if (err) {

                    Logger.error(err);
                    return next(err, false);
                }

                if (data.sessionToken !== session.sessionToken) {
                    return next(null, false);
                }

                return next(null, true, data);

            });

        }
    });

});

server.register(require('vision'), function (err) {

    if (err) {
        Logger.error('Failed to register \'vision\' plugin in server.');
    }

    server.views({
        engines: {
            html: {
                compile: function (src, options) {

                    var template = Nunjucks.compile(src, options.environment);
                    return function (context) {
                        return template.render(context);
                    };
                },
                prepare: function (options, next) {

                    options.compileOptions.environment = Nunjucks.configure(options.path, { watch: true });

                    SettingsController.getAllSettings(function (err, result) {

                        if (err) {

                            Logger.error('An error occurred while retrieving settings', err);
                            options.compileOptions.environment.addGlobal('settings', false);
                        }

                        options.compileOptions.environment.addGlobal('settings', result);
                        return next();

                    });

                }
            }
        },
        path: Path.join(__dirname, '../app/views')
    });
});

// i18n
I18n.configure({
    locales: ['en', 'it'],
    defaultLocale: 'en',
    directory: Path.join(__dirname, '../app/i18n'),
    cookie: 'pp_language',
    updateFiles: false,
    objectNotation: true
});

server.methods.locals('__', function (i18nRequest) {

    try {
        return I18n.__(i18nRequest);
    } catch (ex) {
        Logger.error('An exception occured with i18n in the template.', ex);
        return '##missing language##';
    }
});

server.ext('onPreResponse', function (request, reply) {

    var response = request.response;
    if (!response.isBoom) {
        return reply.continue();
    }

    // Only Really Care about 5xx errors; not 404 and stuff.
    if (response.output.statusCode >= 500 && response.stack) {
        Logger.warn(response.stack);
    }

    try {
        var stat = Fs.lstatSync(Path.join(__dirname, '../app/views/code/dynamic', response.output.statusCode.toString() + '.html'));
        return reply.view('code/dynamic/' + response.output.statusCode.toString(), response.output.payload);
    } catch (ex) {
        Logger.error('A server error occured but was not able to be rendered for the user.', ex.stack);
        return reply.view('code/500');
    }

});

// Handle public folder
server.route({
    method: 'GET',
    path: '/assets/{param*}',
    handler: {
        directory: {
            path: 'app/assets'
        }
    }
});

// Base Routes
server.route([
    {
        method: 'GET',
        path: '/',
        config: {
            handler: function (request, reply) {

                reply.redirect('/servers');
            },
            auth: 'session'
        }
    },
    {
        method: 'GET',
        path: '/servers',
        config: {
            handler: Routes.Base.get.servers,
            auth: 'session'
        }
    },
    {
        method: 'GET',
        path: '/totp',
        config: {
            handler: Routes.Base.get.totp,
            auth: 'session'
        }
    },
    {
        method: 'POST',
        path: '/totp/generate-token',
        config: {
            handler: Routes.Base.post.totp.generateToken,
            auth: 'session'
        }
    },
    {
        method: 'POST',
        path: '/totp/verify-token',
        config: {
            handler: Routes.Base.post.totp.verifyToken,
            auth: 'session'
        }
    },
    {
        method: 'GET',
        path: '/account',
        config: {
            handler: Routes.Base.get.account,
            auth: 'session'
        }
    },
    {
        method: 'POST',
        path: '/account/update/{action}',
        config: {
            handler: Routes.Base.post.accountUpdate,
            auth: 'session'
        }
    }
]);

// Server Routes
server.route([
    {
        method: 'GET',
        path: '/server/{server}',
        config: {
            handler: Routes.Server.get.index,
            auth: 'session'
        }
    }
]);

// Routes for Administration
server.route([
    {
        method: 'GET',
        path: '/admin',
        config: {
            handler: Routes.Admin.get.index,
            auth: {
                strategy: 'session',
                scope: 'admin'
            }
        }
    },
    {
        method: 'GET',
        path: '/admin/settings/{action}',
        config: {
            handler: Routes.Admin.get.settings.router,
            auth: {
                strategy: 'session',
                scope: 'admin'
            }
        }
    },
    {
        method: 'POST',
        path: '/admin/settings/general/updateCompanyName',
        config: {
            handler: Routes.Admin.post.settings.general.updateCompanyName,
            auth: {
                strategy: 'session',
                scope: 'admin'
            }
        }
    },
    {
        method: 'POST',
        path: '/admin/settings/general/generalSettings',
        config: {
            handler: Routes.Admin.post.settings.general.generalSettings,
            auth: {
                strategy: 'session',
                scope: 'admin'
            }
        }
    }
]);

// Routes for Authentication
server.route([
    {
        method: 'GET',
        path: '/auth/login',
        config: {
            handler: Routes.Authentication.get.login,
            auth: {
                mode: 'try',
                strategy: 'session'
            },
            plugins: {
                'hapi-auth-cookie': {
                    redirectTo: false
                }
            }
        }
    }, {
        method: 'POST',
        path: '/auth/login',
        config: {
            handler: Routes.Authentication.post.login,
            auth: {
                mode: 'try',
                strategy: 'session'
            },
            plugins: {
                'hapi-auth-cookie': {
                    redirectTo: false
                }
            }
        }
    }, {
        method: 'GET',
        path: '/auth/logout',
        config: {
            handler: Routes.Authentication.get.logout,
            auth: 'session'
        }
    }, {
        method: 'POST',
        path: '/auth/login/totp',
        config: {
            handler: Routes.Authentication.post.totp,
            auth: {
                mode: 'try',
                strategy: 'session'
            },
            plugins: {
                'hapi-auth-cookie': {
                    redirectTo: false
                }
            }
        }
    }, {
        method: 'GET',
        path: '/auth/register/{token*}',
        config: {
            handler: Routes.Authentication.get.register
        }
    }, {
        method: 'POST',
        path: '/auth/register',
        config: {
            handler: Routes.Authentication.post.register
        }
    }, {
        method: 'GET',
        path: '/auth/password',
        config: {
            handler: Routes.Authentication.get.password
        }
    }, {
        method: 'POST',
        path: '/auth/password',
        config: {
            handler: Routes.Authentication.get.password
        }
    }
]);

process.env.NODE_TLS_REJECT_UNAUTHORIZED = 0;
server.route({
    method: '*',
    path: '/scales/{server}/{path*}',
    config: {
        handler: {
            proxy: {
                mapUri: Routes.Scales.mapUri,
                onResponse: Routes.Scales.onResponse
            }
        },
        auth: 'session',
        plugins: {
            'hapi-auth-cookie': {
                redirectTo: false
            }
        }
    }
});

server.start(function () {

    Logger.info('Server running at: ' + server.info.uri);
});
