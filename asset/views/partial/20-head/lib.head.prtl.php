<!-- Lib -->
<?= !empty($PAGE->lib_hljs) ? csslink('/res/lib/hljs/styles/atom-one-dark.css') . jslink('/res/lib/hljs/hl.min.js', 1) . '<script>hljs.highlightAll();</script>' : ''; ?> 

<?= !empty($PAGE->lib_prism) ? csslink('/res/lib/prism/prism.css') : ''; ?> 


<?= !empty($PAGE->lib->tw) ? csslink('/res/lib//tailwind/tailwind.css') : ''; ?>
<?= !empty($PAGE->lib->fa) ? csslink('https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.5.1/css/all.min.css') : ''; ?>