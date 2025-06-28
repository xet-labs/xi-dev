<?php
use xet\Loc;

require(Loc::file('PRTL', 'head'));
?>

<body>
	<?php
	require(Loc::file('PRTL', 'lib.body'));
	require(Loc::file('PRTL', 'script.body'));
	require(Loc::file('PRTL','nav'));
	!empty($PAGE->prtl_stickySocial) ? require(Loc::file('PRTL','sticky-social')) : '';
	?>


	<main>
		<?php
		// Get the file and directory information of the calling file
		$callFilePath = debug_backtrace()[0]['file'];
		$callFile = basename($callFilePath);
		$callPath = dirname($callFilePath);
		$callDir = basename(dirname($callFilePath));

		if (!empty($PAGE->cnt)) {
			echo htmlspecialchars_decode($PAGE->cnt);
		
		} elseif (!empty($blog->cntFile) && file_exists($blog->cntFile)) {
			require_once($blog->cntFile);
		
		} elseif (file_exists($callPath . '/' . pathinfo($callFile, PATHINFO_FILENAME) . '.cnt.php')) {
			require($callPath . '/' . pathinfo($callFile, PATHINFO_FILENAME) . '.cnt.php');
			
		} elseif (file_exists($callPath . '/cnt.php')) {
			require($callPath . '/cnt.php');
		}
		?>

	</main>


	<?php
	require(Loc::FILE('PRTL', 'footer'));
	require(Loc::FILE('PRTL', 'script99.body'));
	require(Loc::FILE('PRTL', 'lib99.body'));
	?>
	
</body>

</html>