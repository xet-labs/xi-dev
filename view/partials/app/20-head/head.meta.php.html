<?php
// dd($blog->featured_img);
function issetx($obj, $property)
{
  return isset($obj) && property_exists($obj, $property) && !is_null($obj->$property) && $obj->$property !== false;
}


// $PAGE->app = rtrim(url('/'), '/');
$PAGE->url = rtrim(url()->current(), '/');
$PAGE->app = url('/');
$PAGE->appName = !empty($PAGE->appName) ? $PAGE->appName : config('app.name');
$PAGE->appLogo = $PAGE->app . '/res/static/brand/favicon.svg';
$PAGE->appImg = $PAGE->app . '/res/static/brand/brand.svg';

$PAGE->type = htmlspecialchars(!empty($PAGE->type) ? $PAGE->type : 'WebSite');
$PAGE->canonical = issetx($PAGE, 'canonical') ? htmlspecialchars($PAGE->canonical) : $PAGE->url;

$PAGE->description = htmlspecialchars(!empty($PAGE->description) ? $PAGE->description : (!empty($PAGE->excerpt) ? $PAGE->excerpt : ''));
$PAGE->featured_img = !empty($PAGE->featured_img) ? (is_array($PAGE->featured_img) ? array_map(fn($img) => htmlspecialchars(asset($img)), $PAGE->featured_img) : htmlspecialchars(asset($PAGE->featured_img))) : '';


if (in_array($PAGE->type, ['BlogPosting', 'NewsArticle', 'Article'])) {
  $PAGE->ogType = "article";
  $PAGE->author = htmlspecialchars($PAGE->author ?? $PAGE->name ?? 'Guest Author');
  $PAGE->authorUrl = htmlspecialchars(empty($PAGE->authorUrl) ? (!empty($PAGE->username) ? $PAGE->app . '/@' . $PAGE->username : '') : $PAGE->authorUrl);
  $PAGE->profile_img = htmlspecialchars(!empty($PAGE->profile_img) ? asset($PAGE->profile_img) : '');
} elseif ($PAGE->type === 'Product') {
  $PAGE->ogType = "product";
} else {
  $PAGE->ogType = 'website';
}

$PAGE->metaTitle = htmlspecialchars(!empty($PAGE->title) ? $PAGE->title . (!empty($PAGE->author) ? ' | by ' . $PAGE->author : '') . ' | ' . $PAGE->appName : $PAGE->appName);
?>


<!-- META -->
<title><?= $PAGE->metaTitle ?></title>

<meta name="robots" content="index, follow" />
<meta name="referrer" content="no-referrer-when-downgrade" />
<?= issetx($PAGE, 'canonical') ? '<link rel="canonical" href="' . htmlspecialchars($PAGE->canonical) . '" />' : '' ?>

<?= '<meta name="title" content="' . $PAGE->metaTitle . '" />' ?>
<?= !empty($PAGE->description) ? '<meta name="description" content="' . htmlspecialchars($PAGE->description) . '" />' : ''; ?>
<?= !empty($PAGE->tags) ? '<meta name="keywords" content="' . implode(', ', $PAGE->tags) . '" />' : ''; ?>
<?= !empty($PAGE->ogType) ? '<meta property="og:type" content="' . $PAGE->ogType . '" />' : ''; ?>
<?= issetx($PAGE, 'appName') ? '<meta property="og:site_name" content="' . $PAGE->appName . '" />' : ''; ?>
<?= !empty($PAGE->title) ? '<meta property="og:title" content="' . htmlspecialchars($PAGE->title) . '" />' : ''; ?>
<?= !empty($PAGE->description) ? '<meta property="og:description" content="' . htmlspecialchars($PAGE->description) . '" />' : ''; ?>
<?= issetx($PAGE, 'canonical') ? '<meta property="og:url" content="' . htmlspecialchars($PAGE->canonical) . '" />' : ''; ?>

<?= !empty($PAGE->featured_img)
  ? '<meta property="og:image" content="' . $PAGE->featured_img[0] . '" />' .
  (!empty($PAGE->featured_imgAlt)
    ? '<meta property="og:image:alt" content="' . htmlspecialchars($PAGE->featured_imgAlt) . '" />'
    : '<meta property="og:image:alt" content="' . htmlspecialchars($PAGE->title) . '" />')
  : '';
?>

<?= !empty($PAGE->title) ? '<meta name="twitter:title" content="' . htmlspecialchars($PAGE->title) . '">' : ''; ?>
<?= !empty($PAGE->description) ? '<meta name="twitter:description" content="' . htmlspecialchars($PAGE->description) . '">' : ''; ?>
<?= issetx($PAGE, 'canonical') ? '<meta name="twitter:url" content="' . htmlspecialchars($PAGE->canonical) . '">' : ''; ?>
<?= !empty($PAGE->x['site']) ? '<meta name="twitter:site" content="' . htmlspecialchars($PAGE->x['site']) . '" />' : ''; ?>
<?= !empty($PAGE->featured_img) ? '<meta name="twitter:image:src" content=' . $PAGE->featured_img[0] . '><meta name="twitter:card" content="summary_large_image" />' : ''; ?>
<?= !empty($PAGE->featured_imgAlt) ? '<meta name="twitter:image:alt"  content="' . htmlspecialchars($PAGE->featured_imgAlt) . '">' : ''; ?>
<?= !empty($PAGE->x['creator']) ? '<meta name="twitter:creator" content="' . htmlspecialchars($PAGE->x['creator']) . '" />' : ''; ?>
<?= !empty($PAGE->site_domain) ? '<meta property="twitter:domain" value="' . htmlspecialchars($PAGE->site_domain) . '" />' : ''; ?>

