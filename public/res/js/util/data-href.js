document.addEventListener("DOMContentLoaded", () => {
  // === Navigate via [data-href] on click or Enter key ===
  document.addEventListener('click', e => {
    const el = e.target.closest('[data-href]');
    if (el && !e.target.closest('a, button')) {
      window.location.href = el.dataset.href;
    }
  });
  document.addEventListener('keydown', e => {
    if (e.key === 'Enter') {
      const el = e.target.closest('[data-href]');
      if (el) {
        window.location.href = el.dataset.href;
      }
    }
  });
});