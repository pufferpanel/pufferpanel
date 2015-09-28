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
var Nodemailer = require('nodemailer');
var Mailcomposer = require('mailcomposer');
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

            return next();

        } else if (data.method === 'direct') {

            transporter = nodemailer.createTransport();

            transporter.sendMail({
                from: data.email,
                to: to,
                subject: subject,
                text: message
            });

            return next();

        } else {
            return next(new Error('The method selected for sending email is not valid.'));
        }

    });

};

EmailController.sendTemplate = function (template, data, next) {

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

// Handle Logic for prebuiltAuth, since we can't have anything that is standardized...

var transporter = Nodemailer.createTransport({
    service: 'postmark',
    auth: prebuiltAuth
});

var emailOptions = {
    from: 'ptero@dactyl.link',
    to: 'fly@pterodactyl.io',
    subject: 'CHECK OUT MY WINGS!',
    html: 'holla!'
};

transporter.sendMail(emailOptions, function (err, message) {

});
