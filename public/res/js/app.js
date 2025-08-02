// CSRF token
$.ajaxSetup({ headers: { 'X-CSRF-TOKEN': $('meta[name="csrf-token"]').attr('content') } });

document.addEventListener("DOMContentLoaded", function () {
  // === Inject SVGs from [data-svg-src] ===
  document.querySelectorAll("[data-svg-src]").forEach(el => {
    const url = el.getAttribute("data-svg-src");
    if (url) {
      fetch(url)
        .then(res => res.text())
        .then(svg => {
          el.innerHTML = svg;
        })
        .catch(err => console.warn("SVG load failed:", url, err));
    }
  });

  // === Theme toggle (light/dark) ===
  const savedTheme = localStorage.getItem("theme");
  if (savedTheme === "dark" || savedTheme === "light") {
    document.documentElement.classList.add(savedTheme);
  }

  const themeSwitch = document.getElementById("id-themeswitch");
  if (themeSwitch) {
    themeSwitch.addEventListener("click", function () {
      document.documentElement.classList.toggle("dark");
      const theme = document.documentElement.classList.contains("dark") ? "dark" : "light";
      localStorage.setItem("theme", theme);
    });
  }

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

  // === Fix modal flicker for login/signup switch ===
  const loginBtn = document.getElementById('id-login-btn');
  const signupBtn = document.getElementById('id-signup-btn');
  const loginWrap = document.querySelector('.login-wrap');
  const signupWrap = document.querySelector('.signup-wrap');

  function toggleAuthForm() {
    if (loginBtn?.checked) {
      if (loginWrap) loginWrap.style.display = "flex";
      if (signupWrap) signupWrap.style.display = "";
    } else if (signupBtn?.checked) {
      if (signupWrap) signupWrap.style.display = "flex";
      if (loginWrap) loginWrap.style.display = "";
    }
  }

  if (loginBtn) loginBtn.addEventListener('click', toggleAuthForm);
  if (signupBtn) signupBtn.addEventListener('click', toggleAuthForm);

  // Nav hide|reveal on scroll
  let lastScrollTop = 0;
  const nav = document.querySelector('nav');
  const navHeight = nav.offsetHeight;
  
  window.addEventListener('scroll', function () {
    const currentScroll = window.pageYOffset || document.documentElement.scrollTop;
    
    if (currentScroll > lastScrollTop) {
      nav.style.top = `-${navHeight}px`;
    } else { nav.style.top = '0'; }
    
    lastScrollTop = currentScroll <= 0 ? 0 : currentScroll;
  });
  
  // === Optional: Add fix for scrollbar shift if needed ===
});

