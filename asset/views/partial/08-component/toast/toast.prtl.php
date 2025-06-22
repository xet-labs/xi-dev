<!-- Hidden checkbox to control the visibility of the toast -->
<input type="checkbox" style="display:none" id="id-toast-btn">

<?php if (session()->has('toast')) : ?>
    <?php $toast = session('toast'); ?>
    <div class="toast-wrap">
        <figure class="toast">
            <div class="toast-body">
                <div class="toast-icon-wrap">
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor" fill="none" stroke-linecap="round" stroke-linejoin="round"> 
                        <path stroke="none" d="M0 0h24v24H0z" fill="none"></path>
                        <path d="M5 12l5 5l10 -10"></path>
                    </svg>
                </div>

                <!-- Toast description -->
                <div class="toast-description">
                    <p><?= $toast['description']; ?></p>
                </div>

                <!-- Button if text is provided -->
                <?php if (!empty($toast['btnText'])): ?>
                    <button class="toast-btn">
                        <a href="<?= htmlspecialchars($toast['btnHref'] ?? '#', ENT_QUOTES, 'UTF-8'); ?>">
                            <?= htmlspecialchars($toast['btnText'], ENT_QUOTES, 'UTF-8'); ?>
                        </a>
                    </button>
                <?php endif; ?>
            </div>
            <div class="toast-progress"></div>
        </figure>
    </div>

    <!-- Clear the session flash after displaying -->
    <?php session()->forget('toast'); ?>
<?php endif; ?>
