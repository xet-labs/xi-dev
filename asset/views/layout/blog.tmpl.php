<?php
use xet\Loc;

require(Loc::file('PRTL', 'head'));
$imgSrc = $PAGE->imgSrc;
?>


<body class="body-blog">

	<?php
	require(Loc::file('PRTL', 'lib.body'));
	require(Loc::file('PRTL', 'script.body'));
	$currentMenu = $subBrand = 'Blog';
	require(Loc::file('PRTL', 'nav'));
	// !empty($PAGE->prtl_stickySocial) ? require(Loc::file('PRTL','sticky-social')) : '';
	?>

	<main>
	<div class="blog-wrap wrap">

	<article class="blog con publisheds">

		<header class="blog-head">
			<?php require(Loc::file('PRTL', 'head.blog')); ?>
		</header>

		<section class="blog-cnt published">
		<?php
		if (!empty($PAGE->cnt)) { echo htmlspecialchars_decode($PAGE->cnt); }
		elseif (!empty($blog->cntFile) && file_exists($blog->cntFile)) { require_once($blog->cntFile); }
		?>
		</section>
		
		<div class="blog-foot">
			<hr class="line-x">
		</div>
		
	</article>

	</div>
	</main>

	<?php
	require(Loc::FILE('PRTL', 'footer'));
	require(Loc::FILE('PRTL', 'script99.body'));
	require(Loc::FILE('PRTL', 'lib99.body'));
	?>

</body>

</html>