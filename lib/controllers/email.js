/*
 * PufferPanel - Reinventing the way game servers are managed.
 * Copyright (c) 2015 PufferPanel
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 */
var _ = require('underscore');
var Rfr = require('rfr');
var Path = require('path');
var Dns = require('dns');
var Nodemailer = require('nodemailer');
var Mailcomposer = require('mailcomposer');
var Nunjucks = require('nunjucks');
var Async = require('async');
var Transports = {
    SMTP: require('nodemailer-smtp-transport'),
    Postmark: require('postmark'),
    Sendgrid: require('sendgrid'),
    Mailgun: require('mailgun').Mailgun
};
var SettingsController = Rfr('lib/controllers/settings.js');
var AuthenticationController = Rfr('lib/controllers/authentication.js');

/** @namespace */
var EmailController = {};

EmailController.send = function (to, subject, message, next) {

    SettingsController.get('emailConfiguration', function (err, data) {

        if (err) {
            return next(err);
        }

        var transporter;
        if (data.method === 'postmark') {
            EmailController.sendViaPostmark(to, subject, message, data, next);
        } else if (data.method === 'mailgun') {
            EmailController.sendViaMailgun(to, subject, message, data, next);
        } else if (data.method === 'sendgrid') {
            EmailController.sendViaSendgrid(to, subject, message, data, next);
        } else if (data.method === 'smtp') {

            // Send via SMTP Server
            transporter = Nodemailer.createTransport(Transports.SMTP({
                host: data.smtp.host,
                port: data.smtp.port,
                auth: {
                    user: data.smtp.user,
                    pass: AuthenticationController.decrypt(data.smtp.pass)
                }
            }));

            transporter.sendMail({
                from: data.email,
                to: to,
                subject: subject,
                text: message
            });

            return next(null);

        } else if (data.method === 'direct') {

            transporter = nodemailer.createTransport();

            transporter.sendMail({
                from: data.email,
                to: to,
                subject: subject,
                text: message
            });

            return next(null);

        } else {
            return next(new Error('The method selected for sending email is not valid.'));
        }

    });

};

EmailController.sendTemplate = function (template, to, data, next) {

    SettingsController.getAllSettings(function (err, settings) {

        if (err) {
            return next(err);
        }

        Async.series([
            function (callback) {
                if (data.ip) {
                    Dns.reverse(data.ip, function (err, response) {

                        if (err) {
                            return callback(err);
                        }

                        return callback(null, response.toString());

                    });
                } else {
                    return callback();
                }
            }
        ], function (err, dnsResult) {

            if (err) {
                return next(err);
            }

            data = _.extend(data, {
                email: to,
                companyName: settings.companyName,
                mainUrl: settings.urls.main,
                ipDetails: dnsResult || 'unknown'
            });

            Nunjucks.configure(Path.join(__dirname, '../../app/emails'));
            EmailController.send(to, data.subject || 'untitled', Nunjucks.render(template + '.html', data), function (err) {
                return next(err);
            });

        });

    });

};

EmailController.sendViaPostmark = function (to, subject, message, data, next) {

    var Client = new Transports.Postmark.Client(data.apiToken);

    Client.sendEmail({
        'From': data.email,
        'To': to,
        'Subject': subject,
        'HtmlBody': message
    }, function (err) {
        return next(err);
    });

};

EmailController.sendViaMailgun = function (to, subject, message, data, next) {

    var Composed = Mailcomposer({
        from: data.email,
        sender: data.email,
        to: to,
        subject: subject,
        html: message
    });

    var Client = new Transports.Mailgun(data.apiToken);

    Client.sendRaw(data.email, to, Composed, function (err) {
        return next(err);
    });

};

EmailController.sendViaSendgrid = function (to, subject, message, data, next) {

    var Client = Transports.Sendgrid(data.apiToken);
    var Email = new Client.Email({
        to: to,
        from: data.email,
        subject: subject,
        html: message
    });

    Client.send(Email, function (err) {
        return next(err);
    });

};

module.exports = EmailController;
