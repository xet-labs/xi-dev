/**
 * Inject content from elements with a specific attribute.
 * @param {string} attrName - Attribute to look for (e.g., "cntnt-src")
 * @param {Object} options
 * @param {"replace"|"append"|"prepend"} [options.mode="replace"] - How to insert content
 * @param {boolean} [options.observe=true] - Watch for dynamically added elements
 */
function injectContent(attrName = "cntnt-src", options = {}) {
    const { mode = "replace", observe = true } = options;

    function processElement(el) {
        const url = el.getAttribute(attrName);
        if (!url) return;
        
        fetch(url)
            .then(res => {
                if (!res.ok) throw new Error(`HTTP ${res.status}`);
                return res.text();
            })
            .then(content => {
                if (mode === "append") el.insertAdjacentHTML("beforeend", content);
                else if (mode === "prepend") el.insertAdjacentHTML("afterbegin", content);
                else el.innerHTML = content;
            })
            .catch(err => console.warn(`Failed to fetch [${attrName}] ${url}:`, err));
    }

    // Initial load
    document.querySelectorAll(`[${attrName}]`).forEach(processElement);

    // Watch for new elements
    if (observe) {
        const observer = new MutationObserver(mutations => {
            for (const mutation of mutations) {
                mutation.addedNodes.forEach(node => {
                    if (node.nodeType === 1 && node.hasAttribute?.(attrName)) {
                        processElement(node);
                    }
                    if (node.querySelectorAll) {
                        node.querySelectorAll(`[${attrName}]`).forEach(processElement);
                    }
                });
            }
        });
        observer.observe(document.body, { childList: true, subtree: true });
    }
}

// Auto-init on DOMContentLoaded
document.addEventListener("DOMContentLoaded", () => injectContent("cntnt-src"));
