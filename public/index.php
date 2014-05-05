<?php 
require(dirname(__DIR__).'/src/framework/framework.core.php');

$app = new Silex\Application();
$app['debug'] = true;

$app->register(new Silex\Provider\TwigServiceProvider(), [
	'twig.path' => dirname(__DIR__).'/app/views'
]);

$app['twig']->addGlobal('lang', $_l->loadTemplates());
$app['twig']->addGlobal('settings', $core->settings->get());
$app['twig']->addGlobal('get', Page\components::twigGET());
if($core->user->getData('root_admin') == 1){ $app['twig']->addGlobal('admin', true); }

$app->get('/', function() use ($app)
{
	return $app['twig']->render('panel/index.html', [
		'footer' => array(
			'queries' => Database\databaseInit::getCount(),
			'seconds' => number_format((microtime(true) - $pageStartTime), 4)
		)
	]);
});

$app->run();