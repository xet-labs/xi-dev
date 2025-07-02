document.addEventListener("DOMContentLoaded", function () {
    svg_copy = `
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="svg-icon lucide-copy"><rect x="8" y="8" width="14" height="14" rx="2" ry="2"></rect><path d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"></path></svg>
    `;
    svg_tick = `
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="svg-icon lucide-check">
            <path d="M20 6L9 17l-5-5"></path>
        </svg>
    `;

    document.querySelectorAll("pre:has(code)").forEach(pre => {
        // Copy button
        let button = document.createElement("button");
        button.className = "icon copy-code-btn";
        button.title = "Copy";
        button.innerHTML = svg_copy;
        // Copy fn
        button.addEventListener("click", () => {
            let code = pre.querySelector("code");
            const timeout = 1500;
            if (code) {
                navigator.clipboard.writeText(code.textContent).then(() => {
                    button.innerHTML = svg_tick;
                    button.style.color = "#44cf6e";
                    button.style.display = "flex";
                    // button.classList.add("hover");
                    // setTimeout(() => {button.classList.remove("hover");}, timeout - 80);
                    setTimeout(() => {
                        button.innerHTML = svg_copy;
                        button.style.color = "";
                        button.style.display = "";  
                    }, timeout);

                }).catch(err => console.error("Copy failed", err));
            }
        });

        pre.appendChild(button);
    });
});
