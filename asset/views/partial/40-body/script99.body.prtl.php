<?php 
    use xet\Loc;
?>

<?php if (!empty($PAGE->jsInc99)){
    foreach ($PAGE->jsInc99 as $_js){ ?> <?= jslink($_js) ?>
<?php }}; ?>
