
    
<!-- Blog Category -->
<?php if (is_array($PAGE->tags)){ ?>
    <div class="blog-head-tags">
        <?php $blogTags = array_slice($PAGE->tags, 0, 4); 
        foreach ($blogTags as $tag){ ?>
            <a href="<?= htmlspecialchars($tag) ?>">
                <span style="color: var(--colr-txt-2);"></span><?= htmlspecialchars($tag) ?>
            </a>
        <?php } ?>
    </div>
<?php } ?>

    

<!-- Blog Headline -->
<h1 class="blog-head-title" title="<?= htmlspecialchars($PAGE->title); ?>"><?= htmlspecialchars($PAGE->title); ?></h1>

<!-- Blog Meta -->
<div class="blog-head-meta">
    
    <a href="/@<?= htmlspecialchars($PAGE->username); ?>" title="<?= htmlspecialchars($PAGE->name); ?>" class="author-img-wrap">
        <figure>
            <img src="<?= htmlspecialchars($PAGE->profile_img); ?>#profile-img" alt="<?= htmlspecialchars($PAGE->name); ?>" loading="lazy" fetchpriority="low">
        </figure>
    </a>

    <div class="meta-info">

        <div class="meta-author">
            <a href="/@<?= htmlspecialchars($PAGE->username); ?>" title="@<?= htmlspecialchars($PAGE->username); ?>">
                <?= htmlspecialchars($PAGE->name); ?>
            </a>

            <span class="seperator"></span>

            <a href="" class="meta-follow"> Follow </a>

            <?php if($PAGE->verified){?>
                <i class="icon meta-verified" title="Verified">
                    <svg xmlns="http://www.w3.org/2000/svg" width="1.15em" viewBox="0 0 24 24" focusable="false" aria-hidden="true" style="pointer-events: none; display: inherit;"><path d="M12 2C6.5 2 2 6.5 2 12s4.5 10 10 10 10-4.5 10-10S17.5 2 12 2zM9.8 17.3l-4.2-4.1L7 11.8l2.8 2.7L17 7.4l1.4 1.4-8.6 8.5z"></path></svg>
                </i>
            <?php } ?>
        </div>

        <div class="meta-time">
            <time>
                <?= htmlspecialchars(date('F j, Y, g:i A', strtotime($PAGE->created_at))); ?>
            </time>
        </div>
    
    </div>
</div>


<figure class="blog-hero">
    <img src="<?= htmlspecialchars($blog->featured_img[0]); ?>#blog-hero" alt="<?= htmlspecialchars($PAGE->title); ?>"  loading="lazy" class="lazyload">
</figure>
