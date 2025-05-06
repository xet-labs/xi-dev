// h2 onclick scroll to top
document.addEventListener('DOMContentLoaded', () => {
    const rootStyle = getComputedStyle(document.documentElement);
    const offset = 1.5 * parseFloat(rootStyle.fontSize);
    
    document.querySelectorAll('h2').forEach(heading => {
        heading.classList.add('h2-clickable');

        heading.addEventListener('click', () => {
            window.scrollTo({
                top: heading.getBoundingClientRect().top + window.scrollY - offset,
                behavior: 'smooth'
            });
        });
    });
});