<?php
if (in_array($PAGE->type, ['BlogPosting', 'NewsArticle', 'Article'])) { ?>
  <?= !empty($PAGE->username) ? '<meta name="author" content="' . htmlspecialchars($PAGE->author) . '" />' : ''; ?>
  <?= !empty($PAGE->username) ? '<meta property="article:author" content="' . htmlspecialchars($PAGE->authorUrl) . '" />' : ''; ?>
  <?= !empty($PAGE->authorUrl) ? '<link rel="author" href="' . htmlspecialchars($PAGE->authorUrl) . '" />' : ''; ?>
  <?= !empty($PAGE->created_at) ? '<meta property="article:published_time" content="' . htmlspecialchars($PAGE->created_at) . '" />' : ''; ?>

  <?= !empty($PAGE->author) ? '<meta name="twitter:label1" content="Written by" /><meta name="twitter:data1" content="@' . htmlspecialchars($PAGE->author) . '" />' : ''; ?>
  <?= !empty($PAGE->x['category']) ? '<meta name="twitter:label2" value="Category" /><meta name="twitter:data2" value="' . htmlspecialchars($PAGE->x['category']) . '" />' : ''; ?>
  <?= !empty($PAGE->created_at) ? '<meta name="twitter:label3" value="Published on" /><meta name="twitter:data3" value="' . htmlspecialchars($PAGE->created_at) . '" />' : ''; ?>
  <?= !empty($PAGE->reading_time) ? '<meta name="twitter:label4" value="Reading time" /><meta name="twitter:data4" value="' . htmlspecialchars($PAGE->reading_time) . '" />' : ''; ?>
<?php } ?>

<script type="application/ld+json">
  {
    "@context": "https://schema.org",
    "@type": "<?= $PAGE->type ?>",
    <?= !empty($PAGE->ld) ? $PAGE->ld . ',' : ''; ?> "url": "<?= $PAGE->url ?>",
    "mainEntityOfPage": "<?= $PAGE->url ?>",
    <?= !empty($PAGE->description) ? '"description": "' . $PAGE->description . '",' : ''; ?>
    <?= !empty($PAGE->tags) ? '"keywords": ' . json_encode($PAGE->tags) . ',' : ''; ?>
    <?= !empty($PAGE->featured_img) ? '"image": ' . json_encode($PAGE->featured_img) . ',' : ''; ?>

    <?php if (in_array($PAGE->type, ['BlogPosting', 'NewsArticle', 'Article'])) { ?> "headline": "<?= $PAGE->metaTitle ?>",
      "dateCreated": "<?= $PAGE->created_at instanceof \Carbon\Carbon ? $blog->created_at->setTimezone('UTC')->format('Y-m-d\TH:i:s.v\Z') : (new \Carbon\Carbon($blog->created_at))->setTimezone('UTC')->format('Y-m-d\TH:i:s.v\Z'); ?>",
      "datePublished": "<?= $PAGE->created_at instanceof \Carbon\Carbon ? $blog->created_at->setTimezone('UTC')->format('Y-m-d\TH:i:s.v\Z') : (new \Carbon\Carbon($blog->created_at))->setTimezone('UTC')->format('Y-m-d\TH:i:s.v\Z'); ?>",
      "dateModified": "<?= $PAGE->updated_at instanceof \Carbon\Carbon ? $blog->updated_at->setTimezone('UTC')->format('Y-m-d\TH:i:s.v\Z') : (new \Carbon\Carbon($blog->updated_at))->setTimezone('UTC')->format('Y-m-d\TH:i:s.v\Z'); ?>",
      "author": {
        "@type": "Person",
        "name": "<?= $PAGE->author ?>",
        <?= !empty($PAGE->profile_description) ? '"description": "' . $PAGE->profile_description . '",' : ''; ?>
        <?= !empty($PAGE->profile_img) ? '"image": "' . $PAGE->profile_img . '",' : ''; ?>
        <?= !empty($PAGE->jobTitle) ? '"jobTitle": "' . $PAGE->jobTitle . '",' : ''; ?>
        <?= !empty($PAGE->sameAs) ? '"sameAs": "' . $PAGE->sameAs . '",' : ''; ?> "url": "<?= $PAGE->authorUrl ?>"
      },

      "articleSection": <?= json_encode($PAGE->tags) ?>,
      "creator": ["<?= $PAGE->author ?>"],
      "editor": "<?= $PAGE->author ?>",
    <?php } ?>
    <?= (isset($PAGE->isAccessibleForFree) && $PAGE->isAccessibleForFree === false) ? '"isAccessibleForFree": false,' : '"isAccessibleForFree": true,'; ?>


    "publisher": {
      "@type": "Organization",
      "name": "<?= !empty($PAGE->ldPublisherName) ?  $PAGE->ldPublisherName : $PAGE->appName ?>",
      <?= !empty($PAGE->altName) ? '"alternateName": ' . json_encode($PAGE->altName) . ',' : ''; ?> "url": "<?= !empty($PAGE->ldPublisherUrl) ?  $PAGE->ldPublisherUrl : $PAGE->app ?>",
      "logo": {
        "@type": "ImageObject",
        "url": "<?= $PAGE->appImg ?>"
      }
    }
    <?= !empty($PAGE->ld99) ? ',' . $PAGE->ld99 : ''; ?>

  }
</script>

