<?php 
session_start();
use xet\Loc;

$PAGE = (object) array_merge(
    !empty($PAGE) ? (array) $PAGE : [], 
    !empty($blog) ? $blog->toArray() : []
);

function csslink($cssUrl, $load = 0) { 
    if ($load == 0) {
        return '
        <link rel="preload" as="style" href="' . htmlspecialchars($cssUrl, ENT_QUOTES, 'UTF-8') . '">
        <link rel="stylesheet" href="' . htmlspecialchars($cssUrl, ENT_QUOTES, 'UTF-8') . '">
        ';
    } elseif ($load == 1) {
        return '
        <link rel="' . htmlspecialchars($cssUrl, ENT_QUOTES, 'UTF-8') . '">
        ';
    }
}
function jslink($jsUrl, $load = 0) { 
    if ($load == 0) {
        return '
        <link rel="preload" as="script" href="' . htmlspecialchars($jsUrl, ENT_QUOTES, 'UTF-8') . '" crossorigin="anonymous">
		<script defer src="' . htmlspecialchars($jsUrl, ENT_QUOTES, 'UTF-8') . '" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
        ';
    } elseif ($load == 1) {
        return '
        <script src="' . htmlspecialchars($jsUrl, ENT_QUOTES, 'UTF-8') . '"></script>
        ';
    }
}
?>

<!DOCTYPE html>
<html lang="en">
<script>document.documentElement.classList.add(localStorage.getItem("theme"));</script>
<head>
	<meta charset="UTF-8">
	<meta http-equiv="X-UA-Compatible" content="IE=edge" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />
	
	<link rel="icon" type="image/svg+xml" href="/favicon.svg" />
	<link rel="icon" type="image/png" sizes="96x96" href="<= asset(/res/static/brand/favicon-192x192.png)?>" />
	<link rel="icon" type="image/png" sizes="48x48" href="/res/static/brand/favicon-180x180.png" />
	<link rel="icon" type="image/png" sizes="16x16" href="/res/static/brand/favicon-48x48.png" />
	<link rel="apple-touch-icon" href="/res/static/brand/favicon-180x180.png" />
	<link rel="icon" type="image/x-icon" href="/favicon.ico" />
	<link rel="manifest" href="/manifest.json" />
	<meta name="theme-color" content="#ff9700">

	<!-- FONTS -->
	<link rel="preload" href="/res/static/fonts/Inter/Inter-VariableFont_slnt_wght.woff2" as="font" type="font/woff2" crossorigin="anonymous">

	<!-- CSS -->
	<?= csslink(route('app.css')) ?>

	<?php if (!empty($PAGE->link) && $PAGE->link !== false) {
		if (is_array($PAGE->link)) { 
			foreach ($PAGE->link as $link) { ?><?= csslink(htmlspecialchars($link)); ?><?php }
		} else { ?><?= csslink(htmlspecialchars($PAGE->link)); ?><?php } ?>
	<?php } ?>
	

	<?php
	require(Loc::file('PRTL', 'meta.head'));
	require(Loc::file('PRTL', 'lib.head'));
	require(Loc::file('PRTL', 'script.head'));
	?>

	<meta name="csrf-token" content="<?= csrf_token(); ?>">

</head>