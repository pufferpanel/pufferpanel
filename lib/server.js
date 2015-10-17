/*
 * PufferPanel - Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var Rfr = require('rfr');
var Hapi = require('hapi');
var Nunjucks = require('nunjucks');
var Path = require('path');
var I18n = require('i18n');
var Fs = require('fs-extra');
var _ = require('underscore');
var RoutingApi = Rfr('lib/api/routing.js');
var SettingsController = Rfr('lib/controllers/settings.js');
var UserController = Rfr('lib/controllers/user.js');
var Configuration = Rfr('lib/config.js');
var Routes = Rfr('lib/routes/core.js');
var Logger = Rfr('lib/logger.js');

var hapiServer = new Hapi.Server();
var serverConfig = Configuration.server || { port: 3000 };
var yarConfig = Configuration.yar || { password: 'default' };
var Server = {};

/**
 * Prepares the Hapi server
 */
Server.prepare = function () {

    hapiServer.connection({ port: serverConfig.port });

    hapiServer.register(require('inert'), function (err) {

        if (err) {
            Logger.error('Failed to register \'inert\' plugin in server.');
        }
    });

    hapiServer.register({
        register: require('hapi-locals'),
        options: {
            isMergeArrays: false
        }
    }, function (err) {

        if (err) {
            Logger.error('Failed to register \'hapi-locals\' plugin in server.');
        }
    });

    hapiServer.register({
        register: require('h2o2')
    }, function (err) {

        if (err) {
            console.log('Failed to load h2o2');
        }
    });

    hapiServer.register({
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

    hapiServer.register({
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

    hapiServer.register(require('hapi-auth-cookie'), function (err) {

        if (err) {
            Logger.error('Failed to register \'hapi-auth-cookie\' plugin in server.');
        }

        hapiServer.auth.strategy('session', 'cookie', {
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

    I18n.configure({
        locales: ['en', 'it'],
        defaultLocale: 'en',
        directory: Path.join(__dirname, '../app/i18n'),
        cookie: 'pp_language',
        updateFiles: false,
        objectNotation: true
    });

    hapiServer.methods.locals('__', function (i18nRequest) {

        try {
            return I18n.__(i18nRequest);
        } catch (ex) {
            Logger.error('An exception occurred with i18n in the template.', ex);
            return '##missing language##';
        }
    });

    hapiServer.ext('onPreResponse', function (request, reply) {

        var response = request.response;
        if (!response.isBoom) {
            return reply.continue();
        }

        // Only Really Care about 5xx errors; not 404 and stuff.
        if ((response.output.statusCode >= 500 && response.stack) || Logger.runningVerbose === true) {
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

    //TODO: We really should probably not do all of this here, maybe have lib/routes/loader.js?
    hapiServer.route({
        method: 'GET',
        path: '/assets/{param*}',
        handler: {
            directory: {
                path: 'app/assets'
            }
        }
    });

    RoutingApi.register({
        method: 'GET',
        path: '/',
        handler: function (request, reply) {

            reply.redirect('/servers');
        },
        auth: 'session'
    });

    RoutingApi.register({
        method: 'GET',
        path: '/servers',
        handler: Routes.Base.get.servers,
        auth: 'session'
    });

    RoutingApi.register({
        method: 'GET',
        path: '/totp',
        handler: Routes.Base.get.totp,
        auth: 'session'
    });
    RoutingApi.register({
        method: 'POST',
        path: '/totp/generate-token',
        handler: Routes.Base.post.totp.generateToken,
        auth: 'session'
    });
    RoutingApi.register({
        method: 'POST',
        path: '/totp/verify-token',
        handler: Routes.Base.post.totp.verifyToken,
        auth: 'session'
    });
    RoutingApi.register({
        method: 'GET',
        path: '/account',
        handler: Routes.Base.get.account,
        auth: 'session'
    });
    RoutingApi.register({
        method: 'POST',
        path: '/account/update/{action}',
        handler: Routes.Base.post.accountUpdate,
        auth: 'session'
    });

    // Server Routes
    RoutingApi.register({
        method: '*',
        path: '/server/{server}/{path*}',
        handler: Routes.Server.handler,
        auth: 'session'
    });

    // Routes for Administration
    RoutingApi.register({
        method: 'GET',
        path: '/admin',
        handler: Routes.Admin.getIndex,
        auth: {
            strategy: 'session',
            scope: 'admin'
        }
    });

    RoutingApi.register({
        method: ['GET', 'POST'],
        path: '/admin/settings/{action}',
        handler: Routes.Admin.settings.router,
        auth: {
            strategy: 'session',
            scope: 'admin'
        }
    });

    RoutingApi.register({
        method: 'GET',
        path: '/admin/servers',
        handler: Routes.Admin.servers.getListAllServers,
        auth: {
            strategy: 'session',
            scope: 'admin'
        }
    });

    RoutingApi.register({
        method: 'GET',
        path: '/admin/servers/new',
        handler: Routes.Admin.servers.getAddNewServer,
        auth: {
            strategy: 'session',
            scope: 'admin'
        }
    });

// Routes for Authentication
    hapiServer.route([
        {
            method: 'GET',
            path: '/auth/login',
            config: {
                handler: Routes.Authentication.getLogin,
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
                handler: Routes.Authentication.postLogin,
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
                handler: Routes.Authentication.getLogout,
                auth: 'session'
            }
        }, {
            method: 'POST',
            path: '/auth/login/totp',
            config: {
                handler: Routes.Authentication.postTotp,
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
                handler: Routes.Authentication.getRegister
            }
        }, {
            method: 'POST',
            path: '/auth/register',
            config: {
                handler: Routes.Authentication.postRegister
            }
        }, {
            method: 'GET',
            path: '/auth/password',
            config: {
                handler: Routes.Authentication.getPassword
            }
        }, {
            method: 'POST',
            path: '/auth/password',
            config: {
                handler: Routes.Authentication.postPassword
            }
        }, {
            method: 'GET',
            path: '/auth/password/{token}',
            config: {
                handler: Routes.Authentication.getValidatePasswordRequest
            }
        }
    ]);

    process.env.NODE_TLS_REJECT_UNAUTHORIZED = 0;
    hapiServer.route({
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
};

/**
 * Finalizes Hapi server preparations.
 * This will register vision and registered routes from the API.
 */
Server.finalize = function () {

    hapiServer.register(require('vision'), function (err) {

        if (err) {
            Logger.error('Failed to register \'vision\' plugin in server.');
        }

        hapiServer.views({
            engines: {
                html: {
                    compile: function (template, options, next) {

                        SettingsController.getAllSettings(function (err, result) {

                            if (err) {
                                Logger.error('An error occurred while retrieving settings', err);
                                options.environment.addGlobal('settings', false);
                            }

                            options.environment.addGlobal('settings', result);
                            var compile = Nunjucks.compile(template, options.environment);

                            return next(null, function (context, opts, callback) {
                                return compile.render(context, callback);
                            });
                        });
                    },
                    prepare: function (options, next) {

                        options.compileOptions.environment = Nunjucks.configure(options.path);
                        return next();
                    }
                }
            },
            compileMode: 'async',
            compileOptions: {
                watch: true
            },
            isCached: false,
            path: RoutingApi._getViews()
        });
    });

    //load API declared routes
    var apiRoutes = RoutingApi._getRoutes();
    _.forEach(_.keys(apiRoutes), function (k) {

        var route = apiRoutes[k];
        hapiServer.route([
            {
                method: route.method,
                path: k,
                config: {
                    handler: route.callback,
                    auth: route.auth
                }
            }
        ]);
    });
};

/**
 * Starts the Hapi server
 */
Server.start = function () {

    hapiServer.start(function () {

        Logger.info('Server running at: ' + hapiServer.info.uri);
    });
};


module.exports = Server;
