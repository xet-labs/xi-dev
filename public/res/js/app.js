document.addEventListener("DOMContentLoaded", () => {

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

