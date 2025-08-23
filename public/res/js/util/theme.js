document.addEventListener("DOMContentLoaded", () => {
  // === Theme toggle (light/dark) ===
  document.documentElement.classList.add(localStorage.getItem('theme') || 'light');

  const themeSwitch = document.getElementById("id-themeswitch");
  if (themeSwitch) {
    themeSwitch.addEventListener("click", () => {
      document.documentElement.classList.toggle("dark");
      const theme = document.documentElement.classList.contains("dark") ? "dark" : "light";
      localStorage.setItem("theme", theme);
    });
  }
});